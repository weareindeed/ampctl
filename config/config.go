package config

import (
	"fmt"
	"os"
	"path"
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
	Default          string                `yaml:"default"`
	ComposerVersions []string              `yaml:"composer_versions"`
	Versions         map[string]PhpVersion `yaml:"versions"`
}

func newPhp() Php {
	return Php{}
}

type Apache struct {
	Workspace                      string `yaml:"workspace"`
	HttpPort                       string `yaml:"http_port"`
	HttpsPort                      string `yaml:"https_port"`
	SslCertificateFile             string `yaml:"ssl_certificate_file"`
	SslCertificateKeyFile          string `yaml:"ssl_certificate_key_file"`
	SslCertificateCn               string `yaml:"ssl_certificate_cn"`
	SslCertificateCountry          string `yaml:"ssl_certificate_county"`
	SslCertificateLocality         string `yaml:"ssl_certificate_locality"`
	SslCertificateOrganization     string `yaml:"ssl_certificate_organization"`
	SslCertificateOrganizationUnit string `yaml:"ssl_certificate_organization_unit"`
	SslCertificateProvince         string `yaml:"ssl_certificate_province"`
}

func newApache() Apache {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err) // or handle error properly
	}

	return Apache{
		Workspace:             "/var/www",
		HttpPort:              "80",
		HttpsPort:             "443",
		SslCertificateFile:    path.Join(home, ".ampctl/CA.pem"),
		SslCertificateKeyFile: path.Join(home, ".ampctl/CA.key"),
	}
}

type PhpVersion struct {
	Enabled bool `yaml:"enabled"`
}

type Config struct {
	Hosts    []Host   `yaml:"hosts"`
	Php      Php      `yaml:"php"`
	Apache   Apache   `yaml:"apache"`
	Database Database `yaml:"database"`
}

func NewConfig() *Config {
	return &Config{
		Hosts:  []Host{},
		Php:    newPhp(),
		Apache: newApache(),
	}
}

type Database struct {
	Versions map[string]DatabaseVersion `yaml:"versions"`
}

type DatabaseVersion struct {
	Port    string `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
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

func (c *Config) LoadConfig() error {
	err := c.ImportConfigFileFromHomeDir()
	if err != nil {
		return err
	}
	return nil
}
