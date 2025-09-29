package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var container = NewContainer()

	if len(os.Args) > 1 && os.Args[1] == "--run-task" {
		err := container.GetConfig().ImportFile(os.Args[3])
		if err != nil {
			fmt.Println(err)
		}

		runTask(container, os.Args[2])
		return
	}

	if os.Geteuid() == 0 {
		fmt.Println("Don't run ampctl as root")
		return
	}

	err := container.GetConfig().LoadConfig()
	if err != nil {
		fmt.Println(err)
	}

	var rootCmd = &cobra.Command{Use: "ampctl"}

	var checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Check settings",
		Run: func(cmd *cobra.Command, args []string) {
			CheckCommand(container)
		},
	}

	var provisionCmd = &cobra.Command{
		Use:   "provision",
		Short: "Run provision",
		Run: func(cmd *cobra.Command, args []string) {
			ProvisionCommand(container)
		},
	}

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Write config and restart services",
		Run: func(cmd *cobra.Command, args []string) {
			DeployCommand(container)
		},
	}

	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart services",
		Run: func(cmd *cobra.Command, args []string) {
			RestartCommand(container)
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop services",
		Run: func(cmd *cobra.Command, args []string) {
			StopCommand(container)
		},
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start services",
		Run: func(cmd *cobra.Command, args []string) {
			StartCommand(container)
		},
	}

	rootCmd.AddCommand(
		checkCmd,
		provisionCmd,
		deployCmd,
		restartCmd,
		startCmd,
		stopCmd,
	)

	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func runTask(container *Container, taskName string) {
	err := container.GetTask(taskName).Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
