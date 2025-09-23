package task

import (
	"ampctl/config"
	"ampctl/util"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

type PhpStartTask struct {
	Config *config.Config
}

func (t *PhpStartTask) Run() error {
	versions := getPhpVersions(t.Config.Php.Versions)
	for _, version := range versions {
		err := util.BrewStartService("php@" + version)
		if err != nil {
			return err
		}
	}
	return nil
}

type PhpRestartTask struct {
	Config *config.Config
}

func (t *PhpRestartTask) Run() error {
	versions := getPhpVersions(t.Config.Php.Versions)
	for _, version := range versions {
		err := util.BrewRestartService("php@" + version)
		if err != nil {
			return err
		}
	}
	return nil
}

type PhpStopTask struct {
	Config *config.Config
}

func (t *PhpStopTask) Run() error {
	versions := getPhpVersions(t.Config.Php.Versions)
	for _, version := range versions {
		err := util.BrewStopService("php@" + version)
		if err != nil {
			return err
		}
	}
	return nil
}

type PhpInstallTask struct {
	Config *config.Config
}

func (t *PhpInstallTask) Run() error {
	versions := getPhpVersions(t.Config.Php.Versions)
	for _, version := range versions {
		err := t.installSingleVersion(version, t.Config.Php.Versions[version])
		if err != nil {
			return err
		}

		dir := "/opt/homebrew/php-bin"

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}

		err = t.installExecutables(version, t.Config.Php.Versions[version])
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

func (t *PhpInstallTask) installExecutables(version string, config config.PhpVersion) error {
	dir := path.Join("/opt/homebrew/php-bin", version)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("brew", "ls", "--versions", "php@"+version)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	binPath := path.Join("/opt/homebrew/Cellar", parts[0], parts[1], "bin/php")

	link := path.Join(dir, "php")

	if _, err := os.Lstat(link); err == nil {
		err = os.Remove(link)
		if err != nil {
			return err
		}
	}

	err = os.Symlink(binPath, link)
	if err != nil {
		return err
	}

	return nil
}

func getPhpVersions(config map[string]config.PhpVersion) []string {
	keys := make([]string, 0, len(config))
	for k, version := range config {
		if !version.Enabled {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

type PhpWriteConfigTask struct {
	Config *config.Config
}

func (t *PhpWriteConfigTask) Run() error {
	versions := getPhpVersions(t.Config.Php.Versions)
	for _, version := range versions {
		fmt.Println("Write config for php " + version)
		err := t.writeFpmConfig(version, t.Config.Php.Versions[version])
		if err != nil {
			return err
		}

		err = t.writeIniConfig(version, t.Config.Php.Versions[version])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *PhpWriteConfigTask) writeFpmConfig(name string, version config.PhpVersion) error {
	wwwConf := path.Join("/opt/homebrew/etc/php", name, "php-fpm.d/www.conf")
	fpmConf := path.Join("/opt/homebrew/etc/php", name, "php-fpm.conf")

	err := setPhpConfig(wwwConf, "listen.mode", "0666")
	if err != nil {
		return err
	}

	sockDir := path.Join("/opt/homebrew/var/run/php", name)
	if _, err := os.Stat(sockDir); os.IsNotExist(err) {
		if err := os.MkdirAll(sockDir, 0755); err != nil {
			return err
		}
	}

	sockFile := path.Join("/opt/homebrew/var/run/php", name, "php-fpm.sock")
	err = setPhpConfig(wwwConf, "listen", sockFile)
	if err != nil {
		return err
	}

	logDir := path.Join("/opt/homebrew/var/log/php", name)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	slowLogPath := path.Join("/opt/homebrew/var/log/php", name, "slow.log")
	err = setPhpConfig(wwwConf, "slowlog", slowLogPath)
	if err != nil {
		return err
	}

	errorLogPath := path.Join("/opt/homebrew/var/log/php", name, "error.log")
	err = setPhpConfig(fpmConf, "error_log", errorLogPath)
	if err != nil {
		return err
	}

	err = os.Remove("/opt/homebrew/var/log/php-fpm.log")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.Remove("/opt/homebrew/var/log/php-fpm.log.default")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (t *PhpWriteConfigTask) writeIniConfig(name string, version config.PhpVersion) error {
	//iniConf := path.Join("/opt/homebrew/etc/php", name, "php.ini")
	return nil
}

func setPhpConfig(filename string, key string, value string) error {
	newLine := fmt.Sprintf("%s = %s", key, value)
	regex := fmt.Sprintf("^;?%s\\s=.+$", key)
	return util.LineInFile(filename, regex, newLine)
}
