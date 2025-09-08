package main

import (
	"ampctl/brew"
	"ampctl/config"
	"fmt"
)

func CheckCommand() {
	if path, ok := brew.CheckPath(); ok {
		fmt.Printf("Homebrew is installed at: %s\n", path)
	} else {
		fmt.Println("Homebrew is NOT installed. Please install first before you go on")
		return
	}

	cfg := &config.Config{}
	err := cfg.ImportHostsFileFromHomeDir()
	if err != nil {
		fmt.Printf("Failed to load hosts file: %v\n", err)
	}

	fmt.Printf("Hosts:")
	for _, host := range cfg.Hosts {
		fmt.Printf(" - Host: %s Path: %s Version: %s SSL: %t\n",
			host.Host, host.Path, host.Version, host.Ssl)
	}
}
