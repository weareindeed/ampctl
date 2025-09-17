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
		fmt.Println("Not added yet, so we add")
		err := util.AddTap("shivammathur/php")
		if err != nil {
			return fmt.Errorf("Error installing shivammathur/php")
		}
	} else {
		fmt.Println("Already added")
	}
	return nil
}
