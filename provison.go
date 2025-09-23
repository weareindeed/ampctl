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
	RunTaskAsRoot("hosts:write", container.GetConfig())
}
