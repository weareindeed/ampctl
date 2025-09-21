package task

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleLineSetApacheConfig(t *testing.T) {

	tmp := filepath.Join(os.TempDir(), "test_blockinfile.txt")
	defer os.Remove(tmp)

	content := "Listen 8080"

	err := os.WriteFile(tmp, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	err = setApacheConfig(tmp, "Listen", "80")
	if err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	bytes, _ := os.ReadFile(tmp)
	content = string(bytes)

	assert.Equal(t, "Listen 80", content)
}

func TestMultipleLineSetApacheConfig(t *testing.T) {

	tmp := filepath.Join(os.TempDir(), "test_blockinfile.txt")
	defer os.Remove(tmp)

	content := "\n" +
		"Listen 8080" +
		"\n"

	err := os.WriteFile(tmp, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	err = setApacheConfig(tmp, "Listen", "80")
	if err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	bytes, _ := os.ReadFile(tmp)
	content = string(bytes)

	assert.Contains(t, content, "Listen 80")
}
