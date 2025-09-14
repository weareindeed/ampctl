package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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

func run() {
	// Check if we are already running as sudo
	if os.Geteuid() == 0 {
		fmt.Println("Already running as root!")
		return
	}

	// get current binary path
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Error determining executable path: %v", err)
	}

	// inherit arguments of current process
	args := os.Args[1:]
	cmdArgs := append([]string{exe}, args...)

	// start new process with sudo
	cmd := exec.Command("sudo", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error starting with sudo: %v", err)
	}
}
