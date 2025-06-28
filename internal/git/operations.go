package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"worktree-manager/internal/executors"
)

type GitOperations struct {
	cmdExecutor executors.CommandExecutor
}

func NewGitOperations() *GitOperations {
	return &GitOperations{
		cmdExecutor: executors.NewSystemCommandExecutor(),
	}
}

var defaultGitOps = NewGitOperations()

func FetchFromOrigin(repoDir string) error {
	return defaultGitOps.FetchFromOrigin(repoDir)
}

func (g *GitOperations) FetchFromOrigin(repoDir string) error {
	ctx := &executors.CommandExecutionContext{
		Command:    "git",
		Args:       []string{"fetch", "origin"},
		WorkingDir: repoDir,
	}
	return g.cmdExecutor.Execute(ctx)
}

func RemoteBranchExists(repoDir, branch string) bool {
	return defaultGitOps.RemoteBranchExists(repoDir, branch)
}

func (g *GitOperations) RemoteBranchExists(repoDir, branch string) bool {
	ctx := &executors.CommandExecutionContext{
		Command:    "git",
		Args:       []string{"ls-remote", "--exit-code", "--heads", "origin", branch},
		WorkingDir: repoDir,
	}
	err := g.cmdExecutor.Execute(ctx)
	return err == nil
}

func GetBaseBranch(repoDir string) (string, error) {
	return defaultGitOps.GetBaseBranch(repoDir)
}

func (g *GitOperations) GetBaseBranch(repoDir string) (string, error) {
	mainCtx := &executors.CommandExecutionContext{
		Command:    "git",
		Args:       []string{"ls-remote", "--exit-code", "--heads", "origin", "main"},
		WorkingDir: repoDir,
	}
	if err := g.cmdExecutor.Execute(mainCtx); err == nil {
		return "origin/main", nil
	}

	masterCtx := &executors.CommandExecutionContext{
		Command:    "git",
		Args:       []string{"ls-remote", "--exit-code", "--heads", "origin", "master"},
		WorkingDir: repoDir,
	}
	if err := g.cmdExecutor.Execute(masterCtx); err == nil {
		return "origin/master", nil
	}

	return "", fmt.Errorf("neither origin/main nor origin/master exists")
}

func IsGitRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	_, err := os.Stat(gitDir)
	return !os.IsNotExist(err)
}

func CreateWorktree(repoDir string, opts WorktreeCreateOptions) error {
	return defaultGitOps.CreateWorktree(repoDir, opts)
}

func (g *GitOperations) CreateWorktree(repoDir string, opts WorktreeCreateOptions) error {
	args := []string{"worktree", "add"}

	if opts.CreateBranch {
		args = append(args, "-b", opts.Branch)
	}

	args = append(args, opts.WorktreePath)

	if opts.SourceBranch != "" {
		args = append(args, opts.SourceBranch)
	}

	ctx := &executors.CommandExecutionContext{
		Command:    "git",
		Args:       args,
		WorkingDir: repoDir,
		ShowOutput: true,
	}
	return g.cmdExecutor.Execute(ctx)
}

func RemoveWorktree(repoDir, worktreePath string) error {
	return defaultGitOps.RemoveWorktree(repoDir, worktreePath)
}

func (g *GitOperations) RemoveWorktree(repoDir, worktreePath string) error {
	ctx := &executors.CommandExecutionContext{
		Command:    "git",
		Args:       []string{"worktree", "remove", "--force", worktreePath},
		WorkingDir: repoDir,
	}
	return g.cmdExecutor.Execute(ctx)
}

func ListWorktrees(repoDir string) ([]Worktree, error) {
	return defaultGitOps.ListWorktrees(repoDir)
}

func (g *GitOperations) ListWorktrees(repoDir string) ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = repoDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return parseWorktreeList(string(output)), nil
}

func parseWorktreeList(output string) []Worktree {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var worktrees []Worktree
	var current *Worktree

	for _, line := range lines {
		if line == "" {
			if current != nil {
				worktrees = append(worktrees, *current)
				current = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			if current != nil {
				worktrees = append(worktrees, *current)
			}
			current = &Worktree{
				Path: strings.TrimPrefix(line, "worktree "),
			}
		} else if current != nil {
			if strings.HasPrefix(line, "branch ") {
				current.Branch = strings.TrimPrefix(line, "branch ")
			}
		}
	}

	if current != nil {
		worktrees = append(worktrees, *current)
	}

	return worktrees
}
