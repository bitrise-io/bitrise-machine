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
		return nil
	}

	if err := config.DeleteSSHConfigFileFromDir(MachineWorkdir); err != nil {
		return fmt.Errorf("Failed to delete SSH Configuration file: %s", err)
	}

	return fmt.Errorf("Unsupported CleanupMode: %s", configModel.CleanupMode)
}

func cleanup(c *cli.Context) {
	log.Infoln("Cleanup")
}
