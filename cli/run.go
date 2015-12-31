package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-tools/bitrise-machine/config"
	"github.com/bitrise-tools/bitrise-machine/utils"
	"github.com/codegangsta/cli"
)

const (
	logChunkRuneLenght         = 10000 // ~ 10 KB
	abortCheckFrequencySeconds = 10.0

	buildFinishedWithErrorExitCode    = 10
	buildAbortedByTimeoutExitCode     = 2
	buildAbortedByUserRequestExitCode = 3

	logSummaryMetaInfoID = "/logs/summary"

	// LogFormatJSON ...
	LogFormatJSON = "json"
)

// RunResults ...
type RunResults struct {
	IsTimeoutError       bool
	IsUserRequestedAbort bool
	RunError             error
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

// LogSummaryModel ...
type LogSummaryModel struct {
	GeneratedChunkCount int64 `json:"generated_chunk_count"`
}

func printJSONControlMetaInfo(metaInfoID string, metaInfoObj interface{}) error {
	// format:
	//  :{{metaInfoID}}:{{metaInfoObj as json}}
	// ex: :/log/summary:{generated_chunk_count:123}

	metaInfoJSONBytes, err := json.Marshal(metaInfoObj)
	if err != nil {
		return err
	}
	fmt.Printf("\n:%s:%s\n", metaInfoID, metaInfoJSONBytes)
	return nil
}

func printLogSummary(logChunkNum int64) {
	logSummaryModel := LogSummaryModel{GeneratedChunkCount: logChunkNum}
	if err := printJSONControlMetaInfo(logSummaryMetaInfoID, logSummaryModel); err != nil {
		log.Errorf("Failed to generate Log Summary: %s", err)
	}
}

// AbortCheckModel ...
type AbortCheckModel struct {
	StatusStr    string `json:"status"`
	IsAborted    bool   `json:"is_aborted"`
	ErrorMessage string `json:"error_msg"`
}

func requestGetJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println(" [!] Failed to close r.Body:", err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

func performRun(sshConfig config.SSHConfigModel, commandToRunStr string,
	timeoutSeconds int64, abortCheckURL string, logFormat string,
) RunResults {

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
						if err != nil {
							log.Errorf("Failed to convert log chunk. Error: %s", err)
							log.Errorf(" Log chunk was: %s", chunkStr)
						} else {
							fmt.Printf("%s\n", logChunkBytes)
							logChunkIndex++
						}
					} else {
						fmt.Printf("%s", chunkStr)
						logChunkIndex++
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

	logTickFn := func() {
		isFlushLogs := false
		// force flush logs if we did not generated log chunks in the last
		//  couple of seconds - if the process did not generate enough logs
		//  to trigger a chunk generation
		timeDiffSec := time.Now().Sub(lastLogChunkSentAt).Seconds()
		if timeDiffSec > 3.0 {
			isFlushLogs = true
		}
		if processLogs(isFlushLogs) || isFlushLogs {
			lastLogChunkSentAt = time.Now()
		}
	}

	lastAbortCheckAt := time.Now()
	abortCheckTickFN := func() {
		if abortCheckURL == "" {
			// no abort check URL defined
			return
		}
		timeDiffSec := time.Now().Sub(lastAbortCheckAt).Seconds()
		if timeDiffSec < abortCheckFrequencySeconds {
			return
		}
		lastAbortCheckAt = time.Now()
		log.Debug("=> Abort check")

		abortCheckModel := AbortCheckModel{}
		jsonErr := requestGetJSON(
			abortCheckURL,
			&abortCheckModel)
		log.Debugf("==> Result: %#v | err: %s", abortCheckModel, jsonErr)

		if jsonErr == nil && abortCheckModel.StatusStr == "ok" && abortCheckModel.IsAborted {
			runRes = RunResults{RunError: fmt.Errorf("Build was aborted"), IsUserRequestedAbort: true}
			isRunFinished = true
		}
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
			logTickFn()
			abortCheckTickFN()
		}
	}

	processLogs(true)
	printLogSummary(logChunkIndex)

	return runRes
}

func run(c *cli.Context) {
	log.Infoln("Run")

	if len(c.Args()) < 1 {
		log.Fatalln("No command to run specified!")
	}

	inCmdArgs := c.Args()
	log.Debugf("inCmdArgs: %v", inCmdArgs)
	cmdToRun := inCmdArgs[0]
	cmdToRunArgs := []string{}
	if len(inCmdArgs) > 1 {
		cmdToRunArgs = inCmdArgs[1:]
	}

	sshConfigModel, err := config.ReadSSHConfigFileFromDir(MachineWorkdir)
	if err != nil {
		log.Fatalln("Failed to read SSH configs - you should probably call 'setup' first!")
	}

	// fullCmdToRunStr := fmt.Sprintf("%s %s", cmdToRun, strings.Join(cmdToRunArgs, " "))
	fullCmdToRunStr := cmdex.LogPrintableCommandArgs(append([]string{cmdToRun}, cmdToRunArgs...))
	log.Infoln("fullCmdToRunStr: ", fullCmdToRunStr)

	timeoutSecs := c.Int(TimeoutFlagKey)
	if timeoutSecs > 0 {
		log.Infof("Timeout parameter: %d seconds", timeoutSecs)
	} else {
		log.Infoln("No timeout defined.")
		timeoutSecs = 0
	}

	abortCheckURL := c.String(AbortCheckURLFlagKey)
	if abortCheckURL != "" {
		log.Infof("Abort check URL: %s", abortCheckURL)
	} else {
		log.Infoln("No abort-check-URL defined.")
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

	runResult := performRun(sshConfigModel, fullCmdToRunStr, int64(timeoutSecs), abortCheckURL, logFormat)
	if runResult.RunError != nil {
		if runResult.IsTimeoutError {
			log.Errorf("[!] Timeout: %s", runResult.RunError)
			os.Exit(buildAbortedByTimeoutExitCode)
		} else if runResult.IsUserRequestedAbort {
			log.Errorf("[!] User requested abort: %s", runResult.RunError)
			os.Exit(buildAbortedByUserRequestExitCode)
		} else {
			log.Errorf("Run failed: %s", runResult.RunError)
			os.Exit(buildFinishedWithErrorExitCode)
		}
	}

	log.Infoln("Run finished - OK")
}
