package task

import (
	"ampctl/config"
	"ampctl/util"
	"fmt"
	"strings"
)

type HostsWriteTask struct {
	Config *config.Config
}

func (t *HostsWriteTask) Run() error {
	fmt.Println("Write hosts")
	err := writeHosts(t.Config.Hosts)
	if err != nil {
		return err
	}
	return nil
}

func writeHosts(hosts []config.Host) error {
	var sb strings.Builder

	for _, host := range hosts {
		sb.WriteString("127.0.0.1 " + host.Host + "\n")
		sb.WriteString("::1 " + host.Host + "\n")
	}

	content := sb.String()

	return util.BlockInFile("/etc/hosts", content)
}
