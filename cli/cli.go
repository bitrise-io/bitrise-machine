package cli

import (
	"fmt"
	"os"
	"path"
	"runtime/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/freezable"
	"github.com/codegangsta/cli"
)

var (
	// MachineWorkdir ...
	MachineWorkdir = freezable.String{}
	// MachineParamsAdditionalEnvs ...
	MachineParamsAdditionalEnvs = freezable.StringSlice{}
	// MachineConfigTypeID ...
	MachineConfigTypeID = freezable.String{}
)

func before(c *cli.Context) error {
	// Log level
	if logLevel, err := log.ParseLevel(c.String(LogLevelKey)); err != nil {
		log.Fatal("Failed to parse log level:", err)
	} else {
		log.SetLevel(logLevel)
	}

	if len(c.Args()) != 0 && !c.Bool(HelpKey) && !c.Bool(VersionKey) {
		if err := MachineWorkdir.Set(c.String(WorkdirKey)); err != nil {
			log.Fatalf("Failed to set MachineWorkdir: %s", err)
		}
		if MachineWorkdir.String() == "" {
			log.Fatalln("No Workdir specified!")
		}
	}
	MachineWorkdir.Freeze()

	if err := MachineConfigTypeID.Set(c.String(ConfigTypeIDParamKey)); err != nil {
		log.Fatalf("Failed to set MachineConfigTypeID: %s", err)
	}
	log.Debugf("MachineConfigTypeID: %s", MachineConfigTypeID)

	if err := MachineParamsAdditionalEnvs.Set(c.StringSlice(EnvironmentParamKey)); err != nil {
		log.Fatalf("Failed to set MachineParamsAdditionalEnvs: %s", err)
	}
	log.Debugf("MachineParamsAdditionalEnvs: %s", MachineParamsAdditionalEnvs)
	MachineParamsAdditionalEnvs.Freeze()

	return nil
}

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

// Run the CLI
func Run() {
	cli.VersionPrinter = printVersion

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "bitrise-machine"
	app.Version = "0.9.9"

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = appFlags
	app.Commands = commands

	if cpuProfileFilePath := os.Getenv(CPUProfileFilePathEnvKey); cpuProfileFilePath != "" {
		log.Infof("Enabling CPU Profiler, writing results into file: %s", cpuProfileFilePath)
		f, err := os.Create(cpuProfileFilePath)
		if err != nil {
			log.Fatalf("CPU Profile: failed to create file: %s", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("CPU Profile: failed to start CPU profiler: %s", err)
		}
		defer pprof.StopCPUProfile()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Finished with error:", err)
	}
}
