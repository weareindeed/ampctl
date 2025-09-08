package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "ampctl"}

	var checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Check settings",
		Run: func(cmd *cobra.Command, args []string) {
			CheckCommand()
		},
	}

	var provisionCmd = &cobra.Command{
		Use:   "provision",
		Short: "Run provision",
		Run: func(cmd *cobra.Command, args []string) {
			ProvisionCommand()
		},
	}

	rootCmd.AddCommand(
		checkCmd,
		provisionCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
