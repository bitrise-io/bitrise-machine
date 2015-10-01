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
	// CleanupModeCustomCommand ...
	CleanupModeCustomCommand = "custom-command"
)

// EnvItemsModel ...
type EnvItemsModel map[string]string

// MachineConfigModel ...
type MachineConfigModel struct {
	CleanupMode          string        `json:"cleanup_mode"`
	IsCleanupBeforeSetup bool          `json:"is_cleanup_before_setup"`
	IsDoTimesyncAtSetup  bool          `json:"is_do_timesync_at_setup"`
	CustomCleanupCommand string        `json:"custom_cleanup_command"`
	Envs                 EnvItemsModel `json:"envs"`
}

func (configModel *MachineConfigModel) normalizeAndValidate() error {
	if configModel.CleanupMode != CleanupModeRollback &&
		configModel.CleanupMode != CleanupModeRecreate &&
		configModel.CleanupMode != CleanupModeCustomCommand {
		return fmt.Errorf("Invalid CleanupMode: %s", configModel.CleanupMode)
	}

	return nil
}

// ToCmdEnvs ...
func (envItmsModel *EnvItemsModel) ToCmdEnvs() []string {
	res := make([]string, len(*envItmsModel))
	idx := 0
	for key, value := range *envItmsModel {
		res[idx] = fmt.Sprintf("%s=%s", key, value)
		idx++
	}
	return res
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
