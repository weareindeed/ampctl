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

func LoadModule(name string) {

}
