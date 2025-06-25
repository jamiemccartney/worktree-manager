package repo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/git"
	"worktree-manager/internal/output"
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

	cfg, err := config.Load()
	if err != nil {
		output.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	repo, err := cfg.FindRepoByAlias(alias)
	if err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}

	// Get worktree information
	worktrees, err := git.ListWorktrees(repo.Dir)
	if err != nil {
		output.Warning("Could not list worktrees: %v", err)
	}

	// Ask for confirmation
	output.Warning("This will remove repository '%s' from the configuration.", alias)
	fmt.Printf("Directory: %s", repo.Dir)
	if len(worktrees) > 0 {
		output.Info("Found %d worktrees associated with this repository.", len(worktrees))
	}

	fmt.Print("Do you also want to delete the repository and its worktrees directories? [y/N]: ")

	var response string
	fmt.Scanln(&response)

	deleteDir := response == "y" || response == "Y" || response == "yes"

	if deleteDir {
		// Remove repository directory
		if err := os.RemoveAll(repo.Dir); err != nil {
			output.Error("Failed to delete repository directory: %v", err)
			os.Exit(1)
		}
		output.Cleanup("Deleted repository directory: %s", repo.Dir)

		// Remove worktrees directory
		worktreesDir := filepath.Join(cfg.WorktreesDir, repo.Alias)
		if _, err := os.Stat(worktreesDir); !os.IsNotExist(err) {
			if err := os.RemoveAll(worktreesDir); err != nil {
				output.Error("Failed to delete worktrees directory: %v", err)
				os.Exit(1)
			}
			output.Cleanup("Deleted worktrees directory: %s", worktreesDir)
		}
	}

	// Remove from config
	if err := cfg.RemoveRepo(alias); err != nil {
		output.Error("Failed to remove repository from config: %v", err)
		os.Exit(1)
	}

	output.Success("Repository '%s' removed from configuration", alias)
	return nil
}
