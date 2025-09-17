package task

import (
	"ampctl/config"
	"ampctl/util"
	"fmt"
	"os"
	"regexp"
	"sort"
)

type PhpInstallTask struct {
	Config *config.Config
}

func (t *PhpInstallTask) Run() error {
	versions := t.getPhpVersions(t.Config.Php.Versions)
	for _, version := range versions {
		if !t.Config.Php.Versions[version].Enabled {
			continue
		}
		err := t.installSingleVersion(version, t.Config.Php.Versions[version])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *PhpInstallTask) installSingleVersion(name string, version config.PhpVersion) error {
	fmt.Print(fmt.Sprintf("Check PHP version %s: ", name))

	packageName := fmt.Sprintf("shivammathur/php/php@%s", name)

	if util.IsPackageInstalled(packageName) {
		fmt.Println("is already installed")
	} else {
		fmt.Println("not installed yet, so we install")
		err := util.InstallPackage(packageName)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (t *PhpInstallTask) getPhpVersions(config map[string]config.PhpVersion) []string {
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (t *PhpInstallTask) setPhpConfig(version string, key string, value string) error {
	file := fmt.Sprintf("/opt/homebrew/etc/php/%s/php.ini", version)

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(fmt.Sprintf("^;?%s.*$", key))
	newContent := re.ReplaceAllString(string(content), fmt.Sprintf("%s = %s", key, value))

	err = os.WriteFile(file, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
