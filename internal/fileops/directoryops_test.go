package fileops

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWithDir_Success(t *testing.T) {
	// Get original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "withdir_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Resolve symlinks for consistent comparison
	tempDir, err = filepath.EvalSymlinks(tempDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for temp directory: %v", err)
	}

	var executedInDir string
	var functionCalled bool

	// Test successful execution
	err = WithDir(tempDir, func() error {
		functionCalled = true
		executedInDir, _ = os.Getwd()
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !functionCalled {
		t.Error("Function was not called")
	}

	// Verify function was executed in the correct directory
	if executedInDir != tempDir {
		t.Errorf("Function executed in wrong directory. Expected %s, got %s", tempDir, executedInDir)
	}

	// Verify we're back in the original directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory after test: %v", err)
	}

	if currentDir != originalDir {
		t.Errorf("Directory not restored. Expected %s, got %s", originalDir, currentDir)
	}
}

func TestWithDir_FunctionError(t *testing.T) {
	// Get original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "withdir_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	expectedError := errors.New("test function error")
	var functionCalled bool

	// Test function returning error
	err = WithDir(tempDir, func() error {
		functionCalled = true
		return expectedError
	})

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if !functionCalled {
		t.Error("Function was not called")
	}

	// Verify we're still back in the original directory even after error
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory after test: %v", err)
	}

	if currentDir != originalDir {
		t.Errorf("Directory not restored after error. Expected %s, got %s", originalDir, currentDir)
	}
}

func TestWithDir_InvalidDirectory(t *testing.T) {
	// Get original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	invalidDir := "/nonexistent/directory/that/should/not/exist"
	var functionCalled bool

	// Test with invalid directory
	err = WithDir(invalidDir, func() error {
		functionCalled = true
		return nil
	})

	if err == nil {
		t.Error("Expected error for invalid directory, got nil")
	}

	if functionCalled {
		t.Error("Function should not have been called for invalid directory")
	}

	// Verify error message contains expected information
	expectedMsg := "failed to change to directory"
	if err != nil && len(err.Error()) > 0 {
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Error message should contain '%s', got: %s", expectedMsg, err.Error())
		}
		if !contains(err.Error(), invalidDir) {
			t.Errorf("Error message should contain directory '%s', got: %s", invalidDir, err.Error())
		}
	}

	// Verify we're still in the original directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory after test: %v", err)
	}

	if currentDir != originalDir {
		t.Errorf("Directory changed unexpectedly. Expected %s, got %s", originalDir, currentDir)
	}
}

func TestWithDir_DirectoryChangeDuringExecution(t *testing.T) {
	// Get original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create temporary directories for testing
	tempDir1, err := os.MkdirTemp("", "withdir_test1")
	if err != nil {
		t.Fatalf("Failed to create temp directory 1: %v", err)
	}
	defer os.RemoveAll(tempDir1)

	tempDir2, err := os.MkdirTemp("", "withdir_test2")
	if err != nil {
		t.Fatalf("Failed to create temp directory 2: %v", err)
	}
	defer os.RemoveAll(tempDir2)

	// Resolve symlinks for consistent comparison
	tempDir1, err = filepath.EvalSymlinks(tempDir1)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for temp directory 1: %v", err)
	}

	tempDir2, err = filepath.EvalSymlinks(tempDir2)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for temp directory 2: %v", err)
	}

	var executedInDir, changedToDir string
	var functionCalled bool

	// Test that changing directory inside the function doesn't affect restoration
	err = WithDir(tempDir1, func() error {
		functionCalled = true
		executedInDir, _ = os.Getwd()

		// Change to another directory inside the function
		if err := os.Chdir(tempDir2); err != nil {
			return err
		}
		changedToDir, _ = os.Getwd()

		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !functionCalled {
		t.Error("Function was not called")
	}

	// Verify function started in the correct directory
	if executedInDir != tempDir1 {
		t.Errorf("Function started in wrong directory. Expected %s, got %s", tempDir1, executedInDir)
	}

	// Verify function changed to the second directory
	if changedToDir != tempDir2 {
		t.Errorf("Function didn't change to second directory. Expected %s, got %s", tempDir2, changedToDir)
	}

	// Verify we're restored to the original directory (not tempDir2)
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory after test: %v", err)
	}

	if currentDir != originalDir {
		t.Errorf("Directory not restored to original. Expected %s, got %s", originalDir, currentDir)
	}
}

func TestWithDir_EmptyDirectory(t *testing.T) {
	// Test with empty string directory
	var functionCalled bool

	err := WithDir("", func() error {
		functionCalled = true
		return nil
	})

	if err == nil {
		t.Error("Expected error for empty directory string, got nil")
	}

	if functionCalled {
		t.Error("Function should not have been called for empty directory")
	}
}

func TestWithDir_RelativePath(t *testing.T) {
	// Get original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "withdir_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Resolve symlinks for consistent comparison
	tempDir, err = filepath.EvalSymlinks(tempDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for temp directory: %v", err)
	}

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Change to temp directory first
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	var executedInDir string
	var functionCalled bool

	// Test with relative path
	err = WithDir("subdir", func() error {
		functionCalled = true
		executedInDir, _ = os.Getwd()
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error for relative path, got %v", err)
	}

	if !functionCalled {
		t.Error("Function was not called")
	}

	// Verify function was executed in the subdirectory
	if executedInDir != subDir {
		t.Errorf("Function executed in wrong directory. Expected %s, got %s", subDir, executedInDir)
	}

	// Verify we're back in the temp directory (where we started)
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory after test: %v", err)
	}

	if currentDir != tempDir {
		t.Errorf("Directory not restored. Expected %s, got %s", tempDir, currentDir)
	}

	// Restore original directory for cleanup
	if err := os.Chdir(originalDir); err != nil {
		t.Logf("Warning: Failed to restore original directory: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())
}
