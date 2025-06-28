package consts

import (
	"os"
	"path/filepath"
)

type DirectoryPaths struct {
	WorktreeManagerDir  string
	DefaultGitReposDir  string
	DefaultWorktreesDir string
	ScriptsDir          string
	RepoScriptsDir      func(string) string
}

func GetDirectoryPaths() DirectoryPaths {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	worktreeManagerDir := filepath.Join(homeDir, ".worktree-manager")
	scriptsDir := filepath.Join(worktreeManagerDir, "scripts")

	return DirectoryPaths{
		WorktreeManagerDir:  worktreeManagerDir,
		DefaultGitReposDir:  filepath.Join(worktreeManagerDir, "repos"),
		DefaultWorktreesDir: filepath.Join(worktreeManagerDir, "worktrees"),
		ScriptsDir:          scriptsDir,
		RepoScriptsDir: func(repoAlias string) string {
			return filepath.Join(scriptsDir, repoAlias)
		},
	}
}
