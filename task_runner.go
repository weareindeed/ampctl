package main

import (
	"ampctl/config"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

func RunTask(name string, config *config.Config) {
	fmt.Printf("=== %s === \n", name)
	exePath, _ := os.Executable()

	file := writeConfigFile(config)
	defer os.Remove(file)

	cmd := exec.Command(exePath, "--run-task", name, file)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func RunTaskAsRoot(name string, config *config.Config) {
	fmt.Printf("=== %s === \n", name)
	exePath, _ := os.Executable()

	file := writeConfigFile(config)
	defer os.Remove(file)

	cmd := exec.Command("sudo", "-S", exePath, "--run-task", name, file)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

// writeConfigFile will write the config as yaml to a tmp file and return the path
func writeConfigFile(config *config.Config) string {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	tmpFile, err := os.CreateTemp("", "task-*.yaml")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(bytes); err != nil {
		panic(err)
	}

	return tmpFile.Name()
}
