package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-machine/config"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/codegangsta/cli"
)

func doCleanup(configModel config.MachineConfigModel) error {
	log.Infoln("==> doCleanup")

	if configModel.CleanupMode == config.CleanupModeRollback {
		if err := cmdex.RunCommandInDir(MachineWorkdir, "vagrant", "sandbox", "rollback"); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unsupported CleanupMode: %s", configModel.CleanupMode)
	}

	if err := config.DeleteSSHFilesFromDir(MachineWorkdir); err != nil {
		return fmt.Errorf("Failed to delete SSH file from workdir: %s", err)
	}

	return nil
}

func cleanup(c *cli.Context) {
	log.Infoln("Cleanup")

	configModel, err := config.ReadMachineConfigFileFromDir(MachineWorkdir)
	if err != nil {
		log.Fatalln("Failed to read Config file: ", err)
	}

	log.Infof("configModel: %#v", configModel)

	if err := doCleanup(configModel); err != nil {
		log.Fatalf("Failed to Cleanup: %s", err)
	}

	log.Infoln("Cleanup - DONE - OK")
}
