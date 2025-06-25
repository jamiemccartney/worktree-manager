package repo

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured repositories",
	Long:  `List all repositories configured in worktree-manager.`,
	RunE:  runRepoList,
}

func runRepoList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		output.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	output.PrintRepoList(cfg.Repos)
	return nil
}