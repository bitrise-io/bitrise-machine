package cli

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/bitrise-tools/bitrise-machine/config"
	"github.com/bitrise-tools/bitrise-machine/session"
	"github.com/bitrise-tools/bitrise-machine/vagrant"
	"github.com/urfave/cli"
)

func getVagrantStatus(configModel config.MachineConfigModel, sessionStore session.StoreModel) (vagrant.MachineReadableItem, error) {
	// Read `vagrant status` log/output
	outputs, err := runVagrantCommandAndReturnCombinedOutput(configModel, sessionStore, "status", "--machine-readable")
	if err != nil {
		return vagrant.MachineReadableItem{}, fmt.Errorf("'vagrant status' failed. Output was: %s", outputs)
	}
	statusItms := vagrant.ParseMachineReadableItemsFromString(outputs, "", "state")
	if len(statusItms) != 1 {
		return vagrant.MachineReadableItem{}, fmt.Errorf("Failed to determine the 'status' of the machine. Output was: %s", outputs)
	}
	return statusItms[0], nil
}

// cleanupDestroyCommon ...
//  common code, cleanup's destroy
func cleanupDestroyCommon(configModel config.MachineConfigModel, sessionStore session.StoreModel) error {
	machineStatus, err := getVagrantStatus(configModel, sessionStore)
	if err != nil {
		return fmt.Errorf("Failed to get vagrant status: %s", err)
	}

	if machineStatus.Data != "not_created" {
		// destroy
		logrus.Infoln("Destroying machine...")
		if err := doDestroy(configModel, sessionStore); err != nil {
			return fmt.Errorf("'vagrant destroy' failed with error: %s", err)
		}
		logrus.Infoln("Machine destroyed.")
	} else {
		logrus.Infoln("Machine is in not-created state, skipping destroy.")
	}

	return nil
}

func doRecreateCleanup(configModel config.MachineConfigModel, previousSession session.StoreModel) (session.StoreModel, error) {
	// destroy
	if err := cleanupDestroyCommon(configModel, previousSession); err != nil {
		return previousSession, fmt.Errorf("doRecreateCleanup: failed to destroy: %s", err)
	}

	// re-create
	sessionStore, err := vagrantUpAndSessionInit(configModel)
	if err != nil {
		return sessionStore, fmt.Errorf("'vagrant up' failed with error: %s", err)
	}

	logrus.Infoln("Machine created and ready!")
	return sessionStore, nil
}

func doDestroyCleanup(configModel config.MachineConfigModel, sessionStore session.StoreModel) error {
	// destroy
	if err := cleanupDestroyCommon(configModel, sessionStore); err != nil {
		return fmt.Errorf("doDestroyCleanup: failed to destroy: %s", err)
	}

	logrus.Infoln("Machine destroyed, clean!")
	return nil
}

func doCustomCleanup(configModel config.MachineConfigModel, previousSession session.StoreModel) (session.StoreModel, error) {
	logrus.Infoln("Cleanup mode: custom-command")
	if configModel.CustomCleanupCommand == "" {
		return previousSession, errors.New("cleanup mode was custom-command, but no custom cleanup command specified")
	}
	logrus.Infof("=> Specified custom command: %s", configModel.CustomCleanupCommand)

	// Read `vagrant status` log/output
	machineStatus := vagrant.MachineReadableItem{}
	if outputs, err := runVagrantCommandAndReturnCombinedOutput(configModel, previousSession, "status", "--machine-readable"); err != nil {
		if err != nil {
			logrus.Errorf("'vagrant status' failed with output: %s", outputs)
			return previousSession, err
		}
	} else {
		statusItms := vagrant.ParseMachineReadableItemsFromString(outputs, "", "state")
		if len(statusItms) != 1 {
			return previousSession, fmt.Errorf("Failed to determine the 'status' of the machine. Output was: %s", outputs)
		}
		machineStatus = statusItms[0]
	}

	sessionStore := previousSession
	if machineStatus.Data == "not_created" {
		logrus.Infoln("Machine not yet created - creating with 'vagrant up'...")
		sessStore, err := vagrantUpAndSessionInit(configModel)
		if err != nil {
			return sessStore, fmt.Errorf("'vagrant up' failed with error: %s", err)
		}
		sessionStore = sessStore
		logrus.Infoln("Machine created!")
	} else {
		logrus.Infof("Machine already created - using the specified custom-command (%s) to clean it up...", configModel.CustomCleanupCommand)
		if err := runVagrantCommand(configModel, previousSession, configModel.CustomCleanupCommand); err != nil {
			return previousSession, fmt.Errorf("'vagrant %s' failed with error: %s", configModel.CustomCleanupCommand, err)
		}
		logrus.Infoln("Successful custom cleanup")
	}

	logrus.Infoln("Machine created and ready!")
	return sessionStore, nil
}

// doCleanup ...
// @isSkipHostCleanup : !!! should only be specified in case the host will be destroyed right after
//   the cleanup. 'will-be-destroyed' will leave the host as-it-is, uncleared!!
func doCleanup(configModel config.MachineConfigModel, isSkipHostCleanup string, sessionStore session.StoreModel) error {
	logrus.Infof("==> doCleanup (mode: %s)", configModel.CleanupMode)

	if isSkipHostCleanup != "will-be-destroyed" {
		switch configModel.CleanupMode {
		case config.CleanupModeRollback:
			if err := runVagrantCommand(configModel, sessionStore, "snapshot", "pop", "--no-delete"); err != nil {
				return err
			}
		case config.CleanupModeRecreate:
			sessStore, err := doRecreateCleanup(configModel, sessionStore)
			if err != nil {
				return err
			}
			sessionStore = sessStore
		case config.CleanupModeDestroy:
			if err := doDestroyCleanup(configModel, sessionStore); err != nil {
				return err
			}
		case config.CleanupModeCustomCommand:
			sessStore, err := doCustomCleanup(configModel, sessionStore)
			if err != nil {
				return err
			}
			sessionStore = sessStore
		default:
			return fmt.Errorf("Unsupported CleanupMode: %s", configModel.CleanupMode)
		}
	} else {
		logrus.Warnln("Skipping Host Cleanup! This option should only be used if the Host is destroyed immediately after this cleanup!!")
	}

	if err := config.DeleteSSHFilesFromDir(MachineWorkdir.Get()); err != nil {
		return fmt.Errorf("Failed to delete SSH file from workdir: %s", err)
	}

	return nil
}

func cleanup(c *cli.Context) {
	logrus.Infoln("Cleanup")

	// --- Configs and inputs

	additionalEnvs, err := config.CreateEnvItemsModelFromSlice(MachineParamsAdditionalEnvs.Get())
	if err != nil {
		logrus.Fatalf("Invalid Environment parameter: %s", err)
	}
	logrus.Debugf("additionalEnvs: %#v", additionalEnvs)

	configModel, err := config.ReadMachineConfigFileFromDir(MachineWorkdir.Get(), additionalEnvs)
	if err != nil {
		logrus.Fatalln("Failed to read Config file: ", err)
	}

	logrus.Infof("configModel: %#v", configModel)

	sessionStore, err := loadSessionIfSupportedAndExists(configModel.CleanupMode)
	if err != nil {
		logrus.Fatalf("Failed to load session, error: %s", err)
	}

	// ---

	if err := doCleanup(configModel, "", sessionStore); err != nil {
		logrus.Fatalf("Failed to Cleanup: %s", err)
	}

	logrus.Infoln("Cleanup - DONE - OK")
}
