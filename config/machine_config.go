package config

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

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

// CreateEnvItemsModelFromSlice ...
func CreateEnvItemsModelFromSlice(envsArr []string) (EnvItemsModel, error) {
	envItemsModel := EnvItemsModel{}
	for _, aEnvStr := range envsArr {
		splits := strings.Split(aEnvStr, "=")
		key := splits[0]
		if key == "" {
			return EnvItemsModel{}, fmt.Errorf("Invalid item, empty key. (Parameter was: %s)", aEnvStr)
		}
		if len(splits) < 2 {
			return EnvItemsModel{}, fmt.Errorf("Invalid item, no value defined. Key was: %s", splits[0])
		}

		value := strings.Join(splits[1:], "=")
		envItemsModel[key] = value
	}
	return envItemsModel, nil
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
func ReadMachineConfigFileFromDir(workdirPth string, appendEnvs EnvItemsModel) (MachineConfigModel, error) {
	configBytes, err := fileutil.ReadBytesFromFile(path.Join(workdirPth, machineConfigFileName))
	if err != nil {
		return MachineConfigModel{}, fmt.Errorf("ReadMachineConfigFileFromDir: failed to read file: %s", err)
	}

	machineConfig, err := readMachineConfigFromBytes(configBytes)
	if err != nil {
		return MachineConfigModel{}, fmt.Errorf("ReadMachineConfigFileFromDir: failed to parse configuration: %s", err)
	}

	for k, v := range appendEnvs {
		machineConfig.Envs[k] = v
	}

	return machineConfig, nil
}
