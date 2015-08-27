package utils

import "testing"

func Test_GenerateSSHKeypair(t *testing.T) {
	privBytes, pubBytes, err := GenerateSSHKeypair()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Logf("PRIVATE KEY: %s", privBytes)
	t.Logf("PUBLIC KEY: %s", pubBytes)
}
