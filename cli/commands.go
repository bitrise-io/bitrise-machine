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

	// --- Command flags

	// TimeoutFlagKey ...
	TimeoutFlagKey = "timeout"
	// LogFormatFlagKey ...
	LogFormatFlagKey = "logformat"
	// ForceFlagKey ...
	ForceFlagKey = "force"
)

var (
	commands = []cli.Command{
		{
			Name:   "setup",
			Usage:  "Setup/initialize the Host.",
			Action: setup,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  ForceFlagKey,
					Usage: "Force setup",
				},
			},
		},
		{
			Name:            "run",
			Usage:           "Run command on a Host - have to be initialized with setup first!",
			Action:          run,
			SkipFlagParsing: true,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  TimeoutFlagKey,
					Value: 0,
					Usage: "Timeout, in seconds",
				},
				cli.StringFlag{
					Name:  LogFormatFlagKey,
					Value: "",
					Usage: "Log format for the executed command's output. Default is 'no transform'. Options: json",
				},
			},
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
