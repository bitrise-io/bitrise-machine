package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"

	"github.com/bitrise-io/bitrise-machine/config"
)

// RunCommandThroughSSH ...
func RunCommandThroughSSH(sshConfigModel config.SSHConfigModel, cmdToRunWithSSH string) error {
	sshArgs := sshConfigModel.SSHCommandArgs()
	fullArgs := append(sshArgs, cmdToRunWithSSH)

	cmd := exec.Command("ssh", fullArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GenerateSSHKeypair ...
func GenerateSSHKeypair() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	privateKeyPemBytes := pem.EncodeToMemory(&privateKeyBlock)
	publicKey := privateKey.PublicKey

	pub, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	pubBytes := ssh.MarshalAuthorizedKey(pub)

	return privateKeyPemBytes, pubBytes, nil
}
