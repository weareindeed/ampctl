package brew

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	cmd := exec.Command("brew", "list", "--formula", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If `brew list` fails, the package is likely not installed
		return false
	}

	// remove path / tap
	parts := strings.Split(pkg, "/")
	pkg = parts[len(parts)-1]

	// Sometimes brew lists other info, so check if the package name appears in output
	return strings.Contains(string(output), pkg)
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
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	fmt.Println(string(output))
	return true
}
