package config

import (
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/stretchr/testify/require"
)

func Test_EnvItemsModel_ToCmdEnvs(t *testing.T) {
	empty := EnvItemsModel{}
	require.Equal(t, []string{}, empty.ToCmdEnvs())

	one := EnvItemsModel{"key": "value"}
	require.Equal(t, []string{"key=value"}, one.ToCmdEnvs())

	two := EnvItemsModel{"key1": "value 1", "key2": "value 2"}
	testutil.EqualSlicesWithoutOrder(t, []string{"key1=value 1", "key2=value 2"}, two.ToCmdEnvs())

	envRef := EnvItemsModel{"key": "value with $HOME env ref"}
	require.Equal(t, []string{"key=value with $HOME env ref"}, envRef.ToCmdEnvs())
}

func Test_readMachineConfigFromBytes(t *testing.T) {
	configContent := `{
"cleanup_mode": "rollback",
"is_cleanup_before_setup": false,
"is_do_timesync_at_setup": true
}`

	t.Log("configContent: ", configContent)
	configModel, err := readMachineConfigFromBytes([]byte(configContent))
	if err != nil {
		t.Fatalf("Failed to read Config: %s", err)
	}
	t.Logf("configModel: %#v", configModel)

	if configModel.CleanupMode != "rollback" {
		t.Fatal("Invalid CleanupMode!")
	}
	if configModel.IsCleanupBeforeSetup != false {
		t.Fatal("Invalid IsCleanupBeforeSetup!")
	}
}

func Test_MachineConfigModel_normalizeAndValidate(t *testing.T) {
	configModel := MachineConfigModel{CleanupMode: ""}
	t.Log("Invalid CleanupMode")
	if err := configModel.normalizeAndValidate(); err == nil {
		t.Fatal("Should return a validation error!")
	}

	configModel = MachineConfigModel{
		CleanupMode:          CleanupModeRollback,
		IsCleanupBeforeSetup: true,
		IsDoTimesyncAtSetup:  false,
	}

	t.Logf("configModel: %#v", configModel)
	if err := configModel.normalizeAndValidate(); err != nil {
		t.Fatalf("Failed with error: %s", err)
	}
	if configModel.IsCleanupBeforeSetup != true {
		t.Fatal("Invalid IsCleanupBeforeSetup")
	}
	if configModel.IsDoTimesyncAtSetup != false {
		t.Fatal("Invalid IsDoTimesyncAtSetup")
	}
}
