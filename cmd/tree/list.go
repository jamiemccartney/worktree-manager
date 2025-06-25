package tree

import (
	"os"

	"github.com/spf13/cobra"
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
	if err := worktree.ListWorktrees(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return nil
}


// Command is added to tree command in tree.go