package tree

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/contextkeys"
	"worktree-manager/internal/output"
	"worktree-manager/internal/worktree"
)

var AddCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "Add a new worktree for the specified branch",
	Long:  `Create a new worktree for the specified branch in the current repository. Must be run from within a repository managed by worktree-manager.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	branch := args[0]
	cfg := cmd.Context().Value(contextkeys.ConfigKey).(*config.Config)
	
	if err := worktree.AddWorktree(cfg, branch); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return nil
}

