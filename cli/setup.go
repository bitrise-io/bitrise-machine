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

	//
	outputs, err := cmdex.RunCommandInDirAndReturnCombinedStdoutAndStderr(MachineWorkdir, "vagrant", "ssh-config")
	if err != nil {
		log.Errorf("'vagrant ssh-config' failed with output: %s", outputs)
		return sshConfigModel, err
	}

	sshConfigModel, err = config.CreateSSHConfigFromVagrantSSHConfigLog(outputs)
	if err != nil {
		log.Errorf("'vagrant ssh-config' returned an invalid output (failed to scan SSH Config): %s", outputs)
		return sshConfigModel, err
	}

	if err := sshConfigModel.WriteIntoFileInDir(MachineWorkdir); err != nil {
		return sshConfigModel, err
	}

	sshConfigModel, err = config.ReadSSHConfigFileFromDir(MachineWorkdir)
	if err != nil {
		return sshConfigModel, err
	}

	log.Fatalln("--> IMPLEMENT SSH KEY HANDLING!!")

	return sshConfigModel, nil
}

func doTimesync(sshConfigModel config.SSHConfigModel) error {
	log.Infoln("==> doTimesync")

	const layout = "2006-01-02 15:04:05 MST"
	timeNow := time.Now()
	timeAsString := timeNow.UTC().Format(layout)
	log.Infoln(" (debug) timeAsString: ", timeAsString)
	timeSyncCmd := fmt.Sprintf("%s %s",
		`sudo date -uf "%Y-%m-%d %H:%M:%S UTC"`,
		`"`+timeAsString+`"`)
	log.Infoln("timeSyncCmd: ", timeSyncCmd)

	if err := utils.RunCommandThroughSSH(sshConfigModel, timeSyncCmd); err != nil {
		return err
	}

	// if err := bitruncommon.RunCommandThroughSSHWithWriters([]string{timeSyncCmd},
	// 	prov.sshConfig, logStreamerOut, logStreamerErr, prov.IsVerbose); err != nil {
	// 	return err
	// }
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

	// doCleanup
	if *configModel.IsCleanupBeforeSetup {
		if err := doCleanup(configModel); err != nil {
			log.Fatalf("Failed to Cleanup: %s", err)
		}
	}

	// ssh
	sshConfigModel, err := doSetupSSH(configModel)
	if err != nil {
		log.Fatalf("Failed to Setup SSH: %s", err)
	}

	// time sync
	if err := doTimesync(sshConfigModel); err != nil {
		log.Fatalf("Failed to do Time Sync: %s", err)
	}

	log.Infoln("=> Setup DONE - OK")
}
