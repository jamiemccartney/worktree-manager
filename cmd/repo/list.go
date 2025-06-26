package repo

import (
	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/contextkeys"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured repositories",
	Long:  `List all repositories configured in worktree-manager.`,
	RunE:  runRepoList,
}

func runRepoList(cmd *cobra.Command, args []string) error {
	cfg := cmd.Context().Value(contextkeys.ConfigKey).(*config.Config)

	config.PrintRepoList(cfg.Repos)
	return nil
}
