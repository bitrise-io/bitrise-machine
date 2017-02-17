package cli

import (
	"fmt"

	"github.com/bitrise-tools/bitrise-machine/config"
	"github.com/bitrise-tools/bitrise-machine/session"
	"github.com/bitrise-tools/bitrise-machine/vagrant"
)

func allVagrantEnvs(configModel config.MachineConfigModel, sessionStore session.StoreModel) []string {
	configEnvs := configModel.AllCmdEnvsForConfigType(MachineConfigTypeID.Get())
	return append(configEnvs, sessionStore.Envs()...)
}

func vagrantUp(configModel config.MachineConfigModel, sessionStore session.StoreModel) error {
	envs := allVagrantEnvs(configModel, sessionStore)
	if err := vagrant.Up(MachineWorkdir.Get(), envs); err != nil {
		return fmt.Errorf("vagrantUp: Failed to vagrant up, error: %s", err)
	}
	return nil
}
