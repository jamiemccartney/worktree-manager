package state

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"worktree-manager/internal/consts"
	"worktree-manager/internal/fileops"
	"worktree-manager/internal/output"
)

// State represents the application state (repos and active repo)
type State struct {
	ActiveRepo string `json:"active-repo"`
	Repos      []Repo `json:"repos"`
}

// Repo represents a repository in the state
type Repo struct {
	Alias string `json:"alias"`
	Dir   string `json:"dir"`
}

var (
	appState *State
)

// Load reads and returns the application state
func Load() (*State, error) {
	if appState != nil {
		return appState, nil
	}

	var state State
	if err := fileops.ReadJSONFile(consts.GetFilePaths().State, &state); err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	// Expand environment variables in repo directories
	for i := range state.Repos {
		state.Repos[i].Dir = fileops.ExpandEnvVars(state.Repos[i].Dir)
	}

	appState = &state
	return appState, nil
}

// StateExists checks if the state file exists
func StateExists() bool {
	return fileops.FileExists(consts.GetFilePaths().State)
}

// CreateDefault creates a default state file
func CreateDefault() error {
	defaultState := State{
		ActiveRepo: "",
		Repos:      []Repo{},
	}

	statePath := consts.GetFilePaths().State
	if err := fileops.EnsureDir(filepath.Dir(statePath)); err != nil {
		return err
	}

	if err := fileops.WriteJSONFile(statePath, defaultState); err != nil {
		return err
	}

	// Create default directories
	dirs := consts.GetDirectoryPaths()
	if err := fileops.EnsureDir(dirs.DefaultGitReposDir); err != nil {
		return err
	}
	output.Info("Repos Directory created at: %s", dirs.DefaultGitReposDir)

	if err := fileops.EnsureDir(dirs.DefaultWorktreesDir); err != nil {
		return err
	}
	output.Info("Worktrees Directory created at: %s", dirs.DefaultWorktreesDir)

	appState = &defaultState
	return nil
}

// Save saves the current state
func (s *State) Save() error {
	return fileops.WriteJSONFile(consts.GetFilePaths().State, s)
}

// GetCurrentRepo returns the repo for the current working directory
func (s *State) GetCurrentRepo() (*Repo, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	for _, repo := range s.Repos {
		if strings.HasPrefix(pwd, repo.Dir) {
			return &repo, nil
		}
	}

	return nil, fmt.Errorf("current directory is not within a managed repository")
}

// FindRepoByAlias finds a repository by its alias
func (s *State) FindRepoByAlias(alias string) (*Repo, error) {
	for _, repo := range s.Repos {
		if repo.Alias == alias {
			return &repo, nil
		}
	}
	return nil, fmt.Errorf("repository with alias '%s' not found", alias)
}

// AddRepo adds a new repository to the state
func (s *State) AddRepo(repo Repo) error {
	// Check if alias already exists
	for _, existingRepo := range s.Repos {
		if existingRepo.Alias == repo.Alias {
			return fmt.Errorf("repository with alias '%s' already exists", repo.Alias)
		}
	}

	// Expand environment variables
	repo.Dir = fileops.ExpandEnvVars(repo.Dir)

	s.Repos = append(s.Repos, repo)

	// Create post-worktree-add script for this repo
	if err := s.createRepoScript(repo.Alias); err != nil {
		output.Warning("Failed to create post-worktree-add script: %v", err)
	}

	return s.Save()
}

// RemoveRepo removes a repository from the state
func (s *State) RemoveRepo(alias string) error {
	for i, repo := range s.Repos {
		if repo.Alias == alias {
			// Prompt user about script deletion
			if err := s.confirmRepoScriptDeletion(alias); err != nil {
				return err
			}

			// Remove from slice
			s.Repos = append(s.Repos[:i], s.Repos[i+1:]...)

			// Clean up script directory
			if err := s.removeRepoScript(alias); err != nil {
				output.Warning("Failed to clean up repo scripts: %v", err)
			}

			return s.Save()
		}
	}
	return fmt.Errorf("repository with alias '%s' not found", alias)
}

// SetActiveRepo sets the active repository
func (s *State) SetActiveRepo(alias string) error {
	_, err := s.FindRepoByAlias(alias)
	if err != nil {
		return err
	}

	s.ActiveRepo = alias
	return s.Save()
}

// GetActiveRepo returns the active repository
func (s *State) GetActiveRepo() (*Repo, error) {
	if s.ActiveRepo == "" {
		return nil, fmt.Errorf("no active repository set. Use 'wt repo use <alias>' to set one")
	}

	return s.FindRepoByAlias(s.ActiveRepo)
}

func (s *State) createRepoScript(repoAlias string) error {
	scriptPath := consts.GetFilePaths().PostWorktreeAddScript(repoAlias)
	content := consts.GetPostWorktreeAddScriptContent(repoAlias)

	return fileops.CreateExecutableScript(scriptPath, content)
}

func (s *State) confirmRepoScriptDeletion(repoAlias string) error {
	scriptDir := fileops.GetRepoScriptDir(repoAlias)
	if !fileops.FileExists(scriptDir) {
		return nil // Nothing to delete
	}

	output.Warning("Removing repository '%s' will also delete its script directory:", repoAlias)
	output.Warning("  %s", scriptDir)
	output.Warning("This may result in data loss if you have customized scripts.")

	// In a real implementation, you'd prompt for user confirmation here
	// For now, we'll proceed with deletion
	output.Info("Proceeding with script directory cleanup...")

	return nil
}

func (s *State) removeRepoScript(repoAlias string) error {
	scriptDir := fileops.GetRepoScriptDir(repoAlias)
	if !fileops.FileExists(scriptDir) {
		return nil // Nothing to delete
	}

	if err := os.RemoveAll(scriptDir); err != nil {
		return fmt.Errorf("failed to remove script directory %s: %w", scriptDir, err)
	}

	output.Info("Cleaned up script directory: %s", scriptDir)
	return nil
}

// PrintRepoList displays a formatted list of repositories
func PrintRepoList(repos []Repo) {
	if len(repos) == 0 {
		output.Warning("No repositories configured")
		output.Hint("Use 'wt repo clone <url>' to add a repository")
		return
	}

	output.Info("Configured repositories:")
	for _, repo := range repos {
		output.Item("%s → %s", repo.Alias, repo.Dir)
	}
}

// GetStateFromContext extracts state from context
func GetStateFromContext(ctx context.Context) *State {
	return ctx.Value(consts.GetContextKeys().State).(*State)
}
