package config

import "testing"

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
