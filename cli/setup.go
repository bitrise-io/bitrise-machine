package cli

import (
	"fmt"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-machine/config"
	"github.com/bitrise-io/bitrise-machine/utils"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/codegangsta/cli"
)

func doSetupSSH(configModel config.MachineConfigModel) (config.SSHConfigModel, error) {
	log.Infoln("==> doSetupSSH")
	sshConfigModel := config.SSHConfigModel{}

	// Read `vagrant ssh-config` log/output
	outputs, err := cmdex.RunCommandInDirAndReturnCombinedStdoutAndStderr(MachineWorkdir, "vagrant", "ssh-config")
	if err != nil {
		log.Errorf("'vagrant ssh-config' failed with output: %s", outputs)
		return sshConfigModel, err
	}

	// Convert `vagrant ssh-config` to our SSHConfigModel
	sshConfigModel, err = config.CreateSSHConfigFromVagrantSSHConfigLog(outputs)
	if err != nil {
		log.Errorf("'vagrant ssh-config' returned an invalid output (failed to scan SSH Config): %s", outputs)
		return sshConfigModel, err
	}

	// Generate SSH Keypair
	privBytes, pubBytes, err := utils.GenerateSSHKeypair()
	if err != nil {
		return sshConfigModel, err
	}

	// Write the SSH Keypair to file
	privKeyFilePth, _, err := config.WriteSSHKeypairToFiles(MachineWorkdir, privBytes, pubBytes)
	if err != nil {
		return sshConfigModel, err
	}

	// Replace the ~/.ssh/authorized_keys inside the VM to only allow
	//  the new keypair
	replaceAuthKeysCmd := fmt.Sprintf(`printf "%s" > ~/.ssh/authorized_keys`, pubBytes)
	if err := utils.RunCommandThroughSSH(sshConfigModel, replaceAuthKeysCmd); err != nil {
		return sshConfigModel, err
	}

	// Save private key as the new identity
	sshConfigModel.IdentityPath = privKeyFilePth
	if err := sshConfigModel.WriteIntoFileInDir(MachineWorkdir); err != nil {
		return sshConfigModel, err
	}

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
	timeSyncCmd := fmt.Sprintf("%s %s",
		`sudo date -uf "%Y-%m-%d %H:%M:%S UTC"`,
		`"`+timeAsString+`"`)
	log.Infoln("timeSyncCmd: ", timeSyncCmd)

	if err := utils.RunCommandThroughSSH(sshConfigModel, timeSyncCmd); err != nil {
		return err
	}

	return nil
}

func setup(c *cli.Context) {
	log.Infoln("Setup")

	configModel, err := config.ReadMachineConfigFileFromDir(MachineWorkdir)
	if err != nil {
		log.Fatalln("Failed to read Config file: ", err)
	}

	isOK, err := pathutil.IsPathExists(path.Join(MachineWorkdir, "Vagrantfile"))
	if err != nil {
		log.Fatalln("Failed to check 'Vagrantfile' in the WorkDir: ", err)
	}
	if !isOK {
		log.Fatalln("Vagrantfile not found in the WorkDir!")
	}

	log.Infof("configModel: %#v", configModel)

	isSkipSetups := false
	if config.IsSSHKeypairFileExistInDirectory(MachineWorkdir) && !c.Bool(ForceFlagKey) {
		log.Info("Host is already prepared and no --force flag was specified, skipping setup.")
		isSkipSetups = true
	}

	if !isSkipSetups {
		// doCleanup
		if configModel.IsCleanupBeforeSetup {
			if err := doCleanup(configModel); err != nil {
				log.Fatalf("Failed to Cleanup: %s", err)
			}
		}

		// ssh
		_, err := doSetupSSH(configModel)
		if err != nil {
			log.Fatalf("Failed to Setup SSH: %s", err)
		}
	}

	// time sync
	sshConfigModel, err := config.ReadSSHConfigFileFromDir(MachineWorkdir)
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
