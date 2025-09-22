package task

import (
	"ampctl/config"
	"ampctl/util"
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type ApacheInstallTask struct {
}

func (t *ApacheInstallTask) Run() error {
	fmt.Print("Check if apache is installed: ")
	if !util.IsPackageInstalled("httpd") {
		fmt.Println("Not installed yet, so we install it")
		err := util.InstallPackage("httpd")
		if err != nil {
			return fmt.Errorf("Error installing apache (httpd)")
		}
	} else {
		fmt.Println("Already installed")
	}

	return nil
}

type ApacheConfigWriteTask struct {
	Config *config.Config
}

func (t *ApacheConfigWriteTask) Run() error {
	fmt.Println("Write apache config")
	file := fmt.Sprintf("/opt/homebrew/etc/httpd/httpd.conf")
	err := writeConfig(file, t.Config.Apache)
	if err != nil {
		return fmt.Errorf("Error installing apache (httpd)")
	}

	file = fmt.Sprintf("/opt/homebrew/etc/httpd/extra/httpd-ssl.conf")
	err = writeSslConfig(file, t.Config.Apache)
	if err != nil {
		return fmt.Errorf("Error writing httpd-ssl.conf")
	}

	fmt.Println("Write apache hosts")
	file = fmt.Sprintf("/opt/homebrew/etc/httpd/extra/httpd-vhosts.conf")
	err = writeApacheHosts(file, t.Config.Hosts, t.Config.Apache)
	if err != nil {
		return fmt.Errorf("Error writing httpd-vhosts.conf")
	}

	fmt.Println("Test config")
	cmd := exec.Command("/opt/homebrew/bin/apachectl", "configtest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type ApacheStartTask struct {
}

func (t *ApacheStartTask) Run() error {
	return util.BrewStartService("httpd")
}

type ApacheRestartTask struct {
}

func (t *ApacheRestartTask) Run() error {
	return util.BrewRestartService("httpd")
}

type ApacheStopTask struct {
}

func (t *ApacheStopTask) Run() error {
	return util.BrewStopService("httpd")
}

//go:embed templates/*
var templatesFS embed.FS

func writeApacheHosts(filepath string, hosts []config.Host, apache config.Apache) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/apache/hosts.tmpl")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	for _, host := range hosts {
		data := map[string]any{
			"Port":       apache.HttpPort,
			"SslPort":    apache.HttpsPort,
			"Host":       host.Host,
			"Path":       host.Path,
			"Ssl":        host.Ssl,
			"PhpVersion": host.Version,
		}

		err = tmpl.Execute(buf, data)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(filepath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func writeConfig(filepath string, config config.Apache) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/apache/httpd-config.tmpl")
	if err != nil {
		panic(err)
	}

	data := map[string]string{
		"Workspace": config.Workspace,
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}

	err = util.BlockInFile(filepath, buf.String())
	if err != nil {
		return err
	}

	err = setApacheConfig(filepath, "Listen", config.HttpPort)
	if err != nil {
		return err
	}

	err = loadApacheModule(filepath, "rewrite_module", "lib/httpd/modules/mod_rewrite.so")
	if err != nil {
		return err
	}

	err = loadApacheModule(filepath, "proxy_module", "lib/httpd/modules/mod_proxy.so")
	if err != nil {
		return err
	}

	err = loadApacheModule(filepath, "ssl_module", "lib/httpd/modules/mod_ssl.so")
	if err != nil {
		return err
	}

	err = includeApacheConfig(filepath, "/opt/homebrew/etc/httpd/extra/httpd-vhosts.conf")
	if err != nil {
		return err
	}

	err = includeApacheConfig(filepath, "/opt/homebrew/etc/httpd/extra/httpd-ssl.conf")
	if err != nil {
		return err
	}

	return nil
}

func writeSslConfig(filepath string, apache config.Apache) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/apache/httpd-ssl.tmpl")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	data := map[string]any{
		"SslPort":               apache.HttpsPort,
		"SslCertificateFile":    apache.SslCertificateFile,
		"SslCertificateKeyFile": apache.SslCertificateKeyFile,
	}

	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadApacheModule(filepath string, name string, modulePath string) error {
	directive := fmt.Sprintf("LoadModule %s %s", name, modulePath)

	pattern := fmt.Sprintf(`^#?LoadModule\s+%s`, name)

	err := util.LineInFile(filepath, pattern, directive)
	if err != nil {
		return err
	}

	return nil
}

func includeApacheConfig(filepath string, configPath string) error {
	file := fmt.Sprintf("Include %s", configPath)

	pattern := fmt.Sprintf(`^#?Include\s+%s`, configPath)

	err := util.LineInFile(filepath, pattern, file)
	if err != nil {
		return err
	}

	return nil
}

func setApacheConfig(filepath string, name string, value string) error {
	// Read file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	lines := bytes.Split(content, []byte("\n"))
	found := false
	directive := []byte(name)

	for i, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if bytes.HasPrefix(trimmed, directive) {
			// Replace existing directive with new value
			lines[i] = []byte(fmt.Sprintf("%s %s", name, value))
			found = true
			break
		}
	}

	if !found {
		// Append directive if not found
		lines = append(lines, []byte(fmt.Sprintf("%s %s", name, value)))
	}

	// Write file back
	err = os.WriteFile(filepath, bytes.Join(lines, []byte("\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}
