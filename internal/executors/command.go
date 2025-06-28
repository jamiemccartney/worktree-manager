package executors

import (
	"fmt"
	"os/exec"

	"worktree-manager/internal/output"
)

// CommandExecutor defines the interface for executing system commands
type CommandExecutor interface {
	Execute(ctx *CommandExecutionContext) error
}

// CommandExecutionContext contains all parameters needed for command execution
type CommandExecutionContext struct {
	Command     string
	Args        []string
	WorkingDir  string
	Env         []string
	ProgressMsg string
	ShowOutput  bool
}

// SystemCommandExecutor implements CommandExecutor for system commands
type SystemCommandExecutor struct{}

func NewSystemCommandExecutor() *SystemCommandExecutor {
	return &SystemCommandExecutor{}
}

func (e *SystemCommandExecutor) Execute(ctx *CommandExecutionContext) error {
	if ctx.Command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	cmd := exec.Command(ctx.Command, ctx.Args...)

	if ctx.WorkingDir != "" {
		cmd.Dir = ctx.WorkingDir
	}

	if len(ctx.Env) > 0 {
		cmd.Env = ctx.Env
	}

	if ctx.ProgressMsg != "" {
		output.Progress(ctx.ProgressMsg)
	}

	if ctx.ShowOutput {
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("command failed: %v\nOutput: %s", err, string(output))
		}
		return nil
	}

	return cmd.Run()
}
