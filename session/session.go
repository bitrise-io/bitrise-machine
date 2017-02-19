package session

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/bitrise-machine/config"
)

const (
	machineSessionFileName = "bitrise-machine-session.json"

	//
	machineSessionEnvVarTimeID = "BITRISE_MACHINE_SESSION_TIME_ID"
)

// StoreModel ...
type StoreModel struct {
	// SessionTimeID - in UTC, including date and time (up to seconds)
	// format: YYYYMMDDHHMMSS
	// example: 20170215093215
	SessionTimeID string `json:"session_time_id"`
}

func sessionTimeIDForTime(t time.Time) string {
	return t.UTC().Format("20060102150405")
}

// newStoreModel ...
func newStoreModel() StoreModel {
	return StoreModel{
		SessionTimeID: sessionTimeIDForTime(time.Now()),
	}
}

func readSessionStoreFromBytes(sessionStoreBytes []byte) (StoreModel, error) {
	storeModel := StoreModel{}

	if err := json.Unmarshal(sessionStoreBytes, &storeModel); err != nil {
		return storeModel, err
	}

	return storeModel, nil
}

// readStoreFileFromDir ...
func readStoreFileFromDir(workdirPth string) (StoreModel, error) {
	sessionStoreBytes, err := fileutil.ReadBytesFromFile(filepath.Join(workdirPth, machineSessionFileName))
	if err != nil {
		return StoreModel{}, fmt.Errorf("session.readStoreFileFromDir: failed to read file: %s", err)
	}

	machineSessionStore, err := readSessionStoreFromBytes(sessionStoreBytes)
	if err != nil {
		return StoreModel{}, fmt.Errorf("session.readStoreFileFromDir: failed to parse session store data: %s", err)
	}

	return machineSessionStore, nil
}

// writeStoreFileToDir ...
func (model StoreModel) writeStoreFileToDir(workdirPth string) error {
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("Failed to serialize session store, error: %s", err)
	}

	if err := fileutil.WriteBytesToFile(filepath.Join(workdirPth, machineSessionFileName), jsonBytes); err != nil {
		return fmt.Errorf("Failed to write session store into file, error: %s", err)
	}

	return nil
}

// Envs - environment variables defined by the session
func (model StoreModel) Envs() []string {
	return []string{
		machineSessionEnvVarTimeID + "=" + model.SessionTimeID,
	}
}

// Start ...
func Start(workdirPth string) (StoreModel, error) {
	sessionStore := newStoreModel()
	if err := sessionStore.writeStoreFileToDir(workdirPth); err != nil {
		return sessionStore, fmt.Errorf("session.Start error: %s", err)
	}
	return sessionStore, nil
}

// IsSessionSupportedForCleanupType ...
func IsSessionSupportedForCleanupType(cleanupType string) bool {
	switch cleanupType {
	case config.CleanupModeRecreate:
		return true
	case config.CleanupModeDestroy:
		return true
	}
	return false
}

// IsSessionStoreFileExists ...
func IsSessionStoreFileExists(workdirPth string) (bool, error) {
	sessionStoreFilePth := filepath.Join(workdirPth, machineSessionFileName)
	return pathutil.IsPathExists(sessionStoreFilePth)
}

// Load ...
func Load(workdirPth string, cleanupType string) (StoreModel, error) {
	if !IsSessionSupportedForCleanupType(cleanupType) {
		return StoreModel{}, fmt.Errorf("Session Store is not supported for cleanup type: %s", cleanupType)
	}
	return readStoreFileFromDir(workdirPth)
}
