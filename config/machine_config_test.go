package config

import (
	"testing"

	"github.com/bitrise-io/go-utils/pointers"
)

func Test_readMachineConfigFromBytes(t *testing.T) {
	configContent := `{
"cleanup_mode": "rollback",
"is_cleanup_before_setup": false
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
	if configModel.IsCleanupBeforeSetup == nil || *configModel.IsCleanupBeforeSetup != false {
		t.Fatal("Invalid IsCleanupBeforeSetup!")
	}
}

func Test_MachineConfigModel_normalizeAndValidate(t *testing.T) {
	configModel := MachineConfigModel{CleanupMode: ""}
	t.Log("Invalid CleanupMode")
	if err := configModel.normalizeAndValidate(); err == nil {
		t.Fatal("Should return a validation error!")
	}

	configModel = MachineConfigModel{CleanupMode: CleanupModeRollback}
	t.Logf("Default IsCleanupBeforeSetup: %#v", configModel)
	if err := configModel.normalizeAndValidate(); err != nil {
		t.Fatalf("Failed with error: %s", err)
	}
	if *configModel.IsCleanupBeforeSetup != true {
		t.Fatal("Invalid IsCleanupBeforeSetup - default value check")
	}

	configModel = MachineConfigModel{CleanupMode: CleanupModeRollback, IsCleanupBeforeSetup: pointers.NewBoolPtr(false)}
	t.Logf("IsCleanupBeforeSetup=false: %#v", configModel)
	if err := configModel.normalizeAndValidate(); err != nil {
		t.Fatalf("Failed with error: %s", err)
	}
	if *configModel.IsCleanupBeforeSetup != false {
		t.Fatal("Invalid IsCleanupBeforeSetup")
	}
}
