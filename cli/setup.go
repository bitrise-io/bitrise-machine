package cli

import (
	"fmt"
	"path"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/bitrise-machine/config"
	"github.com/bitrise-io/bitrise-machine/session"
	"github.com/bitrise-io/bitrise-machine/utils"
	"github.com/urfave/cli"
)

func doSetupSSH(configModel config.MachineConfigModel, sessionStore session.StoreModel) (config.SSHConfigModel, error) {
	logrus.Infoln("==> doSetupSSH")
	sshConfigModel := config.SSHConfigModel{}

	// Read `vagrant ssh-config` log/output
	outputs, err := runVagrantCommandAndReturnCombinedOutput(configModel, sessionStore, "ssh-config")
	if err != nil {
		logrus.Errorf("'vagrant ssh-config' failed with output: %s", outputs)
		return sshConfigModel, err
	}
	logrus.Debugln("===> (raw) vagrant ssh-config retrieved")

	// Convert `vagrant ssh-config` to our SSHConfigModel
	sshConfigModel, err = config.CreateSSHConfigFromVagrantSSHConfigLog(outputs)
	if err != nil {
		logrus.Errorf("'vagrant ssh-config' returned an invalid output (failed to scan SSH Config): %s", outputs)
		return sshConfigModel, err
	}
	logrus.Debugln("===> vagrant ssh-config parsed")

	// Generate SSH Keypair
	privBytes, pubBytes, err := utils.GenerateSSHKeypair()
	if err != nil {
		return sshConfigModel, err
	}
	logrus.Debugln("===> SSH Keypair generated")

	// Write the SSH Keypair to file
	privKeyFilePth, _, err := config.WriteSSHKeypairToFiles(MachineWorkdir.Get(), privBytes, pubBytes)
	if err != nil {
		return sshConfigModel, err
	}
	logrus.Debugln("===> SSH Keypair written to file")

	// Replace the ~/.ssh/authorized_keys inside the VM to only allow
	//  the new keypair
	replaceAuthKeysCmd := fmt.Sprintf(`printf "%s" > ~/.ssh/authorized_keys`, pubBytes)
	logrus.Debugf("===> Running command through SSH: %s", replaceAuthKeysCmd)
	if err := utils.RunCommandThroughSSH(sshConfigModel, replaceAuthKeysCmd); err != nil {
		return sshConfigModel, err
	}
	logrus.Debugln("===> SSH Keypair is now authorized to access the VM")

	// Save private key as the new identity
	sshConfigModel.IdentityPath = privKeyFilePth
	if err := sshConfigModel.WriteIntoFileInDir(MachineWorkdir.Get()); err != nil {
		return sshConfigModel, err
	}
	logrus.Debugln("===> New identity (private SSH key) saved into config in workdir")

	logrus.Debugln("==> doSetupSSH [done]")
	return sshConfigModel, nil
}

// doTimesync generates an appropriately formatted
//  time string, then calls `sudo date` through SSH
// This is required for Virtual Machines which are
//  simply rolled back to a snapshot state which might
//  mess up the VM's time (restores it to the snapshot time)
func doTimesync(sshConfigModel config.SSHConfigModel) error {
	logrus.Infoln("==> doTimesync")

	const layout = "2006-01-02 15:04:05 MST"
	timeNow := time.Now()
	timeAsString := timeNow.UTC().Format(layout)
	// THIS ONLY WORKS ON OSX VM!!
	timeSyncCmd := fmt.Sprintf("%s %s",
		`sudo date -uf "%Y-%m-%d %H:%M:%S UTC"`,
		`"`+timeAsString+`"`)
	logrus.Infoln("timeSyncCmd: ", timeSyncCmd)

	if err := utils.RunCommandThroughSSH(sshConfigModel, timeSyncCmd); err != nil {
		return err
	}

	return nil
}

// doCreateIfRequired creates the VM if necessary (if it does not exist),
// and initializes a new session; if already created it does nothing
func doCreateIfRequired(configModel config.MachineConfigModel, previousSession session.StoreModel) (session.StoreModel, error) {
	machineStatus, err := getVagrantStatus(configModel, previousSession)
	if err != nil {
		return previousSession, fmt.Errorf("Failed to get vagrant status: %s", err)
	}

	logrus.Debugf("doCreateIfRequired: machineStatus: %#v", machineStatus)

	sessionStore := previousSession
	if machineStatus.Type == "state" && machineStatus.Data == "not_created" {
		logrus.Infoln("Machine not yet created - creating with 'vagrant up'...")

		sessStore, err := vagrantUpAndSessionInit(configModel)
		if err != nil {
			return sessionStore, fmt.Errorf("doCreateIfRequired: Failed to vagrant up, error: %s", err)
		}
		sessionStore = sessStore

		logrus.Infoln("Machine created!")
	}

	return sessionStore, nil
}

func setup(c *cli.Context) {
	logrus.Infoln("Setup")

	additionalEnvs, err := config.CreateEnvItemsModelFromSlice(MachineParamsAdditionalEnvs.Get())
	if err != nil {
		logrus.Fatalf("Invalid Environment parameter: %s", err)
	}

	configModel, err := config.ReadMachineConfigFileFromDir(MachineWorkdir.Get(), additionalEnvs)
	if err != nil {
		logrus.Fatalln("Failed to read Config file: ", err)
	}

	isOK, err := pathutil.IsPathExists(path.Join(MachineWorkdir.Get(), "Vagrantfile"))
	if err != nil {
		logrus.Fatalln("Failed to check 'Vagrantfile' in the WorkDir: ", err)
	}
	if !isOK {
		logrus.Fatalln("Vagrantfile not found in the WorkDir!")
	}

	logrus.Infof("configModel: %#v", configModel)

	previousSession, err := loadSessionIfSupportedAndExists(configModel.CleanupMode)
	if err != nil {
		logrus.Fatalf("Failed to load session, error: %s", err)
	}

	isSkipSetups := false
	if config.IsSSHKeypairFileExistInDirectory(MachineWorkdir.Get()) && !c.Bool(ForceFlagKey) {
		logrus.Info("Host is already prepared and no --force flag was specified, skipping setup.")
		isSkipSetups = true
	}

	if !isSkipSetups {
		// doCleanup
		sessionStore := previousSession
		if configModel.IsCleanupBeforeSetup {
			if sessStore, err := doCleanup(configModel, "", sessionStore); err != nil {
				logrus.Fatalf("Failed to Cleanup: %s", err)
			} else {
				sessionStore = sessStore
			}
		}

		if configModel.CleanupMode == config.CleanupModeDestroy || configModel.IsAllowVagrantCreateInSetup {
			if sessStore, err := doCreateIfRequired(configModel, sessionStore); err != nil {
				logrus.Fatalf("Failed to Create the VM: %s", err)
			} else {
				// overwrite session with the new
				sessionStore = sessStore
			}
		}

		// ssh
		_, err := doSetupSSH(configModel, sessionStore)
		if err != nil {
			logrus.Fatalf("Failed to Setup SSH: %s", err)
		}
	}

	// time sync
	sshConfigModel, err := config.ReadSSHConfigFileFromDir(MachineWorkdir.Get())
	if err != nil {
		logrus.Fatalf("Failed to read SSH Config file! Error: %s", err)
	}
	if configModel.IsDoTimesyncAtSetup {
		if err := doTimesync(sshConfigModel); err != nil {
			logrus.Fatalf("Failed to do Time Sync: %s", err)
		}
	}

	logrus.Infoln("=> Setup DONE - OK")
}
