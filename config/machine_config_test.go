package config

import (
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/stretchr/testify/require"
)

func Test_EnvItemsModel_ToCmdEnvs(t *testing.T) {
	empty := EnvItemsModel{}
	require.Equal(t, []string{}, empty.toCmdEnvs())

	one := EnvItemsModel{"key": "value"}
	require.Equal(t, []string{"key=value"}, one.toCmdEnvs())

	two := EnvItemsModel{"key1": "value 1", "key2": "value 2"}
	testutil.EqualSlicesWithoutOrder(t, []string{"key1=value 1", "key2=value 2"}, two.toCmdEnvs())

	envRef := EnvItemsModel{"key": "value with $HOME env ref"}
	require.Equal(t, []string{"key=value with $HOME env ref"}, envRef.toCmdEnvs())
}

func Test_readMachineConfigFromBytes(t *testing.T) {
	configContent := `{
"cleanup_mode": "rollback",
"is_cleanup_before_setup": false,
"is_do_timesync_at_setup": true
}`

	t.Log("configContent: ", configContent)

	t.Log("Base Config")
	{
		configModel, err := readMachineConfigFromBytes([]byte(configContent), EnvItemsModel{})
		require.NoError(t, err)
		t.Logf("configModel: %#v", configModel)

		if configModel.CleanupMode != "rollback" {
			t.Fatal("Invalid CleanupMode!")
		}
		if configModel.IsCleanupBeforeSetup != false {
			t.Fatal("Invalid IsCleanupBeforeSetup!")
		}

		require.Equal(t, []string{}, configModel.Envs.toCmdEnvs())
		require.Equal(t, EnvItemsModel{}, configModel.Envs)
		require.Equal(t, map[string]EnvItemsModel{}, configModel.ConfigTypeEnvs)
	}

	t.Log("Additional Env items")
	{
		configModel, err := readMachineConfigFromBytes([]byte(configContent), EnvItemsModel{"key": "my value"})
		require.NoError(t, err)
		t.Logf("configModel: %#v", configModel)
		require.Equal(t, EnvItemsModel{"key": "my value"}, configModel.Envs)
	}

	t.Log("Additional Env items - overwrite a config defined Env")
	{
		configContent := `{
"cleanup_mode": "rollback",
"is_cleanup_before_setup": false,
"is_do_timesync_at_setup": true,
"envs": {
  "MY_KEY": "config value"
}
}`
		configModel, err := readMachineConfigFromBytes([]byte(configContent), EnvItemsModel{"MY_KEY": "additional env value"})
		require.NoError(t, err)
		t.Logf("configModel: %#v", configModel)
		require.Equal(t, EnvItemsModel{"MY_KEY": "additional env value"}, configModel.Envs)
	}

	t.Log("Config specific Env items")
	{
		configContent := `{
"cleanup_mode": "rollback",
"is_cleanup_before_setup": false,
"is_do_timesync_at_setup": true,
"envs": {
  "MY_KEY": "config value"
},
"config_type_envs": {
  "config-1": {
    "CONFIG_1_KEY": "config 1 value"
  },
  "config-2": {
    "CONFIG_2_KEY": "config 2 value"
  }
}
}`
		configModel, err := readMachineConfigFromBytes([]byte(configContent), EnvItemsModel{})
		require.NoError(t, err)
		t.Logf("configModel: %#v", configModel)
		require.Equal(t, EnvItemsModel{"MY_KEY": "config value"}, configModel.Envs)
		require.Equal(t, map[string]EnvItemsModel{
			"config-1": EnvItemsModel{"CONFIG_1_KEY": "config 1 value"},
			"config-2": EnvItemsModel{"CONFIG_2_KEY": "config 2 value"},
		}, configModel.ConfigTypeEnvs)
	}

	t.Log("Env Var priority/overwrite : Config Specific one wins")
	{
		configContent := `{
"cleanup_mode": "rollback",
"is_cleanup_before_setup": false,
"is_do_timesync_at_setup": true,
"envs": {
  "MY_KEY": "config value",
  "ENV_KEY_OVERWRITE": "value from Envs"
},
"config_type_envs": {
  "config-1": {
    "CONFIG_1_KEY": "config 1 value",
	"ENV_KEY_OVERWRITE": "value from ConfigTypeEnvs"
  }
}
}`
		configModel, err := readMachineConfigFromBytes([]byte(configContent), EnvItemsModel{"ENV_KEY_OVERWRITE": "User defined value"})
		require.NoError(t, err)
		t.Logf("configModel: %#v", configModel)
		require.Equal(t, EnvItemsModel{
			"MY_KEY":            "config value",
			"ENV_KEY_OVERWRITE": "User defined value",
		}, configModel.Envs)
		require.Equal(t, map[string]EnvItemsModel{
			"config-1": EnvItemsModel{
				"CONFIG_1_KEY":      "config 1 value",
				"ENV_KEY_OVERWRITE": "value from ConfigTypeEnvs",
			},
		}, configModel.ConfigTypeEnvs)

		t.Log("config type envs overwrite base envs, even the user defined ones")
		{
			allEnvs := configModel.allEnvsForConfigType("config-1")
			require.Equal(t, EnvItemsModel{
				"MY_KEY":            "config value",
				"CONFIG_1_KEY":      "config 1 value",
				"ENV_KEY_OVERWRITE": "value from ConfigTypeEnvs",
			}, allEnvs)
		}

		t.Log("but it does not overwrite the value if no config-type or invalid one is defined")
		{
			allEnvs := configModel.allEnvsForConfigType("")
			require.Equal(t, EnvItemsModel{
				"MY_KEY":            "config value",
				"ENV_KEY_OVERWRITE": "User defined value",
			}, allEnvs)
		}
	}
}

func Test_MachineConfigModel_normalizeAndValidate(t *testing.T) {

	t.Log("Invalid CleanupMode")
	{
		configModel := MachineConfigModel{CleanupMode: ""}
		require.EqualError(t, configModel.normalizeAndValidate(), "Invalid CleanupMode: ")
	}

	t.Log("Minimal valid")
	{
		configModel := MachineConfigModel{
			CleanupMode:          CleanupModeRollback,
			IsCleanupBeforeSetup: true,
			IsDoTimesyncAtSetup:  false,
		}

		t.Logf("configModel: %#v", configModel)
		require.NoError(t, configModel.normalizeAndValidate())
		require.Equal(t, true, configModel.IsCleanupBeforeSetup)
		require.Equal(t, false, configModel.IsDoTimesyncAtSetup)
		// defaults
		require.Equal(t, EnvItemsModel{}, configModel.Envs)
		require.Equal(t, map[string]EnvItemsModel{}, configModel.ConfigTypeEnvs)
	}
}

func Test_allEnvsForConfigType(t *testing.T) {
	t.Log("Empty")
	{
		machineConfig := MachineConfigModel{}
		allEnvs := machineConfig.allEnvsForConfigType("")
		require.Equal(t, EnvItemsModel{}, allEnvs)
	}

	t.Log("Only Envs")
	{
		machineConfig := MachineConfigModel{
			Envs: EnvItemsModel{
				"ENV_KEY_1": "env value 1",
			},
		}
		allEnvs := machineConfig.allEnvsForConfigType("")
		require.Equal(t, EnvItemsModel{
			"ENV_KEY_1": "env value 1",
		}, allEnvs)
	}

	t.Log("Only Config Specific Envs")
	{
		machineConfig := MachineConfigModel{
			ConfigTypeEnvs: map[string]EnvItemsModel{
				"config-1": EnvItemsModel{
					"CONFIG_1_ENV_KEY_1": "config 1 - env value 1",
				},
			},
		}
		allEnvs := machineConfig.allEnvsForConfigType("config-1")
		require.Equal(t, EnvItemsModel{
			"CONFIG_1_ENV_KEY_1": "config 1 - env value 1",
		}, allEnvs)
	}

	t.Log("Invalid ConfigTypeID - no envs specified for it")
	{
		machineConfig := MachineConfigModel{
			ConfigTypeEnvs: map[string]EnvItemsModel{
				"config-1": EnvItemsModel{
					"CONFIG_1_ENV_KEY_1": "config 1 - env value 1",
				},
			},
		}
		// empty config type ID
		allEnvs := machineConfig.allEnvsForConfigType("")
		require.Equal(t, EnvItemsModel{}, allEnvs)
		// non existing config type ID
		allEnvs = machineConfig.allEnvsForConfigType("not-defined-config-type-id")
		require.Equal(t, EnvItemsModel{}, allEnvs)
	}

	// test: priority of envs - which one wins?
	t.Log("Priority of Envs - ConfigType Envs overwrite base Envs (other Envs also defined)")
	{
		machineConfig := MachineConfigModel{
			Envs: EnvItemsModel{
				"ENV_KEY_1":         "env value 1",
				"ENV_KEY_OVERWRITE": "value from Envs",
			},
			ConfigTypeEnvs: map[string]EnvItemsModel{
				"config-1": EnvItemsModel{
					"CONFIG_1_ENV_KEY_1": "config 1 - env value 1",
					"ENV_KEY_OVERWRITE":  "value from ConfigTypeEnvs",
				},
			},
		}
		allEnvs := machineConfig.allEnvsForConfigType("config-1")
		require.Equal(t, EnvItemsModel{
			"ENV_KEY_1":          "env value 1",
			"CONFIG_1_ENV_KEY_1": "config 1 - env value 1",
			"ENV_KEY_OVERWRITE":  "value from ConfigTypeEnvs",
		}, allEnvs)
	}

	t.Log("Priority of Envs - ConfigType Envs overwrite base Envs (no other Envs/keys)")
	{
		machineConfig := MachineConfigModel{
			Envs: EnvItemsModel{
				"ENV_KEY_OVERWRITE": "value from Envs",
			},
			ConfigTypeEnvs: map[string]EnvItemsModel{
				"config-1": EnvItemsModel{
					"ENV_KEY_OVERWRITE": "value from ConfigTypeEnvs",
				},
			},
		}
		allEnvs := machineConfig.allEnvsForConfigType("config-1")
		require.Equal(t, EnvItemsModel{
			"ENV_KEY_OVERWRITE": "value from ConfigTypeEnvs",
		}, allEnvs)
	}
}
func TestAllCmdEnvsForConfigType(t *testing.T) {
	t.Log("Empty")
	{
		machineConfig := MachineConfigModel{}
		allCmdEnvs := machineConfig.AllCmdEnvsForConfigType("")
		require.Equal(t, []string{}, allCmdEnvs)
	}

	t.Log("Priority of Envs - ConfigType Envs overwrite base Envs (other Envs also defined)")
	{
		machineConfig := MachineConfigModel{
			Envs: EnvItemsModel{
				"ENV_KEY_1":         "env value 1",
				"ENV_KEY_OVERWRITE": "value from Envs",
			},
			ConfigTypeEnvs: map[string]EnvItemsModel{
				"config-1": EnvItemsModel{
					"CONFIG_1_ENV_KEY_1": "config 1 - env value 1",
					"ENV_KEY_OVERWRITE":  "value from ConfigTypeEnvs",
				},
			},
		}
		allCmdEnvs := machineConfig.AllCmdEnvsForConfigType("config-1")
		testutil.EqualSlicesWithoutOrder(t, []string{
			"ENV_KEY_1=env value 1",
			"CONFIG_1_ENV_KEY_1=config 1 - env value 1",
			"ENV_KEY_OVERWRITE=value from ConfigTypeEnvs",
		}, allCmdEnvs)
	}
}

func TestCreateEnvItemsModelFromSlice(t *testing.T) {
	t.Log("Empty")
	{
		envsItmModel, err := CreateEnvItemsModelFromSlice([]string{})
		require.NoError(t, err)
		require.Equal(t, EnvItemsModel{}, envsItmModel)
	}

	t.Log("One item - but invalid, empty")
	{
		_, err := CreateEnvItemsModelFromSlice([]string{""})
		require.EqualError(t, err, "Invalid item, empty key. (Parameter was: )")
	}

	t.Log("One item - but invalid, value provided but empty key")
	{
		_, err := CreateEnvItemsModelFromSlice([]string{"=hello"})
		require.EqualError(t, err, "Invalid item, empty key. (Parameter was: =hello)")
	}

	t.Log("One item, no value - error")
	{
		_, err := CreateEnvItemsModelFromSlice([]string{"a"})
		require.EqualError(t, err, "Invalid item, no value defined. Key was: a")
	}

	t.Log("One item, with value")
	{
		envsItmModel, err := CreateEnvItemsModelFromSlice([]string{"a=b"})
		require.NoError(t, err)
		require.Equal(t, EnvItemsModel{"a": "b"}, envsItmModel)
	}

	t.Log("One item, with empty value")
	{
		envsItmModel, err := CreateEnvItemsModelFromSlice([]string{"a="})
		require.NoError(t, err)
		require.Equal(t, EnvItemsModel{"a": ""}, envsItmModel)
	}

	t.Log("One item, with value which includes spaces")
	{
		envsItmModel, err := CreateEnvItemsModelFromSlice([]string{"a=b c  d"})
		require.NoError(t, err)
		require.Equal(t, EnvItemsModel{"a": "b c  d"}, envsItmModel)
	}

	t.Log("One item, with value which includes equal signs")
	{
		envsItmModel, err := CreateEnvItemsModelFromSlice([]string{"a=b c=d  =e"})
		require.NoError(t, err)
		require.Equal(t, EnvItemsModel{"a": "b c=d  =e"}, envsItmModel)
	}

	t.Log("Multiple values")
	{
		envsItmModel, err := CreateEnvItemsModelFromSlice([]string{"a=b c d", "1=2 3 4"})
		require.NoError(t, err)
		require.Equal(t, EnvItemsModel{"a": "b c d", "1": "2 3 4"}, envsItmModel)
	}
}
