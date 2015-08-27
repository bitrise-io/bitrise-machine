package config

import (
	"strings"
	"testing"
)

func Test_readSSHConfigFromBytes(t *testing.T) {
	configContent := `{
"ip": "123.123.123.123",
"port": "123",
"identity_path": "/some/path",
"loginname": "usern"
}`

	t.Log("configContent: ", configContent)
	configModel, err := readSSHConfigFromBytes([]byte(configContent))
	if err != nil {
		t.Fatalf("Failed to read Config: %s", err)
	}
	t.Logf("configModel: %#v", configModel)

	if configModel.IP != "123.123.123.123" {
		t.Fatal("Invalid IP!")
	}
}

func Test_SSHConfigModel_normalizeAndValidate(t *testing.T) {
	configModel := SSHConfigModel{}
	t.Log("Invalid, empty SSH Config Model")
	if err := configModel.normalizeAndValidate(); err == nil {
		t.Fatal("Should return a validation error!")
	}

	configModel = SSHConfigModel{
		IP:           "123.123.123.123",
		Port:         "123",
		IdentityPath: "/some/path",
		Loginname:    "usern",
	}
	t.Logf("Proper SSH Config Model: %#v", configModel)
	if err := configModel.normalizeAndValidate(); err != nil {
		t.Fatalf("Failed with error: %s", err)
	}
}

func Test_SSHConfigModel_serializeIntoJSONBytes(t *testing.T) {
	configModel := SSHConfigModel{
		IP:           "123.321.123.321",
		Port:         "123",
		IdentityPath: "/some/path",
		Loginname:    "usern",
	}

	jsonBytes, err := configModel.serializeIntoJSONBytes()
	if err != nil {
		t.Fatalf("Failed to serialize SSH Config as JSON: %s", err)
	}

	expectedOutput := `{
  "ip": "123.321.123.321",
  "port": "123",
  "identity_path": "/some/path",
  "loginname": "usern"
}`

	if string(jsonBytes) != expectedOutput {
		t.Error("Generated JSON doesn't match the expected output.")
		t.Errorf("-> Expected: %s", expectedOutput)
		t.Fatalf("-> Got: %s", jsonBytes)
	}
}

func Test_CreateSSHConfigFromVagrantSSHConfigLog(t *testing.T) {
	vagrantSSHConfigOutput := `Host default
  HostName 123.321.123.321
  User vagrant
  Port 123
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /user/home/.vagrant.d/insecure_private_key
  IdentitiesOnly yes
  LogLevel FATAL
`

	t.Log("CreateSSHConfigFromVagrantSSHConfigLog")
	_, err := CreateSSHConfigFromVagrantSSHConfigLog(vagrantSSHConfigOutput)
	if err != nil {
		t.Fatalf("Failed to create SSH Config from Vagrant SSH Config Log: %s", err)
	}

	t.Log("CreateSSHConfigFromVagrantSSHConfigLog - Missing Port")
	vagrantSSHConfigOutput = `Host default
  HostName 123.321.123.321
  User vagrant
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /user/home/.vagrant.d/insecure_private_key
  IdentitiesOnly yes
  LogLevel FATAL
`
	_, err = CreateSSHConfigFromVagrantSSHConfigLog(vagrantSSHConfigOutput)
	if err == nil {
		t.Fatal("No error returned - should fail without a valid Port!")
	}
}

func Test_SSHConfigModel_SSHCommandArgs(t *testing.T) {
	configModel := SSHConfigModel{
		IP:           "123.321.123.321",
		Port:         "123",
		IdentityPath: "/some/path",
		Loginname:    "usern",
	}

	fullArgsStr := strings.Join(configModel.SSHCommandArgs(), "|")
	expectedStr := "123.321.123.321|-p|123|-oUserKnownHostsFile=/dev/null|-oStrictHostKeyChecking=no|-oPasswordAuthentication=no|-oIdentitiesOnly=yes|-oLogLevel=FATAL|-l|usern|-i|/some/path"
	if fullArgsStr != expectedStr {
		t.Error("Generated args doesn't match the expected output.")
		t.Errorf("-> Expected: %s", expectedStr)
		t.Fatalf("-> Got: %s", fullArgsStr)
	}
}
