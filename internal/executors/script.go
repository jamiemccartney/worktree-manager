package executors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"worktree-manager/internal/consts"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
)

// ScriptExecutor defines the interface for executing scripts
type ScriptExecutor interface {
	Execute(ctx *ScriptExecutionContext) error
}

// ScriptExecutionContext contains all parameters needed for script execution
type ScriptExecutionContext struct {
	ScriptPath   string
	Repo         *state.Repo
	WorktreePath string
	WorkingDir   string
	ProgressMsg  string
}

// BashScriptExecutor implements ScriptExecutor for bash scripts
type BashScriptExecutor struct{}

func NewBashScriptExecutor() *BashScriptExecutor {
	return &BashScriptExecutor{}
}

func (e *BashScriptExecutor) Execute(ctx *ScriptExecutionContext) error {
	if ctx.ScriptPath == "" {
		return nil
	}

	resolvedPath, err := resolveScriptPath(ctx.ScriptPath, ctx.Repo.Dir)
	if err != nil {
		return err
	}

	env := buildScriptEnvironment(ctx.Repo, ctx.WorktreePath)
	cmd := createScriptCommand(resolvedPath, env, ctx.WorkingDir)

	if ctx.ProgressMsg != "" {
		output.Progress(ctx.ProgressMsg, resolvedPath)
	}

	return cmd.Run()
}

func resolveScriptPath(scriptPath, repoDir string) (string, error) {
	if !filepath.IsAbs(scriptPath) {
		scriptPath = filepath.Join(repoDir, scriptPath)
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("script does not exist: %s", scriptPath)
	}

	return scriptPath, nil
}

func buildScriptEnvironment(repo *state.Repo, worktreePath string) []string {
	env := os.Environ()
	envVars := consts.GetEnvironmentVariables()
	env = append(env, fmt.Sprintf("%s=%s", envVars.RepoAlias.Name, repo.Alias))
	env = append(env, fmt.Sprintf("%s=%s", envVars.RepoDir.Name, repo.Dir))
	env = append(env, fmt.Sprintf("%s=%s", envVars.WorktreePath.Name, worktreePath))
	return env
}

func createScriptCommand(scriptPath string, env []string, workingDir string) *exec.Cmd {
	cmd := exec.Command("bash", scriptPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	return cmd
}
