package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// BlockInFile inserts or updates a block of text in a file.
func BlockInFile(filename, blockContent string) error {
	startMarker := fmt.Sprintf("# Beginn block managed by ampctl")
	endMarker := fmt.Sprintf("# End block managed by ampctl")

	// Read file content (if it exists)
	content, _ := os.ReadFile(filename)
	text := string(content)

	// Regex to find existing block
	re := regexp.MustCompile(fmt.Sprintf(`(?s)%s.*?%s`, regexp.QuoteMeta(startMarker), regexp.QuoteMeta(endMarker)))

	newBlock := fmt.Sprintf("%s\n%s\n%s", startMarker, blockContent, endMarker)

	if re.MatchString(text) {
		// Replace existing block
		text = re.ReplaceAllString(text, newBlock)
	} else {
		// Append new block
		if !strings.HasSuffix(text, "\n") {
			text += "\n"
		}
		text += newBlock + "\n"
	}

	// Write back to file
	return os.WriteFile(filename, []byte(text), 0644)
}

// LineInFile replaces the first line matching regexpPattern with 'line', or appends 'line' if no match.
func LineInFile(path string, regexpPattern string, line string) error {
	content, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	text := string(content)
	lines := strings.Split(text, "\n")
	re, err := regexp.Compile(regexpPattern)
	if err != nil {
		return err
	}
	found := false
	for i, l := range lines {
		if re.MatchString(l) {
			lines[i] = line
			found = true
			break
		}
	}
	if !found {
		// Remove possible trailing empty line before appending
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		lines = append(lines, line)
	}
	// Rejoin with \n (add final newline)
	result := strings.Join(lines, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return os.WriteFile(path, []byte(result), 0644)
}

func NotSudoCommand(name string, arg ...string) *exec.Cmd {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser == "" {
		cmd := exec.Command(name, arg...)
		return cmd
	}

	u, err := user.Lookup(sudoUser)
	if err != nil {
		log.Fatal(err)
	}

	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	// Prepare command
	cmd := exec.Command(name, arg...)

	// Run as specific user
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(uid),
		Gid: uint32(gid),
	}

	return cmd
}
