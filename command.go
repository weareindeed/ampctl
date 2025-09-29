package main

import (
	"fmt"
)

func ProvisionCommand(container *Container) {
	fmt.Printf("Start provision\n")

	RunTaskAsRoot("root:privilege", container.GetConfig())
	RunTask("brew:shivammathur:install", container.GetConfig())
	RunTask("ssl:ca:generate", container.GetConfig())
	RunTask("ssl:hosts:generate", container.GetConfig())
	RunTask("apache:install", container.GetConfig())
	RunTask("apache:config:write", container.GetConfig())
	RunTask("apache:restart", container.GetConfig())
	RunTask("php:install", container.GetConfig())
	RunTask("php:config:write", container.GetConfig())
	RunTask("php:restart", container.GetConfig())
	RunTask("database:install", container.GetConfig())
	RunTask("database:config:write", container.GetConfig())
	RunTask("database:restart", container.GetConfig())
	RunTaskAsRoot("hosts:write", container.GetConfig())

	fmt.Printf("\nProvision finished!\n")
}

func DeployCommand(container *Container) {
	fmt.Printf("Start deploy\n")

	RunTaskAsRoot("root:privilege", container.GetConfig())
	RunTask("ssl:hosts:generate", container.GetConfig())
	RunTask("apache:config:write", container.GetConfig())
	RunTask("apache:restart", container.GetConfig())
	RunTask("php:config:write", container.GetConfig())
	RunTask("php:restart", container.GetConfig())
	RunTask("database:config:write", container.GetConfig())
	RunTask("database:restart", container.GetConfig())
	RunTaskAsRoot("hosts:write", container.GetConfig())

	fmt.Printf("\nDeploy finished!\n")
}

func RestartCommand(container *Container) {
	fmt.Printf("Restart\n")

	RunTaskAsRoot("root:privilege", container.GetConfig())
	RunTask("apache:restart", container.GetConfig())
	RunTask("php:restart", container.GetConfig())
	RunTask("database:restart", container.GetConfig())

	fmt.Printf("\nRestart finished!\n")
}

func StartCommand(container *Container) {
	fmt.Printf("Start\n")

	RunTaskAsRoot("root:privilege", container.GetConfig())
	RunTask("apache:start", container.GetConfig())
	RunTask("php:start", container.GetConfig())
	RunTask("database:start", container.GetConfig())

	fmt.Printf("\nStart finished!\n")
}

func StopCommand(container *Container) {
	fmt.Printf("Stop\n")

	RunTaskAsRoot("root:privilege", container.GetConfig())
	RunTask("apache:stop", container.GetConfig())
	RunTask("php:stop", container.GetConfig())
	RunTask("database:stop", container.GetConfig())

	fmt.Printf("\nStop finished!\n")
}
