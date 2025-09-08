package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Host struct {
	Host    string `yaml:"host"`
	Path    string `yaml:"path"`
	Version string `yaml:"version"`
	Ssl     bool   `yaml:"ssl"`
}

type Php struct {
	Default   string                `yaml:"default"`
	Composer1 string                `yaml:"composer1"`
	Composer2 string                `yaml:"composer2"`
	Versions  map[string]PhpVersion `yaml:"versions"`
}

type Apache struct {
	Workspace string `yaml:"workspace"`
}

type PhpVersion struct {
	Enabled bool `yaml:"enabled"`
}

type Config struct {
	Hosts  []Host `yaml:"hosts"`
	Php    Php    `yaml:"php"`
	Apache Apache `yaml:"apache"`
}

// ImportFile reads the hosts file from the given path
func (c *Config) ImportFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read hosts file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse yaml: %w", err)
	}

	return nil
}

// ImportHostsFileFromHomeDir reads the hosts from path . ~/.ampctl/hosts.yaml
func (c *Config) ImportHostsFileFromHomeDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	hostsFile := filepath.Join(homeDir, ".ampctl", "hosts.yaml")
	_, statErr := os.Stat(hostsFile)
	if os.IsNotExist(statErr) {
		// Try to create the file (and parent dir if needed)
		if err := os.MkdirAll(filepath.Dir(hostsFile), 0755); err != nil {
			return fmt.Errorf("failed to create directory for hosts file: %w", err)
		}
		file, err := os.Create(hostsFile)
		if err != nil {
			return fmt.Errorf("failed to create hosts file: %w", err)
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	if err := c.ImportFile(hostsFile); err != nil {
		return fmt.Errorf("failed to load hosts file: %w", err)
	}

	return nil
}

// ImportConfigFileFromHomeDir reads the hosts from path . ~/.ampctl/hosts.yaml
func (c *Config) ImportConfigFileFromHomeDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configFile := filepath.Join(homeDir, ".ampctl", "config.yaml")
	_, statErr := os.Stat(configFile)
	if os.IsNotExist(statErr) {
		// Try to create the file (and parent dir if needed)
		if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
			return fmt.Errorf("failed to create directory for config file: %w", err)
		}
		file, err := os.Create(configFile)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	if err := c.ImportFile(configFile); err != nil {
		return fmt.Errorf("failed to load hosts file: %w", err)
	}

	return nil
}
