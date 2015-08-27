package cli

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-machine/config"
	"github.com/bitrise-io/bitrise-machine/utils"
	"github.com/codegangsta/cli"
)

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

	if err := utils.RunCommandThroughSSH(sshConfigModel, fullCmdToRunStr); err != nil {
		log.Fatalln("Failed to run command: ", err)
	}

	log.Infoln("Run finished - OK")
}
