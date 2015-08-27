package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	sshConfigFileName = "bitrise.machine.ssh.json"
)

// SSHConfigModel ...
type SSHConfigModel struct {
	IP           string `json:"ip"`
	Port         string `json:"port"`
	IdentityPath string `json:"identity_path"`
	Loginname    string `json:"loginname"`
}

func fullSSHConfigFilePath(dirPath string) string {
	return path.Join(dirPath, sshConfigFileName)
}

func (model *SSHConfigModel) normalizeAndValidate() error {
	if model.IP == "" {
		return fmt.Errorf("Invalid, empty IP")
	}
	if model.Port == "" {
		return fmt.Errorf("Invalid, empty Port")
	}
	if model.IdentityPath == "" {
		return fmt.Errorf("Invalid, empty IdentityPath")
	}
	if model.Loginname == "" {
		return fmt.Errorf("Invalid, empty Loginname")
	}

	return nil
}

func readSSHConfigFromBytes(configBytes []byte) (SSHConfigModel, error) {
	model := SSHConfigModel{}

	if err := json.Unmarshal(configBytes, &model); err != nil {
		return model, err
	}

	if err := model.normalizeAndValidate(); err != nil {
		return model, err
	}

	return model, nil
}

// DeleteSSHConfigFileFromDir ...
func DeleteSSHConfigFileFromDir(workdirPth string) error {
	fullConfPath := fullSSHConfigFilePath(workdirPth)
	isExists, err := pathutil.IsPathExists(fullConfPath)
	if err != nil {
		return err
	}
	if !isExists {
		return nil
	}
	return os.Remove(fullConfPath)
}

// ReadSSHConfigFileFromDir ...
func ReadSSHConfigFileFromDir(workdirPth string) (SSHConfigModel, error) {
	configBytes, err := fileutil.ReadBytesFromFile(fullSSHConfigFilePath(workdirPth))
	if err != nil {
		return SSHConfigModel{}, err
	}

	return readSSHConfigFromBytes(configBytes)
}

func (model SSHConfigModel) serializeIntoJSONBytes() ([]byte, error) {
	return json.MarshalIndent(model, "", "  ")
}

// WriteIntoFileInDir ...
func (model SSHConfigModel) WriteIntoFileInDir(workdirPth string) error {
	configBytes, err := model.serializeIntoJSONBytes()
	if err != nil {
		return err
	}

	if err := fileutil.WriteBytesToFile(fullSSHConfigFilePath(workdirPth), configBytes); err != nil {
		return err
	}

	return nil
}

// CreateSSHConfigFromVagrantSSHConfigLog ...
func CreateSSHConfigFromVagrantSSHConfigLog(vagrantSSHConfigLog string) (SSHConfigModel, error) {
	configModel := SSHConfigModel{}
	// process it line by line
	scanner := bufio.NewScanner(strings.NewReader(vagrantSSHConfigLog))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineSplits := strings.Split(line, " ")
		if len(lineSplits) == 2 {
			if lineSplits[0] == "HostName" {
				configModel.IP = lineSplits[1]
			} else if lineSplits[0] == "Port" {
				configModel.Port = lineSplits[1]
			} else if lineSplits[0] == "IdentityFile" {
				configModel.IdentityPath = lineSplits[1]
			} else if lineSplits[0] == "User" {
				configModel.Loginname = lineSplits[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return configModel, err
	}

	if err := configModel.normalizeAndValidate(); err != nil {
		return configModel, err
	}

	return configModel, nil
}

// SSHCommandArgs ...
func (model SSHConfigModel) SSHCommandArgs() []string {
	sshArgs := []string{model.IP, "-p", model.Port,
		"-oUserKnownHostsFile=/dev/null", "-oStrictHostKeyChecking=no",
		"-oPasswordAuthentication=no", "-oIdentitiesOnly=yes", "-oLogLevel=FATAL",
		"-l", model.Loginname, "-i", model.IdentityPath,
	}
	return sshArgs
}
