package repo

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

var UseCmd = &cobra.Command{
	Use:   "use <alias>",
	Short: "Set the active repository for tree commands",
	Long:  `Set the specified repository as the active repository for tree commands.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRepoUse,
}

func runRepoUse(cmd *cobra.Command, args []string) error {
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

	// Check if directory exists
	if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
		output.Error("Repository directory does not exist: %s", repo.Dir)
		os.Exit(1)
	}

	// Set this repository as the active repository
	if err := cfg.SetActiveRepo(alias); err != nil {
		output.Error("Failed to set active repository: %v", err)
		os.Exit(1)
	}

	output.Success("Set '%s' as the active repository", alias)
	output.Info("Repository location: %s", repo.Dir)
	output.Hint("You can now use 'wt tree' commands without specifying a repository")

	return nil
}