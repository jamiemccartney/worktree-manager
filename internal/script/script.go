package script

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

// ExecuteScript runs a script with environment variables set
func ExecuteScript(scriptPath string, repo *config.Repo, worktreePath string) error {
	if scriptPath == "" {
		return nil // No script to execute
	}

	// Resolve script path
	if !filepath.IsAbs(scriptPath) {
		scriptPath = filepath.Join(repo.Dir, scriptPath)
	}

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script does not exist: %s", scriptPath)
	}

	// Prepare environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("REPO_ALIAS=%s", repo.Alias))
	env = append(env, fmt.Sprintf("REPO_DIR=%s", repo.Dir))
	env = append(env, fmt.Sprintf("WORKTREE_PATH=%s", worktreePath))
	env = append(env, fmt.Sprintf("POST_WORKTREE_ADD_SCRIPT=%s", repo.PostWorktreeAddScript))

	// Execute the script
	cmd := exec.Command("bash", scriptPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	output.Progress("Executing script: %s", scriptPath)
	return cmd.Run()
}

// ExecuteWorkOnScript runs the work-on script with environment variables
func ExecuteWorkOnScript(cfg *config.Config, repo *config.Repo, worktreePath string) error {
	if cfg.WorkOnScript == "" {
		return nil // No work-on script configured
	}

	scriptPath := cfg.WorkOnScript
	if !filepath.IsAbs(scriptPath) {
		scriptPath = filepath.Join(repo.Dir, scriptPath)
	}

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("work-on script does not exist: %s", scriptPath)
	}

	// Prepare environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("REPO_ALIAS=%s", repo.Alias))
	env = append(env, fmt.Sprintf("REPO_DIR=%s", repo.Dir))
	env = append(env, fmt.Sprintf("WORKTREE_PATH=%s", worktreePath))
	env = append(env, fmt.Sprintf("POST_WORKTREE_ADD_SCRIPT=%s", repo.PostWorktreeAddScript))

	// Execute the script
	cmd := exec.Command("bash", scriptPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = worktreePath // Set working directory to the worktree

	output.Progress("Executing work-on script: %s", scriptPath)
	return cmd.Run()
}
