package tree

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/contextkeys"
	"worktree-manager/internal/output"
	"worktree-manager/internal/worktree"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove <branch>",
	Short: "Remove a worktree for the specified branch",
	Long:  `Remove the worktree for the specified branch in the current repository. Must be run from within a repository managed by worktree-manager.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	branch := args[0]
	cfg := cmd.Context().Value(contextkeys.ConfigKey).(*config.Config)
	
	if err := worktree.RemoveWorktree(cfg, branch); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return nil
}

