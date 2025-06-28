package fileops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"worktree-manager/internal/consts"
)

func TestReadWriteJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	// Test data
	testData := map[string]interface{}{
		"name":  "test",
		"count": 42,
	}

	// Test write
	err := WriteJSONFile(testFile, testData)
	if err != nil {
		t.Fatalf("WriteJSONFile failed: %v", err)
	}

	// Test read
	var readData map[string]interface{}
	err = ReadJSONFile(testFile, &readData)
	if err != nil {
		t.Fatalf("ReadJSONFile failed: %v", err)
	}

	if readData["name"] != "test" {
		t.Errorf("Expected name=test, got %v", readData["name"])
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.txt")
	nonExistingFile := filepath.Join(tmpDir, "nonexisting.txt")

	// Create a file
	os.WriteFile(existingFile, []byte("test"), 0644)

	if !FileExists(existingFile) {
		t.Error("Expected FileExists to return true for existing file")
	}

	if FileExists(nonExistingFile) {
		t.Error("Expected FileExists to return false for non-existing file")
	}
}

func TestExpandEnvVars(t *testing.T) {
	// Test HOME expansion
	result := ExpandEnvVars("$HOME/test")
	if !strings.Contains(result, "/test") {
		t.Errorf("Expected HOME to be expanded, got %s", result)
	}

	// Test no expansion needed
	result = ExpandEnvVars("/absolute/path")
	if result != "/absolute/path" {
		t.Errorf("Expected no change, got %s", result)
	}
}

func TestGetWorktreeManagerDir(t *testing.T) {
	dir := consts.GetDirectoryPaths().WorktreeManagerDir
	if !strings.HasSuffix(dir, ".worktree-manager") {
		t.Errorf("Expected directory to end with .worktree-manager, got %s", dir)
	}
}

func TestCreateExecutableScript(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test-script.sh")
	content := "#!/bin/bash\necho 'test'"

	err := CreateExecutableScript(scriptPath, content)
	if err != nil {
		t.Fatalf("CreateExecutableScript failed: %v", err)
	}

	// Check file exists
	if !FileExists(scriptPath) {
		t.Error("Script file was not created")
	}

	// Check file is executable
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("Failed to stat script file: %v", err)
	}

	mode := info.Mode()
	if mode&0111 == 0 {
		t.Error("Script file is not executable")
	}
}
