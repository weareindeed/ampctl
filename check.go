package main

import (
	"ampctl/util"
	"fmt"
)

func CheckCommand(container *Container) {
	if path, ok := util.CheckPath(); ok {
		fmt.Printf("Homebrew is installed at: %s\n", path)
	} else {
		fmt.Println("Homebrew is NOT installed. Please install first before you go on")
		return
	}

	fmt.Printf("Hosts:")
	for _, host := range container.GetConfig().Hosts {
		fmt.Printf(" - Host: %s Path: %s Version: %s SSL: %t\n",
			host.Host, host.Path, host.Version, host.Ssl)
	}
}
