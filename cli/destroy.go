package cli

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/bitrise-machine/config"
	"github.com/bitrise-tools/bitrise-machine/utils"
	"github.com/codegangsta/cli"
)

func doDestroy(configModel config.MachineConfigModel) error {
	log.Infoln("==> doDestroy")

	if err := utils.Run(MachineWorkdir, configModel.Envs.ToCmdEnvs(), "vagrant", "destroy", "-f"); err != nil {
		return fmt.Errorf("'vagrant destroy' failed with error: %s", err)
	}

	return nil
}

func destroy(c *cli.Context) {
	log.Infoln("Destroy")

	additionalEnvs, err := config.CreateEnvItemsModelFromSlice(MachineParamsAdditionalEnvs)
	if err != nil {
		log.Fatalf("Invalid Environment parameter: %s", err)
	}

	configModel, err := config.ReadMachineConfigFileFromDir(MachineWorkdir, additionalEnvs)
	if err != nil {
		log.Fatalln("Failed to read Config file: ", err)
	}

	isOK, err := pathutil.IsPathExists(path.Join(MachineWorkdir, "Vagrantfile"))
	if err != nil {
		log.Fatalln("Failed to check 'Vagrantfile' in the WorkDir: ", err)
	}
	if !isOK {
		log.Fatalln("Vagrantfile not found in the WorkDir!")
	}

	log.Infof("configModel: %#v", configModel)

	if err := doCleanup(configModel, "will-be-destroyed"); err != nil {
		log.Fatalf("Failed to Cleanup: %s", err)
	}

	if err := doDestroy(configModel); err != nil {
		log.Fatalf("Failed to Destroy: %s", err)
	}

	log.Infoln("=> Destroy DONE - OK")
}
