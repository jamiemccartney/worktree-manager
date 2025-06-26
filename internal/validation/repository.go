package validation

import (
	"fmt"
	"os"
	"path/filepath"

	"worktree-manager/internal/config"
	"worktree-manager/internal/git"
)

func ValidateGitRepository(path string) error {
	if !git.IsGitRepository(path) {
		return fmt.Errorf("directory is not a git repository: %s", path)
	}
	return nil
}

func ValidateWorktreeStructure(repoPath string) error {
	if err := ValidateGitRepository(repoPath); err != nil {
		return err
	}
	
	worktreesDir := filepath.Join(repoPath, "worktrees")
	if stat, err := os.Stat(worktreesDir); err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("worktrees path exists but is not a directory: %s", worktreesDir)
		}
	}
	
	return nil
}

func ValidateScriptPath(scriptPath, repoDir string) error {
	if scriptPath == "" {
		return nil
	}
	
	if !filepath.IsAbs(scriptPath) {
		scriptPath = filepath.Join(repoDir, scriptPath)
	}
	
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script does not exist: %s", scriptPath)
	}
	
	return nil
}

func ValidateRepositoryConfig(repo *config.Repo) error {
	if repo.Alias == "" {
		return fmt.Errorf("repository alias cannot be empty")
	}
	
	if repo.Dir == "" {
		return fmt.Errorf("repository directory cannot be empty")
	}
	
	if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
		return fmt.Errorf("repository directory does not exist: %s", repo.Dir)
	}
	
	if err := ValidateWorktreeStructure(repo.Dir); err != nil {
		return fmt.Errorf("invalid repository structure: %w", err)
	}
	
	if err := ValidateScriptPath(repo.PostWorktreeAddScript, repo.Dir); err != nil {
		return fmt.Errorf("invalid post-worktree-add script: %w", err)
	}
	
	return nil
}

func ValidateWorktreeExists(worktreePath, branch string) error {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree for branch '%s' does not exist at %s", branch, worktreePath)
	}
	return nil
}

func ValidateWorktreeDoesNotExist(worktreePath, branch string) error {
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		return fmt.Errorf("worktree path '%s' already exists", worktreePath)
	}
	return nil
}

func ValidateConfigurationHealth(cfg *config.Config) []error {
	var errors []error
	
	if _, err := os.Stat(cfg.GitReposDir); os.IsNotExist(err) {
		errors = append(errors, fmt.Errorf("git repos directory does not exist: %s", cfg.GitReposDir))
	}
	
	if cfg.WorkOnScript != "" {
		if _, err := os.Stat(cfg.WorkOnScript); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("work-on script does not exist: %s", cfg.WorkOnScript))
		}
	}
	
	for _, repo := range cfg.Repos {
		if err := ValidateRepositoryConfig(&repo); err != nil {
			errors = append(errors, fmt.Errorf("repository '%s': %w", repo.Alias, err))
		}
	}
	
	return errors
}