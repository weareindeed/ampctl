package task

import (
	"ampctl/config"
	"ampctl/util"
	"fmt"
)

type HomebrewInstallTask struct {
	Config *config.Config
}

func NewHomebrewInstallTask() *HomebrewInstallTask {
	return &HomebrewInstallTask{}
}

func (t HomebrewInstallTask) Run() error {
	if path, ok := util.CheckPath(); ok {
		fmt.Printf("Homebrew is installed at: %s\n", path)
		return nil
	} else {
		return fmt.Errorf("Homebrew is NOT installed. Please install first before you go on")
	}
}
