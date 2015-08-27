package utils

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
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

	log.Infof("Cmd: %#v", cmd)

	return cmd.Run()
}
