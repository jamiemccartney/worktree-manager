package repo

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"worktree-manager/internal/consts"
	gitutils "worktree-manager/internal/git"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
)

var CloneCmd = &cobra.Command{
	Use:   "clone <url>",
	Short: "Clone a repository for use with worktrees",
	Long:  `Clone a repository as a bare repository for use with worktrees. The repository alias will be automatically derived from the repository name, or you can specify a custom alias using the --alias flag.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRepoClone,
}

func runRepoClone(cmd *cobra.Command, args []string) error {
	url := args[0]

	customAlias, err := cmd.Flags().GetString("alias")
	if err != nil {
		output.Error("Failed to read alias flag: %v", err)
		os.Exit(1)
	}

	var alias string
	if customAlias != "" {
		alias = customAlias
		output.Info("Using custom alias '%s'", alias)
	} else {
		alias = gitutils.ExtractRepoNameFromURL(url)
		output.Info("Using alias '%s' derived from repository URL", alias)
	}

	appState := state.GetStateFromContext(cmd.Context())

	if _, err := appState.FindRepoByAlias(alias); err == nil {
		output.Error("Repository with alias '%s' already exists", alias)
		os.Exit(1)
	}

	gitReposDir := consts.GetDirectoryPaths().DefaultGitReposDir
	if err := os.MkdirAll(gitReposDir, 0755); err != nil {
		output.Error("Failed to create git repos directory: %v", err)
		os.Exit(1)
	}

	repoDir := filepath.Join(gitReposDir, alias)

	if _, err := os.Stat(repoDir); !os.IsNotExist(err) {
		output.Error("Directory already exists: %s", repoDir)
		os.Exit(1)
	}

	output.Progress("Cloning repository: %s", url)

	gitCmd := exec.Command("git", "clone", url, repoDir)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		output.Error("Failed to clone repository: %v", err)
		os.Exit(1)
	}

	repo := state.Repo{
		Alias: alias,
		Dir:   repoDir,
	}

	if err := appState.AddRepo(repo); err != nil {
		output.Error("Failed to add repository to state: %v", err)
		os.Exit(1)
	}

	output.Success("Repository '%s' cloned and configured with alias '%s' at: %s", url, alias, repoDir)
	return nil
}

func init() {
	CloneCmd.Flags().StringP("alias", "a", "", "Custom alias for the repository (defaults to repository name)")
}
