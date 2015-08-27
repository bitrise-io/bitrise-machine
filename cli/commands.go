package cli

import "github.com/codegangsta/cli"

const (
	// WorkdirEnvKey ...
	WorkdirEnvKey = "BITRISE_MACHINE_WORKDIR"
	// WorkdirKey ...
	WorkdirKey = "workdir"

	// LogLevelEnvKey ...
	LogLevelEnvKey = "LOGLEVEL"
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

	// HelpKey ...
	HelpKey      = "help"
	helpKeyShort = "h"

	// VersionKey ...
	VersionKey      = "version"
	versionKeyShort = "v"
)

var (
	commands = []cli.Command{
		{
			Name:   "setup",
			Usage:  "Setup/initialize the Host.",
			Action: setup,
		},
		{
			Name:            "run",
			Usage:           "Run command on a Host - have to be initialized with setup first!",
			Action:          run,
			SkipFlagParsing: true,
		},
		{
			Name:   "cleanup",
			Usage:  "Cleanup the Host.",
			Action: cleanup,
		},
	}

	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:   LogLevelKey + ", " + logLevelKeyShort,
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: LogLevelEnvKey,
		},
		cli.StringFlag{
			Name:   WorkdirKey,
			Value:  "",
			Usage:  "Working & config directory path.",
			EnvVar: WorkdirEnvKey,
		},
	}
)

func init() {
	// Override default help and version flags
	cli.HelpFlag = cli.BoolFlag{
		Name:  HelpKey + ", " + helpKeyShort,
		Usage: "Show help.",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  VersionKey + ", " + versionKeyShort,
		Usage: "Print the version.",
	}
}
