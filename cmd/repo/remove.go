package repo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"worktree-manager/internal/consts"
	"worktree-manager/internal/git"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove a repository",
	Long:  `Remove a repository from the configuration and optionally delete its directory.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRepoRemove,
}

func runRepoRemove(cmd *cobra.Command, args []string) error {
	alias := args[0]
	appState := state.GetStateFromContext(cmd.Context())

	repo, err := appState.FindRepoByAlias(alias)
	if err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}

	worktrees, err := git.ListWorktrees(repo.Dir)
	if err != nil {
		output.Warning("Could not list worktrees: %v", err)
	}

	output.Warning("This will remove repository '%s' from the configuration.", alias)
	output.Info("Directory: %s", repo.Dir)
	if len(worktrees) > 0 {
		output.Info("Found %d worktrees associated with this repository.", len(worktrees))
	}

	output.Question("Do you also want to delete the repository and its worktrees directories? [y/N]: ")

	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		output.Error("Failed to read input: %v", err)
		os.Exit(1)
	}

	deleteDir := response == "y" || response == "Y" || response == "yes"

	if deleteDir {
		if err := os.RemoveAll(repo.Dir); err != nil {
			output.Error("Failed to delete repository directory: %v", err)
			os.Exit(1)
		}
		output.Cleanup("Deleted repository directory: %s", repo.Dir)

		worktreesDir := filepath.Join(consts.GetDirectoryPaths().DefaultWorktreesDir, repo.Alias)
		if _, err := os.Stat(worktreesDir); !os.IsNotExist(err) {
			if err := os.RemoveAll(worktreesDir); err != nil {
				output.Error("Failed to delete worktrees directory: %v", err)
				os.Exit(1)
			}
			output.Cleanup("Deleted worktrees directory: %s", worktreesDir)
		}
	}

	if err := appState.RemoveRepo(alias); err != nil {
		output.Error("Failed to remove repository from state: %v", err)
		os.Exit(1)
	}

	output.Success("Repository '%s' removed from configuration", alias)
	return nil
}
