package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FetchFromOrigin(repoDir string) error {
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = repoDir
	return cmd.Run()
}

func RemoteBranchExists(repoDir, branch string) bool {
	cmd := exec.Command("git", "ls-remote", "--exit-code", "--heads", "origin", branch)
	cmd.Dir = repoDir
	return cmd.Run() == nil
}

func GetBaseBranch(repoDir string) (string, error) {
	mainCmd := exec.Command("git", "ls-remote", "--exit-code", "--heads", "origin", "main")
	mainCmd.Dir = repoDir

	if mainCmd.Run() == nil {
		return "origin/main", nil
	}

	masterCmd := exec.Command("git", "ls-remote", "--exit-code", "--heads", "origin", "master")
	masterCmd.Dir = repoDir

	if masterCmd.Run() == nil {
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
	args := []string{"worktree", "add"}

	if opts.CreateBranch {
		args = append(args, "-b", opts.Branch)
	}

	args = append(args, opts.WorktreePath)

	if opts.SourceBranch != "" {
		args = append(args, opts.SourceBranch)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git command failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func RemoveWorktree(repoDir, worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", "--force", worktreePath)
	cmd.Dir = repoDir
	return cmd.Run()
}

func ListWorktrees(repoDir string) ([]Worktree, error) {
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
