package context

import (
	"fmt"
	"os"
	"path/filepath"

	"worktree-manager/internal/config"
)

// RepoContext manages the context of the current repository operation
type RepoContext struct {
	Config      *config.Config
	CurrentRepo *config.Repo
	OriginalDir string
}

// NewRepoContext creates a new repository context using the active repository from config
func NewRepoContext() (*RepoContext, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	activeRepo, err := cfg.GetActiveRepo()
	if err != nil {
		return nil, fmt.Errorf("❌ %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	return &RepoContext{
		Config:      cfg,
		CurrentRepo: activeRepo,
		OriginalDir: originalDir,
	}, nil
}

// NewRepoContextFromWorkingDir creates a new repository context by detecting current repo from working directory
// This is kept for backward compatibility if needed
func NewRepoContextFromWorkingDir() (*RepoContext, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	currentRepo, err := cfg.GetCurrentRepo()
	if err != nil {
		return nil, fmt.Errorf("❌ %v\n\n💡 You must be in a repository managed by worktree-manager.\n   Use 'wt repo clone <url>' to clone a repository first.", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	return &RepoContext{
		Config:      cfg,
		CurrentRepo: currentRepo,
		OriginalDir: originalDir,
	}, nil
}

// ChangeToRepoDir changes the current working directory to the repository directory
func (rc *RepoContext) ChangeToRepoDir() error {
	return os.Chdir(rc.CurrentRepo.Dir)
}

// RestoreOriginalDir restores the original working directory
func (rc *RepoContext) RestoreOriginalDir() error {
	return os.Chdir(rc.OriginalDir)
}

// GetWorktreePath returns the full path to a worktree for the given branch
func (rc *RepoContext) GetWorktreePath(branch string) string {
	return filepath.Join(rc.Config.WorktreesDir, rc.CurrentRepo.Alias, branch)
}

// GetWorktreesDir returns the worktrees directory path for the current repo
func (rc *RepoContext) GetWorktreesDir() string {
	return filepath.Join(rc.Config.WorktreesDir, rc.CurrentRepo.Alias)
}

// EnsureWorktreesDir creates the worktrees directory if it doesn't exist
func (rc *RepoContext) EnsureWorktreesDir() error {
	worktreesDir := rc.GetWorktreesDir()
	return os.MkdirAll(worktreesDir, 0755)
}

// WithRepoDir executes a function while in the repository directory, 
// then restores the original directory
func (rc *RepoContext) WithRepoDir(fn func() error) error {
	if err := rc.ChangeToRepoDir(); err != nil {
		return fmt.Errorf("failed to change to repo directory: %w", err)
	}
	
	defer func() {
		if restoreErr := rc.RestoreOriginalDir(); restoreErr != nil {
			// Log the error but don't override the original error
			fmt.Printf("⚠️  Failed to restore original directory: %v\n", restoreErr)
		}
	}()
	
	return fn()
}