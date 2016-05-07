package config

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/maputil"
)

const (
	machineConfigFileName = "bitrise.machine.config.json"

	// CleanupModeRollback ...
	CleanupModeRollback = "rollback"
	// CleanupModeRecreate ...
	CleanupModeRecreate = "recreate"
	// CleanupModeDestroy ...
	CleanupModeDestroy = "destroy"
	// CleanupModeCustomCommand ...
	CleanupModeCustomCommand = "custom-command"
)

// EnvItemsModel ...
type EnvItemsModel map[string]string

// MachineConfigModel ...
type MachineConfigModel struct {
	CleanupMode string `json:"cleanup_mode"`
	// IsCleanupBeforeSetup - if true do a cleanup before setup, unless
	//  if the host is already in "prepared" state
	// You can force to do a cleanup for every setup if you specify
	//  the --force flag as well for the setup command.
	IsCleanupBeforeSetup bool `json:"is_cleanup_before_setup"`
	// IsAllowVagrantCreateInSetup - if true `vagrant create` will be called
	//  in Setup in case the VM is not yet created
	IsAllowVagrantCreateInSetup bool `json:"is_allow_vagrant_create_in_setup"`
	// IsDoTimesyncAtSetup - OS X only at the moment
	IsDoTimesyncAtSetup  bool   `json:"is_do_timesync_at_setup"`
	CustomCleanupCommand string `json:"custom_cleanup_command"`
	// Envs - these will be set as Environment Variables
	//  for setup, cleanup and destroy
	Envs EnvItemsModel `json:"envs"`
	// ConfigTypeEnvs - these envs will be added to the
	//  other Envs based on the "config-type-id" paramter/flag
	ConfigTypeEnvs map[string]EnvItemsModel `json:"config_type_envs"`
}

func (configModel *MachineConfigModel) normalizeAndValidate() error {
	if configModel.CleanupMode != CleanupModeRollback &&
		configModel.CleanupMode != CleanupModeRecreate &&
		configModel.CleanupMode != CleanupModeCustomCommand &&
		configModel.CleanupMode != CleanupModeDestroy {
		return fmt.Errorf("Invalid CleanupMode: %s", configModel.CleanupMode)
	}

	if configModel.Envs == nil {
		configModel.Envs = EnvItemsModel{}
	}
	if configModel.ConfigTypeEnvs == nil {
		configModel.ConfigTypeEnvs = map[string]EnvItemsModel{}
	}

	return nil
}

// toCmdEnvs ...
func (envItmsModel *EnvItemsModel) toCmdEnvs() []string {
	res := make([]string, len(*envItmsModel))
	idx := 0
	for key, value := range *envItmsModel {
		res[idx] = fmt.Sprintf("%s=%s", key, value)
		idx++
	}
	return res
}

func (configModel MachineConfigModel) allEnvsForConfigType(configTypeID string) EnvItemsModel {
	allEnvs := EnvItemsModel{}

	if configModel.Envs != nil {
		allEnvs = maputil.CloneStringStringMap(configModel.Envs)
	}

	if configTypeID != "" {
		if configModel.ConfigTypeEnvs == nil {
			log.Warningf("No Config Type Envs defined, but Config Type ID was: %s", configTypeID)
		} else {
			configSpecEnvs, isFound := configModel.ConfigTypeEnvs[configTypeID]
			if !isFound {
				log.Warningf("No Config Type Envs found for the specified Config Type ID: %s", configTypeID)
			} else {
				allEnvs = maputil.MergeStringStringMap(allEnvs, configSpecEnvs)
			}
		}
	}

	return allEnvs
}

// AllCmdEnvsForConfigType ...
func (configModel MachineConfigModel) AllCmdEnvsForConfigType(configTypeID string) []string {
	allEnvs := configModel.allEnvsForConfigType(configTypeID)
	return allEnvs.toCmdEnvs()
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

func readMachineConfigFromBytes(configBytes []byte, appendEnvs EnvItemsModel) (MachineConfigModel, error) {
	configModel := MachineConfigModel{}

	if err := json.Unmarshal(configBytes, &configModel); err != nil {
		return configModel, err
	}

	if err := configModel.normalizeAndValidate(); err != nil {
		return configModel, err
	}

	for k, v := range appendEnvs {
		configModel.Envs[k] = v
	}

	return configModel, nil
}

// ReadMachineConfigFileFromDir ...
func ReadMachineConfigFileFromDir(workdirPth string, appendEnvs EnvItemsModel) (MachineConfigModel, error) {
	configBytes, err := fileutil.ReadBytesFromFile(path.Join(workdirPth, machineConfigFileName))
	if err != nil {
		return MachineConfigModel{}, fmt.Errorf("ReadMachineConfigFileFromDir: failed to read file: %s", err)
	}

	machineConfig, err := readMachineConfigFromBytes(configBytes, appendEnvs)
	if err != nil {
		return MachineConfigModel{}, fmt.Errorf("ReadMachineConfigFileFromDir: failed to parse configuration: %s", err)
	}

	return machineConfig, nil
}
