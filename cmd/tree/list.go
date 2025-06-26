package tree

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/contextkeys"
	"worktree-manager/internal/output"
	"worktree-manager/internal/worktree"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List worktrees for the current repository",
	Long:  `List all worktrees for the current repository. Must be run from within a repository managed by worktree-manager.`,
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	cfg := cmd.Context().Value(contextkeys.ConfigKey).(*config.Config)
	
	if err := worktree.ListWorktrees(cfg); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return nil
}


