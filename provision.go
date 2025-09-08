package main

import (
	"ampctl/apache"
	"ampctl/brew"
	"ampctl/config"
	"fmt"
	"sort"
)

func ProvisionCommand() {
	fmt.Printf("Start provision\n")

	success, cfg := ProvisionLoadConfigs()
	if !success {
		return
	}

	if !ProvisionHomebrew() {
		return
	}

	if !ProvisionApache(cfg) {
		return
	}

	if !ProvisionShivammathur() {
		return
	}

	versions := ProvisionGetPhpVersions(cfg.Php.Versions)
	for _, version := range versions {
		if !cfg.Php.Versions[version].Enabled {
			continue
		}
		if !ProvisionInstallPhp(version, cfg.Php.Versions[version]) {
			return
		}
	}
}

func ProvisionHomebrew() bool {
	if path, ok := brew.CheckPath(); ok {
		fmt.Printf("Homebrew is installed at: %s\n", path)
		return true
	} else {
		fmt.Println("Homebrew is NOT installed. Please install first before you go on")
		return false
	}
}

func ProvisionShivammathur() bool {

	fmt.Print("Check if shivammathur/php is added: ")
	if !brew.IsTapInstalled("shivammathur/php") {
		fmt.Println("Not added yet, so we add")
		err := brew.AddTap("shivammathur/php")
		if err != nil {
			fmt.Println("Error installing shivammathur/php")
			return false
		}
	} else {
		fmt.Println("Already added")
	}
	return true
}

func ProvisionLoadConfigs() (bool, *config.Config) {
	fmt.Println("Load config")
	cfg := &config.Config{}
	err := cfg.ImportConfigFileFromHomeDir()

	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	return true, cfg
}

func ProvisionGetPhpVersions(config map[string]config.PhpVersion) []string {
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func ProvisionInstallPhp(version string, config config.PhpVersion) bool {
	fmt.Print(fmt.Sprintf("Check PHP version %s: ", version))

	packageName := fmt.Sprintf("shivammathur/php/php@%s", version)

	if brew.IsPackageInstalled(packageName) {
		fmt.Println("is already installed")
	} else {
		fmt.Println("not installed yet, so we install")
		err := brew.InstallPackage(packageName)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	return true
}

func ProvisionApache(config *config.Config) bool {
	fmt.Print("Check if apache is installed: ")
	if !brew.IsPackageInstalled("httpd") {
		fmt.Println("Not installed yet, so we install it")
		err := brew.InstallPackage("httpd")
		if err != nil {
			fmt.Println("Error installing apache (httpd)")
			return false
		}
	} else {
		fmt.Println("Already installed")
	}
	return ProvisionApacheConfig(config)
}

func ProvisionApacheConfig(config *config.Config) bool {
	fmt.Println("Write apache config")
	file := fmt.Sprintf("/opt/homebrew/etc/httpd/httpd.conf")
	err := apache.WriteConfig(file, config.Apache)
	if err != nil {
		fmt.Println("Error installing apache (httpd)")
		return false
	}

	file = fmt.Sprintf("/opt/homebrew/etc/httpd/extra/ampctl-hosts.conf")
	err = apache.WriteHosts(file, config.Hosts, config.Apache)
	if err != nil {
		fmt.Println("Error installing apache (httpd)")
		return false
	}

	return true
}
