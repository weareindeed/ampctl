package hosts

import (
	"ampctl/util"
	"strings"
)

func WriteHosts(hosts []string) {
	var sb strings.Builder

	for _, host := range hosts {
		sb.WriteString("127.0.0.1 " + host + "\n")
		sb.WriteString("::1 " + host + "\n")
	}

	content := sb.String()

	util.BlockInFile("/etc/hosts", content)
}
