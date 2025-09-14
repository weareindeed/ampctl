package hosts

import (
	"ampctl/config"
	"ampctl/util"
	"strings"
)

func WriteHosts(hosts []config.Host) error {
	var sb strings.Builder

	for _, host := range hosts {
		sb.WriteString("127.0.0.1 " + host.Host + "\n")
		sb.WriteString("::1 " + host.Host + "\n")
	}

	content := sb.String()

	return util.BlockInFile("/etc/hosts", content)
}
