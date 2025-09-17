package task

import (
	"ampctl/config"
	"ampctl/util"
	"bytes"
	"embed"
	"fmt"
	"os"
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
	err := writeApacheHosts(file, t.Config.Hosts, t.Config.Apache)
	if err != nil {
		return fmt.Errorf("Error installing apache (httpd)")
	}

	file = fmt.Sprintf("/opt/homebrew/etc/httpd/extra/ampctl-hosts.conf")
	err = writeApacheHosts(file, t.Config.Hosts, t.Config.Apache)
	if err != nil {
		return fmt.Errorf("Error installing apache (httpd)")
	}

	return nil
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
			"Port":       80,
			"SslPort":    443,
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
	file := fmt.Sprintf(filepath)

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
		fmt.Println(err)
		return err
	}

	err = util.BlockInFile(file, buf.String())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func loadModule(filepath string, name string) {

}

func setConfig(filepath string, name string, value string) error {
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
