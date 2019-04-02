package cli

import (
	"fmt"

	"github.com/bitrise-io/bitrise-machine/config"
	"github.com/bitrise-io/bitrise-machine/session"
	"github.com/bitrise-io/bitrise-machine/utils"

	"github.com/Sirupsen/logrus"
)

func allVagrantEnvs(configModel config.MachineConfigModel, sessionStore session.StoreModel) []string {
	configEnvs := configModel.AllCmdEnvsForConfigType(MachineConfigTypeID.Get())
	return append(configEnvs, sessionStore.Envs()...)
}

func vagrantUpAndSessionInit(configModel config.MachineConfigModel) (session.StoreModel, error) {
	sessionStore := session.StoreModel{}
	if session.IsSessionSupportedForCleanupType(configModel.CleanupMode) {
		sessStore, err := session.Start(MachineWorkdir.Get())
		if err != nil {
			return session.StoreModel{}, fmt.Errorf("Failed to start Session, error: %s", err)
		}
		sessionStore = sessStore
	} else {
		logrus.Infof(" (i) Session handling is not supported for cleanup type: %s", configModel.CleanupMode)
	}

	return sessionStore, runVagrantCommand(configModel, sessionStore, "up")
}

func runVagrantCommand(configModel config.MachineConfigModel, sessionStore session.StoreModel, vagrantCommandArgs ...string) error {
	envs := allVagrantEnvs(configModel, sessionStore)

	if err := utils.Run(MachineWorkdir.Get(), envs, "vagrant", vagrantCommandArgs...); err != nil {
		return fmt.Errorf("vagrant command (%s) failed, error: %s", vagrantCommandArgs, err)
	}
	return nil
}

func runVagrantCommandAndReturnCombinedOutput(configModel config.MachineConfigModel, sessionStore session.StoreModel, vagrantCommandArgs ...string) (string, error) {
	envs := allVagrantEnvs(configModel, sessionStore)

	output, err := utils.RunAndReturnCombinedOutput(MachineWorkdir.Get(), envs, "vagrant", vagrantCommandArgs...)
	if err != nil {
		return output, fmt.Errorf("vagrant command (%s) failed, error: %s", vagrantCommandArgs, err)
	}
	return output, nil
}

// loadSessionIfSupportedAndExists loads the session IF:
// 1. if session is supported for the current Cleanup Mode
//     - if not supported by cleanup mode, prints an info/debug log and returns an empty session store
// 2. if the session file exists
//     - if the session file does not exist, prints a warning and returns an empty session store,
//       except if it fails to determine whether the session store file exists (in which case returns the error)
func loadSessionIfSupportedAndExists(cleanupMode string) (session.StoreModel, error) {
	if !session.IsSessionSupportedForCleanupType(cleanupMode) {
		logrus.Infof(" (i) Session handling is not supported for cleanup type: %s", cleanupMode)
		return session.StoreModel{}, nil
	}

	if isSessExist, err := session.IsSessionStoreFileExists(MachineWorkdir.Get()); err != nil {
		return session.StoreModel{}, fmt.Errorf("Failed to determine whether session exists, error: %s", err)
	} else if !isSessExist {
		logrus.Warn(" (!) Session (file) does not exist")
		return session.StoreModel{}, nil
	}

	sessStore, err := session.Load(MachineWorkdir.Get(), cleanupMode)
	if err != nil {
		return session.StoreModel{}, fmt.Errorf("Failed to load session, error: %s", err)
	}
	return sessStore, nil
}
