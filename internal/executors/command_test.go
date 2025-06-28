package executors

import (
	"testing"
)

func TestSystemCommandExecutor_Execute_EmptyCommand(t *testing.T) {
	executor := NewSystemCommandExecutor()
	ctx := &CommandExecutionContext{
		Command: "",
	}

	err := executor.Execute(ctx)
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}

	expectedMsg := "command cannot be empty"

	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}
