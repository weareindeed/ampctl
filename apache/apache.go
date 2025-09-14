package apache

import (
	"ampctl/config"
	"ampctl/util"
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed templates/*
var templatesFS embed.FS

func WriteHosts(filepath string, hosts []config.Host, apache config.Apache) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/hosts.tmpl")
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

func WriteConfig(filepath string, config config.Apache) error {
	file := fmt.Sprintf(filepath)

	tmpl, err := template.ParseFS(templatesFS, "templates/httpd-config.tmpl")
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

func LoadModule(filepath string, name string) {

}

func SetConfig(filepath string, name string, value string) error {
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
