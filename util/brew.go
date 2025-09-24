package util

import (
	"fmt"
	"os"
	"os/exec"
)

// CheckPath checks if Homebrew is installed and returns its path.
// Returns (path, true) if found, ("", false) otherwise.
func CheckPath() (string, bool) {
	path, err := exec.LookPath("brew")
	if err != nil {
		return "", false
	}
	return path, true
}

// IsPackageInstalled checks if a Homebrew package is installed.
func IsPackageInstalled(pkg string) bool {
	cmd := exec.Command("brew", "ls", "--versions", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If `brew list` fails, the package is likely not installed
		return false
	}
	return len(output) > 0
}

// InstallPackage installs a Homebrew package if it is not already installed.
// Returns an error if the installation fails.
func InstallPackage(pkg string) error {
	if IsPackageInstalled(pkg) {
		return nil
	}
	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// AddTap adds a Homebrew tap.
func AddTap(tap string) error {
	cmd := exec.Command("brew", "tap", tap)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// IsTapInstalled checks if a Homebrew tap is installed.
func IsTapInstalled(tap string) bool {
	cmd := exec.Command("brew", "tap")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return true
}

func BrewStartService(formula string) error {
	cmd := exec.Command("brew", "services", "start", formula)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BrewStopService(formula string) error {
	cmd := exec.Command("brew", "services", "stop", formula)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BrewRestartService(formula string) error {
	cmd := exec.Command("brew", "services", "restart", formula)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
