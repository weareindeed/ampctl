package task

import (
	"ampctl/config"
	"ampctl/util"
	"fmt"
)

type ShivammathurInstallTask struct {
	Config *config.Config
}

func (t *ShivammathurInstallTask) Run() error {
	fmt.Print("Check if shivammathur/php is added: ")
	if !util.IsTapInstalled("shivammathur/php") {
		fmt.Print("Not added yet, so we add\n")
		err := util.AddTap("shivammathur/php")
		if err != nil {
			return fmt.Errorf("Error installing shivammathur/php")
		}
	} else {
		fmt.Print("Already added\n")
	}

	fmt.Print("Check if shivammathur/extensions is added: ")
	if !util.IsTapInstalled("shivammathur/extensions") {
		fmt.Print("Not added yet, so we add\n")
		err := util.AddTap("shivammathur/extensions")
		if err != nil {
			return fmt.Errorf("Error installing shivammathur/extensions")
		}
	} else {
		fmt.Print("Already added\n")
	}

	return nil
}
