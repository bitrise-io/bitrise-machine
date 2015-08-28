package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-machine/config"
	"github.com/bitrise-io/bitrise-machine/utils"
	"github.com/codegangsta/cli"
)

const (
	logChunkRuneLenght = 100

	// LogFormatJSON ...
	LogFormatJSON = "json"
)

// RunResults ...
type RunResults struct {
	IsTimeoutError bool
	RunError       error
}

// LogBuffer ...
type LogBuffer struct {
	logBytes bytes.Buffer
	rwlock   sync.RWMutex
}

// Write ...
func (buff *LogBuffer) Write(p []byte) (n int, err error) {
	buff.rwlock.Lock()
	defer buff.rwlock.Unlock()
	return buff.logBytes.Write(p)
}

func (buff *LogBuffer) Read(n int) []byte {
	buff.rwlock.Lock()
	defer buff.rwlock.Unlock()
	return buff.logBytes.Next(n)
}

// ReadRunes ...
func (buff *LogBuffer) ReadRunes(n int) (string, bool) {
	buff.rwlock.Lock()
	defer buff.rwlock.Unlock()
	res := ""
	isEOF := false
	for i := 0; i < n; i++ {
		r, _, err := buff.logBytes.ReadRune()
		if err == nil {
			res += string(r)
		} else if err == io.EOF {
			isEOF = true
			break
		} else {
			log.Errorf("Failed to read from LogBuffer: %s", err)
		}
	}
	return res, isEOF
}

// LogChunkModel ...
type LogChunkModel struct {
	Data string `json:"data"`
	Pos  int64  `json:"pos"`
}

func logChunkJSONTransform(logChunkData string, logChunkIdx int64) ([]byte, error) {
	logChunk := LogChunkModel{
		Data: logChunkData,
		Pos:  logChunkIdx,
	}
	return json.Marshal(logChunk)
}

func performRun(sshConfig config.SSHConfigModel, commandToRunStr string, timeoutSeconds int64, logFormat string) RunResults {
	logBuff := LogBuffer{}
	var logChunkIndex int64
	lastLogChunkSentAt := time.Now()

	// Log processing
	processLogs := func(isFlush bool) (isLogChunkGenerated bool) {
		for {
			isChunkDone := false
			if logBuff.logBytes.Len() > logChunkRuneLenght || isFlush {
				chunkStr, isEOF := logBuff.ReadRunes(logChunkRuneLenght)
				if chunkStr != "" {
					if logFormat == LogFormatJSON {
						logChunkBytes, err := logChunkJSONTransform(chunkStr, logChunkIndex)
						logChunkIndex++
						if err != nil {
							log.Errorf("Failed to convert log chunk. Error: %s", err)
							log.Errorf(" Log chunk was: %s", chunkStr)
						}
						fmt.Printf("%s\n", logChunkBytes)
					} else {
						fmt.Printf("%s", chunkStr)
					}
					isLogChunkGenerated = true
				}
				if isEOF {
					isChunkDone = true
				}
			} else {
				isChunkDone = true
			}

			if !isFlush || isChunkDone {
				break
			}
		}
		return
	}

	// Run
	isRunFinished := false
	runRes := RunResults{
		IsTimeoutError: false,
		RunError:       nil,
	}
	c1 := make(chan RunResults, 1)
	go func() {
		if err := utils.RunCommandThroughSSHWithWriters(sshConfig, commandToRunStr, &logBuff, &logBuff); err != nil {
			log.Errorf("Failed to run command, error: %s", err)
			c1 <- RunResults{RunError: err, IsTimeoutError: false}
		} else {
			c1 <- RunResults{RunError: nil, IsTimeoutError: false}
		}
	}()

	var timeoutTriggerred <-chan time.Time
	if timeoutSeconds > 0 {
		log.Infof("Timeout registered with %d seconds from now.", timeoutSeconds)
		timeoutTriggerred = time.After(time.Duration(timeoutSeconds) * time.Second)
	}

	for !isRunFinished {
		select {
		case res := <-c1:
			if res.RunError == nil {
				runRes = res
			} else {
				runRes = res
			}
			isRunFinished = true
		case <-timeoutTriggerred:
			runRes = RunResults{RunError: fmt.Errorf("Timeout after %d seconds", timeoutSeconds), IsTimeoutError: true}
			isRunFinished = true
		case <-time.Tick(100 * time.Millisecond):
			isFlushLogs := false
			// force flush logs if we did not generated log chunks in the last
			//  couple of seconds - if the process did not generate enough logs
			//  to trigger a chunk generation
			timeDiffSec := time.Now().Sub(lastLogChunkSentAt).Seconds()
			if timeDiffSec > 5.0 {
				isFlushLogs = true
			}
			if processLogs(isFlushLogs) || isFlushLogs {
				lastLogChunkSentAt = time.Now()
			}
		}
	}

	processLogs(true)

	return runRes
}

func run(c *cli.Context) {
	log.Infoln("Run")

	if len(c.Args()) < 1 {
		log.Fatalln("No command to run specified!")
	}

	inCmdArgs := c.Args()
	cmdToRun := inCmdArgs[0]
	cmdToRunArgs := []string{}
	if len(inCmdArgs) > 1 {
		cmdToRunArgs = inCmdArgs[1:]
	}

	sshConfigModel, err := config.ReadSSHConfigFileFromDir(MachineWorkdir)
	if err != nil {
		log.Fatalln("Failed to read SSH configs - you should probably call 'setup' first!")
	}

	fullCmdToRunStr := fmt.Sprintf("%s %s", cmdToRun, strings.Join(cmdToRunArgs, " "))
	log.Infoln("fullCmdToRunStr: ", fullCmdToRunStr)

	timeoutSecs := c.Int(TimeoutFlagKey)
	if timeoutSecs > 0 {
		log.Infof("Timeout parameter: %d seconds", timeoutSecs)
	} else {
		log.Infoln("No timeout defined.")
		timeoutSecs = 0
	}

	logFormat := c.String(LogFormatFlagKey)
	if logFormat == "json" {
		log.Infof("Log format: %s", logFormat)
	} else {
		if logFormat != "" {
			log.Infof("Invalid Log Format - ignoring: %s", logFormat)
		}
		logFormat = ""
	}

	runResult := performRun(sshConfigModel, fullCmdToRunStr, int64(timeoutSecs), logFormat)
	if runResult.RunError != nil {
		if runResult.IsTimeoutError {
			log.Errorf("[!] Timeout: %s", runResult.RunError)
			os.Exit(2)
		} else {
			log.Errorf("Run failed: %s", runResult.RunError)
			os.Exit(1)
		}
	}

	log.Infoln("Run finished - OK")
}
