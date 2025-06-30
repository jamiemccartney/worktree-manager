package fileops

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"worktree-manager/internal/consts"
)

// ReadJSONFile reads and unmarshals a JSON file
func ReadJSONFile(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	return nil
}

// WriteJSONFile marshals and writes data to a JSON file
func WriteJSONFile(filePath string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}
	return nil
}

// ExpandEnvVars expands environment variables in a string
func ExpandEnvVars(s string) string {
	if strings.Contains(s, "$HOME") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		s = strings.ReplaceAll(s, "$HOME", homeDir)
	}
	return os.ExpandEnv(s)
}

// GetRepoScriptDir returns the script directory for a specific repo
func GetRepoScriptDir(repoAlias string) string {
	return filepath.Join(consts.GetDirectoryPaths().ScriptsDir, repoAlias)
}

// CreateExecutableScript creates a script file and makes it executable
func CreateExecutableScript(scriptPath, content string) error {
	if err := EnsureDir(filepath.Dir(scriptPath)); err != nil {
		return err
	}

	if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("failed to write script %s: %w", scriptPath, err)
	}

	return nil
}
