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

	rootCmd.AddCommand(
		checkCmd,
		provisionCmd,
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
	}
}
