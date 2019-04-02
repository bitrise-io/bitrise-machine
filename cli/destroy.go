package cli

import (
	"fmt"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/bitrise-machine/config"
	"github.com/bitrise-io/bitrise-machine/session"
	"github.com/urfave/cli"
)

func doDestroy(configModel config.MachineConfigModel, sessionStore session.StoreModel) error {
	logrus.Infoln("==> doDestroy")

	if err := runVagrantCommand(configModel, sessionStore, "destroy", "-f"); err != nil {
		return fmt.Errorf("'vagrant destroy' failed with error: %s", err)
	}

	return nil
}

func destroy(c *cli.Context) {
	logrus.Infoln("Destroy")

	// --- Configs and inputs

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

	sessionStore, err := loadSessionIfSupportedAndExists(configModel.CleanupMode)
	if err != nil {
		logrus.Fatalf("Failed to load session, error: %s", err)
	}

	// ---

	if sessStore, err := doCleanup(configModel, "will-be-destroyed", sessionStore); err != nil {
		logrus.Fatalf("Failed to Cleanup: %s", err)
	} else {
		sessionStore = sessStore
	}

	if err := doDestroy(configModel, sessionStore); err != nil {
		logrus.Fatalf("Failed to Destroy: %s", err)
	}

	logrus.Infoln("=> Destroy DONE - OK")
}
