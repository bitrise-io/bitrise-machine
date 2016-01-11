package cli

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	// MachineWorkdir ...
	MachineWorkdir = ""
	// MachineParamsAdditionalEnvs ...
	MachineParamsAdditionalEnvs = []string{}
)

func before(c *cli.Context) error {
	// Log level
	if logLevel, err := log.ParseLevel(c.String(LogLevelKey)); err != nil {
		log.Fatal("Failed to parse log level:", err)
	} else {
		log.SetLevel(logLevel)
	}

	if len(c.Args()) != 0 && !c.Bool(HelpKey) && !c.Bool(VersionKey) {
		MachineWorkdir = c.String(WorkdirKey)
		if MachineWorkdir == "" {
			log.Fatalln("No Workdir specified!")
		}
	}

	MachineParamsAdditionalEnvs = c.StringSlice(EnvironmentParamKey)
	log.Debugf("MachineParamsAdditionalEnvs: %#v", MachineParamsAdditionalEnvs)

	return nil
}

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

// Run the Envman CLI.
func Run() {
	cli.VersionPrinter = printVersion

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "bitrise-machine"
	app.Version = "0.9.6"

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = appFlags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Finished with error:", err)
	}
}
