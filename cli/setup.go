package cli

import (
	"fmt"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/bitrise-machine/config"
	"github.com/bitrise-tools/bitrise-machine/utils"
	"github.com/codegangsta/cli"
)

func doSetupSSH(configModel config.MachineConfigModel) (config.SSHConfigModel, error) {
	log.Infoln("==> doSetupSSH")
	sshConfigModel := config.SSHConfigModel{}

	// Read `vagrant ssh-config` log/output
	outputs, err := utils.RunAndReturnCombinedOutput(MachineWorkdir.Get(), configModel.Envs.ToCmdEnvs(), "vagrant", "ssh-config")
	if err != nil {
		log.Errorf("'vagrant ssh-config' failed with output: %s", outputs)
		return sshConfigModel, err
	}
	log.Debugln("===> (raw) vagrant ssh-config retrieved")

	// Convert `vagrant ssh-config` to our SSHConfigModel
	sshConfigModel, err = config.CreateSSHConfigFromVagrantSSHConfigLog(outputs)
	if err != nil {
		log.Errorf("'vagrant ssh-config' returned an invalid output (failed to scan SSH Config): %s", outputs)
		return sshConfigModel, err
	}
	log.Debugln("===> vagrant ssh-config parsed")

	// Generate SSH Keypair
	privBytes, pubBytes, err := utils.GenerateSSHKeypair()
	if err != nil {
		return sshConfigModel, err
	}
	log.Debugln("===> SSH Keypair generated")

	// Write the SSH Keypair to file
	privKeyFilePth, _, err := config.WriteSSHKeypairToFiles(MachineWorkdir.Get(), privBytes, pubBytes)
	if err != nil {
		return sshConfigModel, err
	}
	log.Debugln("===> SSH Keypair written to file")

	// Replace the ~/.ssh/authorized_keys inside the VM to only allow
	//  the new keypair
	replaceAuthKeysCmd := fmt.Sprintf(`printf "%s" > ~/.ssh/authorized_keys`, pubBytes)
	log.Debugf("===> Running command through SSH: %s", replaceAuthKeysCmd)
	if err := utils.RunCommandThroughSSH(sshConfigModel, replaceAuthKeysCmd); err != nil {
		return sshConfigModel, err
	}
	log.Debugln("===> SSH Keypair is now authorized to access the VM")

	// Save private key as the new identity
	sshConfigModel.IdentityPath = privKeyFilePth
	if err := sshConfigModel.WriteIntoFileInDir(MachineWorkdir.Get()); err != nil {
		return sshConfigModel, err
	}
	log.Debugln("===> New identity (private SSH key) saved into config in workdir")

	log.Debugln("==> doSetupSSH [done]")
	return sshConfigModel, nil
}

// doTimesync generates an appropriately formatted
//  time string, then calls `sudo date` through SSH
// This is required for Virtual Machines which are
//  simply rolled back to a snapshot state which might
//  mess up the VM's time (restores it to the snapshot time)
func doTimesync(sshConfigModel config.SSHConfigModel) error {
	log.Infoln("==> doTimesync")

	const layout = "2006-01-02 15:04:05 MST"
	timeNow := time.Now()
	timeAsString := timeNow.UTC().Format(layout)
	// THIS ONLY WORKS ON OSX VM!!
	timeSyncCmd := fmt.Sprintf("%s %s",
		`sudo date -uf "%Y-%m-%d %H:%M:%S UTC"`,
		`"`+timeAsString+`"`)
	log.Infoln("timeSyncCmd: ", timeSyncCmd)

	if err := utils.RunCommandThroughSSH(sshConfigModel, timeSyncCmd); err != nil {
		return err
	}

	return nil
}

func doCreateIfRequired(configModel config.MachineConfigModel) error {
	machineStatus, err := getVagrantStatus(configModel)
	if err != nil {
		return fmt.Errorf("Failed to get vagrant status: %s", err)
	}

	log.Debugf("doCreateIfRequired: machineStatus: %#v", machineStatus)

	if machineStatus.Type == "state" && machineStatus.Data == "not_created" {
		log.Infoln("Machine not yet created - creating with 'vagrant up'...")

		if err := utils.Run(MachineWorkdir.Get(), configModel.Envs.ToCmdEnvs(), "vagrant", "up"); err != nil {
			return fmt.Errorf("'vagrant up' failed with error: %s", err)
		}

		log.Infoln("Machine created!")
	}

	return nil
}

func setup(c *cli.Context) {
	log.Infoln("Setup")

	additionalEnvs, err := config.CreateEnvItemsModelFromSlice(MachineParamsAdditionalEnvs.Get())
	if err != nil {
		log.Fatalf("Invalid Environment parameter: %s", err)
	}

	configModel, err := config.ReadMachineConfigFileFromDir(MachineWorkdir.Get(), additionalEnvs)
	if err != nil {
		log.Fatalln("Failed to read Config file: ", err)
	}

	isOK, err := pathutil.IsPathExists(path.Join(MachineWorkdir.Get(), "Vagrantfile"))
	if err != nil {
		log.Fatalln("Failed to check 'Vagrantfile' in the WorkDir: ", err)
	}
	if !isOK {
		log.Fatalln("Vagrantfile not found in the WorkDir!")
	}

	log.Infof("configModel: %#v", configModel)

	isSkipSetups := false
	if config.IsSSHKeypairFileExistInDirectory(MachineWorkdir.Get()) && !c.Bool(ForceFlagKey) {
		log.Info("Host is already prepared and no --force flag was specified, skipping setup.")
		isSkipSetups = true
	}

	if !isSkipSetups {
		// doCleanup
		if configModel.IsCleanupBeforeSetup {
			if err := doCleanup(configModel, ""); err != nil {
				log.Fatalf("Failed to Cleanup: %s", err)
			}
		}

		if configModel.CleanupMode == config.CleanupModeDestroy || configModel.IsAllowVagrantCreateInSetup {
			if err := doCreateIfRequired(configModel); err != nil {
				log.Fatalf("Failed to Create the VM: %s", err)
			}
		}

		// ssh
		_, err := doSetupSSH(configModel)
		if err != nil {
			log.Fatalf("Failed to Setup SSH: %s", err)
		}
	}

	// time sync
	sshConfigModel, err := config.ReadSSHConfigFileFromDir(MachineWorkdir.Get())
	if err != nil {
		log.Fatalf("Failed to read SSH Config file! Error: %s", err)
	}
	if configModel.IsDoTimesyncAtSetup {
		if err := doTimesync(sshConfigModel); err != nil {
			log.Fatalf("Failed to do Time Sync: %s", err)
		}
	}

	log.Infoln("=> Setup DONE - OK")
}
