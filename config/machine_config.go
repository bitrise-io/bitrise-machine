package config

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/bitrise-io/go-utils/fileutil"
)

const (
	machineConfigFileName = "bitrise.machine.config.json"

	// CleanupModeRollback ...
	CleanupModeRollback = "rollback"
	// CleanupModeRecreate ...
	CleanupModeRecreate = "recreate"
)

// MachineConfigModel ...
type MachineConfigModel struct {
	CleanupMode          string `json:"cleanup_mode"`
	IsCleanupBeforeSetup bool   `json:"is_cleanup_before_setup"`
	IsDoTimesyncAtSetup  bool   `json:"is_do_timesync_at_setup"`
}

func (configModel *MachineConfigModel) normalizeAndValidate() error {
	if configModel.CleanupMode != CleanupModeRollback && configModel.CleanupMode != CleanupModeRecreate {
		return fmt.Errorf("Invalid CleanupMode: %s", configModel.CleanupMode)
	}

	return nil
}

func readMachineConfigFromBytes(configBytes []byte) (MachineConfigModel, error) {
	configModel := MachineConfigModel{}

	if err := json.Unmarshal(configBytes, &configModel); err != nil {
		return configModel, err
	}

	if err := configModel.normalizeAndValidate(); err != nil {
		return configModel, err
	}

	return configModel, nil
}

// ReadMachineConfigFileFromDir ...
func ReadMachineConfigFileFromDir(workdirPth string) (MachineConfigModel, error) {
	configBytes, err := fileutil.ReadBytesFromFile(path.Join(workdirPth, machineConfigFileName))
	if err != nil {
		return MachineConfigModel{}, err
	}

	return readMachineConfigFromBytes(configBytes)
}
