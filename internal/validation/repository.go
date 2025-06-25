package validation

import (
	"fmt"
	"os"
	"path/filepath"

	"worktree-manager/internal/config"
	"worktree-manager/internal/git"
)

// ValidateGitRepository checks if the given path is a valid git repository
func ValidateGitRepository(path string) error {
	if !git.IsGitRepository(path) {
		return fmt.Errorf("directory is not a git repository: %s", path)
	}
	return nil
}

// ValidateWorktreeStructure checks if the repository has the expected worktree structure
func ValidateWorktreeStructure(repoPath string) error {
	// Check if it's a valid git repository
	if err := ValidateGitRepository(repoPath); err != nil {
		return err
	}
	
	// Check if worktrees directory exists (it's OK if it doesn't exist yet)
	worktreesDir := filepath.Join(repoPath, "worktrees")
	if stat, err := os.Stat(worktreesDir); err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("worktrees path exists but is not a directory: %s", worktreesDir)
		}
	}
	
	return nil
}

// ValidateScriptPath checks if a script exists and validates its path
func ValidateScriptPath(scriptPath, repoDir string) error {
	if scriptPath == "" {
		return nil // No script is valid
	}
	
	// Resolve relative paths relative to repo directory
	if !filepath.IsAbs(scriptPath) {
		scriptPath = filepath.Join(repoDir, scriptPath)
	}
	
	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script does not exist: %s", scriptPath)
	}
	
	return nil
}

// ValidateRepositoryConfig validates a repository configuration
func ValidateRepositoryConfig(repo *config.Repo) error {
	if repo.Alias == "" {
		return fmt.Errorf("repository alias cannot be empty")
	}
	
	if repo.Dir == "" {
		return fmt.Errorf("repository directory cannot be empty")
	}
	
	// Check if directory exists
	if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
		return fmt.Errorf("repository directory does not exist: %s", repo.Dir)
	}
	
	// Validate git repository structure
	if err := ValidateWorktreeStructure(repo.Dir); err != nil {
		return fmt.Errorf("invalid repository structure: %w", err)
	}
	
	// Validate post-worktree-add script if specified
	if err := ValidateScriptPath(repo.PostWorktreeAddScript, repo.Dir); err != nil {
		return fmt.Errorf("invalid post-worktree-add script: %w", err)
	}
	
	return nil
}

// ValidateWorktreeExists checks if a worktree exists at the given path
func ValidateWorktreeExists(worktreePath, branch string) error {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree for branch '%s' does not exist at %s", branch, worktreePath)
	}
	return nil
}

// ValidateWorktreeDoesNotExist checks if a worktree does not already exist at the given path
func ValidateWorktreeDoesNotExist(worktreePath, branch string) error {
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		return fmt.Errorf("worktree path '%s' already exists", worktreePath)
	}
	return nil
}

// ValidateConfigurationHealth performs comprehensive health checks on the configuration
func ValidateConfigurationHealth(cfg *config.Config) []error {
	var errors []error
	
	// Check git-repos-dir
	if _, err := os.Stat(cfg.GitReposDir); os.IsNotExist(err) {
		errors = append(errors, fmt.Errorf("git repos directory does not exist: %s", cfg.GitReposDir))
	}
	
	// Check work-on script if specified
	if cfg.WorkOnScript != "" {
		if _, err := os.Stat(cfg.WorkOnScript); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("work-on script does not exist: %s", cfg.WorkOnScript))
		}
	}
	
	// Check each repository
	for _, repo := range cfg.Repos {
		if err := ValidateRepositoryConfig(&repo); err != nil {
			errors = append(errors, fmt.Errorf("repository '%s': %w", repo.Alias, err))
		}
	}
	
	return errors
}