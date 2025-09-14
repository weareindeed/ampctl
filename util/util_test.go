package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBlockInFile(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "test_blockinfile.txt")
	defer os.Remove(tmp)

	// initial write
	if err := BlockInFile(tmp, "hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(tmp)
	if !strings.Contains(string(content), "hello") {
		t.Errorf("expected content to contain 'hello', got: %s", string(content))
	}

	// update block
	if err := BlockInFile(tmp, "world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ = os.ReadFile(tmp)
	if !strings.Contains(string(content), "world") {
		t.Errorf("expected content to contain 'world', got: %s", string(content))
	}
}

func TestLineInFile(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "test_lineinfile.txt")
	defer os.Remove(tmp)

	// add new line
	if err := LineInFile(tmp, "^foo=", "foo=bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(tmp)
	if !strings.Contains(string(content), "foo=bar") {
		t.Errorf("expected 'foo=bar', got: %s", string(content))
	}

	// replace line
	if err := LineInFile(tmp, "^foo=", "foo=baz"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ = os.ReadFile(tmp)
	if strings.Contains(string(content), "foo=bar") {
		t.Errorf("expected replacement, still contains 'foo=bar'")
	}
	if !strings.Contains(string(content), "foo=baz") {
		t.Errorf("expected 'foo=baz', got: %s", string(content))
	}
}

func TestNotSudoCommand(t *testing.T) {
	cmd := NotSudoCommand("echo", "hello")
	if cmd.Path != "echo" && !strings.Contains(cmd.Path, "echo") {
		t.Errorf("expected echo command, got %s", cmd.Path)
	}

	// Run the command to ensure it executes
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to run command: %v", err)
	}
	if strings.TrimSpace(string(out)) != "hello" {
		t.Errorf("expected 'hello', got %q", string(out))
	}
}
