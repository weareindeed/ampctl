package task

import (
	"ampctl/config"
	"ampctl/util"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"gopkg.in/ini.v1"
	"howett.net/plist"
)

type DatabaseInstallTask struct {
	Config *config.Config
}

func (t *DatabaseInstallTask) Run() error {
	versions := getDatabaseVersions(t.Config.Database.Versions)
	for _, version := range versions {
		err := t.installSingleVersion(version, t.Config.Database.Versions[version])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *DatabaseInstallTask) installSingleVersion(packageName string, version config.DatabaseVersion) error {
	fmt.Print(fmt.Sprintf("Check Database version %s: ", packageName))

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

type DatabaseConfigWriteTask struct {
	Config *config.Config
}

func (t *DatabaseConfigWriteTask) Run() error {
	versions := getDatabaseVersions(t.Config.Database.Versions)
	for _, version := range versions {
		err := t.configSingleVersion(version, t.Config.Database.Versions[version])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *DatabaseConfigWriteTask) configSingleVersion(packageName string, version config.DatabaseVersion) error {
	// Run dir
	runDir := path.Join("/opt/homebrew/var/run/mysql", packageName)
	if _, err := os.Stat(runDir); os.IsNotExist(err) {
		err := os.MkdirAll(runDir, 0755)
		if err != nil {
			return err
		}
	}

	// Data dir
	dataDir := path.Join("/opt/homebrew/var/database", packageName)
	if version.DataDir != "" {
		dataDir = version.DataDir
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			return err
		}
		err = t.initDatabaseDir(dataDir, packageName)
		if err != nil {
			return err
		}
	}

	// Config Dir
	configDir := path.Join("/opt/homebrew/etc/mysql", packageName)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(runDir, 0755)
		if err != nil {
			return err
		}
	}

	cnfFile := path.Join(configDir, "my.cnf")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(runDir, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(cnfFile)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	err = t.writeConfig(cnfFile, dataDir, packageName, version)
	if err != nil {
		return err
	}

	// Daemon
	err = t.updateMacDaemon(packageName, version)
	if err != nil {
		return err
	}

	return nil
}

func (t *DatabaseConfigWriteTask) initDatabaseDir(dataDir string, packageName string) error {
	binPath := path.Join("/opt/homebrew/opt", packageName, "bin/mysql_install_db")

	cmd := exec.Command(
		binPath,
		fmt.Sprintf("--datadir=%s", dataDir),
		"--auth-root-authentication-method=normal",
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to initialize database dir for %s: %v\n", packageName, err)
	}

	return nil
}

func (t *DatabaseConfigWriteTask) writeConfig(cnfFile string, dataDir string, packageName string, version config.DatabaseVersion) error {
	cfg, err := ini.Load(cnfFile)
	if err != nil {
		panic(err)
	}

	// Network
	cfg.Section("mysqld").Key("skip_networking").MustInt(0)
	cfg.Section("mysqld").Key("bind-address").MustString("127.0.0.1")
	cfg.Section("mysqld").Key("port").MustString(version.Port)

	// Data dir
	cfg.Section("mysqld").Key("datadir").MustString(dataDir)

	// Logging
	generalLogFile := path.Join("/opt/homebrew/var/log/mysql", packageName, "mysql.log")
	cfg.Section("mysqld").Key("general_log_file").MustString(generalLogFile)

	generalLog := "0"
	if version.GeneralLog {
		generalLog = "0"
	}
	cfg.Section("mysqld").Key("general_log").MustString(generalLog)

	logErrorFile := path.Join("/opt/homebrew/var/log/mysql", packageName, "error.log")
	cfg.Section("mysqld").Key("log_error").MustString(logErrorFile)

	slowLogFile := path.Join("/opt/homebrew/var/log/mysql", packageName, "slow.log")
	cfg.Section("mysqld").Key("slow_query_log_file").MustString(slowLogFile)
	cfg.Section("mysqld").Key("slow_query_log").MustInt(1)

	// Resources and Timeouts
	cfg.Section("mysqld").Key("long_query_time").MustInt(1)
	cfg.Section("mysqld").Key("interactive_timeout").MustInt(300)
	cfg.Section("mysqld").Key("wait_timeout").MustInt(300)
	cfg.Section("mysqld").Key("max_allowed_packet").MustString("256M")
	cfg.Section("mysqld").Key("table_open_cache").MustInt(250)

	// MariaDB
	if strings.Contains(packageName, "mariadb") {
		cfg.Section("mariadb").Key("lower_case_table_names").MustInt(2)
	}

	// MysqlDB
	if strings.Contains(packageName, "mysql") {
		cfg.Section("mysqld").Key("default_authentication_plugin").MustInt(2)
	}

	var buf bytes.Buffer
	_, err = cfg.WriteTo(&buf)
	if err != nil {
		return err
	}

	err = os.WriteFile(cnfFile, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *DatabaseConfigWriteTask) updateMacDaemon(packageName string, version config.DatabaseVersion) error {
	plistPath := path.Join(
		"/opt/homebrew/opt",
		packageName,
		fmt.Sprintf("homebrew.mxcl.%s.plist", packageName),
	)

	// Read plist
	data, err := os.ReadFile(plistPath)
	if err != nil {
		return fmt.Errorf("could not read plist: %w", err)
	}

	// Decode into a generic map
	var raw map[string]any
	_, err = plist.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("could not unmarshal plist: %w", err)
	}

	// Replace ProgramArguments
	raw["ProgramArguments"] = []string{
		fmt.Sprintf("/opt/homebrew/opt/%s/bin/mysqld_safe", packageName),
		fmt.Sprintf("--defaults-file=/opt/homebrew/etc/mysql/%s/my.cnf", packageName),
	}

	// Marshal back
	out, err := plist.MarshalIndent(raw, plist.XMLFormat, "  ")
	if err != nil {
		return fmt.Errorf("could not marshal plist: %w", err)
	}

	// Write back
	if err := os.WriteFile(plistPath, out, 0644); err != nil {
		return fmt.Errorf("could not write plist: %w", err)
	}

	return nil
}

type DatabaseStartTask struct {
	Config *config.Config
}

func (t *DatabaseStartTask) Run() error {
	for version, _ := range t.Config.Database.Versions {
		if t.Config.Database.Versions[version].Enabled {
			err := util.BrewStartService(version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type DatabaseRestartTask struct {
	Config *config.Config
}

func (t *DatabaseRestartTask) Run() error {
	for version, _ := range t.Config.Database.Versions {
		if t.Config.Database.Versions[version].Enabled {
			err := util.BrewRestartService(version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type DatabaseStopTask struct {
	Config *config.Config
}

func (t *DatabaseStopTask) Run() error {
	for version, _ := range t.Config.Database.Versions {
		if t.Config.Database.Versions[version].Enabled {
			err := util.BrewStopService(version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getDatabaseVersions(config map[string]config.DatabaseVersion) []string {
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
