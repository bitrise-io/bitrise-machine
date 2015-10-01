package utils

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// RunAndReturnCombinedOutput ...
func RunAndReturnCombinedOutput(dir string, envs []string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = append(os.Environ(), envs...)
	outBytes, err := cmd.CombinedOutput()

	if err != nil {
		log.Errorf("Command (args: %#v) failed, outputs was: %s", cmd.Args, outBytes)
		return "", err
	}

	outStr := string(outBytes)
	return strings.TrimSpace(outStr), nil
}

// Run ...
func Run(dir string, envs []string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("Command failed. Args: %#v", cmd.Args)
		return err
	}
	return nil
}
