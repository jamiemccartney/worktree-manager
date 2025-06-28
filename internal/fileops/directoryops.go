package fileops

import (
	"fmt"
	"os"

	"worktree-manager/internal/output"
)

// WithDir executes a function within a specified directory, then restores the original directory
func WithDir(dir string, fn func() error) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("failed to change to directory %s: %w", dir, err)
	}

	defer func() {
		if restoreErr := os.Chdir(originalDir); restoreErr != nil {
			output.Warning("Failed to restore original directory: %v", restoreErr)
		}
	}()

	return fn()
}
