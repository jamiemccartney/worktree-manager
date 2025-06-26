package context

import (
	"fmt"
	"os"
	"path/filepath"

	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

type RepoContext struct {
	Config      *config.Config
	CurrentRepo *config.Repo
	OriginalDir string
}

func NewRepoContext(cfg *config.Config) (*RepoContext, error) {
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

func (rc *RepoContext) ChangeToRepoDir() error {
	return os.Chdir(rc.CurrentRepo.Dir)
}

func (rc *RepoContext) RestoreOriginalDir() error {
	return os.Chdir(rc.OriginalDir)
}

func (rc *RepoContext) GetWorktreePath(branch string) string {
	return filepath.Join(rc.Config.WorktreesDir, rc.CurrentRepo.Alias, branch)
}

func (rc *RepoContext) GetWorktreesDir() string {
	return filepath.Join(rc.Config.WorktreesDir, rc.CurrentRepo.Alias)
}

func (rc *RepoContext) EnsureWorktreesDir() error {
	worktreesDir := rc.GetWorktreesDir()
	return os.MkdirAll(worktreesDir, 0755)
}

func (rc *RepoContext) WithRepoDir(fn func() error) error {
	if err := rc.ChangeToRepoDir(); err != nil {
		return fmt.Errorf("failed to change to repo directory: %w", err)
	}

	defer func() {
		if restoreErr := rc.RestoreOriginalDir(); restoreErr != nil {
			output.Warning("Failed to restore original directory: %v", restoreErr)
		}
	}()

	return fn()
}
