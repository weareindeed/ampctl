package main

import (
	"ampctl/config"
	"ampctl/task"
)

func NewContainer() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

type Container struct {
	services map[string]any
}

func (c *Container) GetConfig() *config.Config {
	if service, ok := c.services["GetConfig"]; ok {
		return service.(*config.Config)
	}

	instance := &config.Config{}
	c.services["GetConfig"] = instance
	return instance
}

func (c *Container) GetTask(name string) task.Task {
	switch name {
	case "brew:shivammathur:install":
		return c.GetShivammathurInstallTask()
	case "apache:install":
		return c.GetApacheInstallTask()
	case "apache:config:write":
		return c.GetApacheConfigWriteTask()
	case "hosts:write":
		return c.GetHostsWriteTask()
	case "php:install":
		return c.GetPhpInstallTask()
	case "root:privilege":
		return c.GetRootPrivilegeTask()
	}

	panic("Task not found")
}

func (c *Container) GetShivammathurInstallTask() *task.ShivammathurInstallTask {
	if service, ok := c.services["GetShivammathurInstallTask"]; ok {
		return service.(*task.ShivammathurInstallTask)
	}

	instance := &task.ShivammathurInstallTask{Config: c.GetConfig()}
	c.services["GetShivammathurInstallTask"] = instance
	return instance
}

func (c *Container) GetApacheInstallTask() *task.ApacheInstallTask {
	if service, ok := c.services["GetApacheInstallTask"]; ok {
		return service.(*task.ApacheInstallTask)
	}

	instance := &task.ApacheInstallTask{}
	c.services["GetApacheInstallTask"] = instance
	return instance
}

func (c *Container) GetApacheConfigWriteTask() *task.ApacheConfigWriteTask {
	if service, ok := c.services["GetApacheWriteConfigTask"]; ok {
		return service.(*task.ApacheConfigWriteTask)
	}

	instance := &task.ApacheConfigWriteTask{Config: c.GetConfig()}
	c.services["GetApacheWriteConfigTask"] = instance
	return instance
}

func (c *Container) GetHostsWriteTask() *task.HostsWriteTask {
	if service, ok := c.services["GetHostsWriteTask"]; ok {
		return service.(*task.HostsWriteTask)
	}

	instance := &task.HostsWriteTask{Config: c.GetConfig()}
	c.services["GetHostsWriteTask"] = instance
	return instance
}

func (c *Container) GetPhpInstallTask() *task.PhpInstallTask {
	if service, ok := c.services["GetPhpInstallTask"]; ok {
		return service.(*task.PhpInstallTask)
	}

	instance := &task.PhpInstallTask{Config: c.GetConfig()}
	c.services["GetPhpInstallTask"] = instance
	return instance
}

func (c *Container) GetRootPrivilegeTask() *task.RootPrivilegeTask {
	if service, ok := c.services["GetRootPrivilegeTask"]; ok {
		return service.(*task.RootPrivilegeTask)
	}

	instance := &task.RootPrivilegeTask{}
	c.services["GetRootPrivilegeTask"] = instance
	return instance
}
