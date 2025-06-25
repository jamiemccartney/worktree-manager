package root

import (
	"github.com/spf13/cobra"
	"worktree-manager/cmd/tree"
)

var TreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Manage worktrees",
	Long:  `Commands for managing worktrees within repositories.`,
}

func init() {
	// Add worktree-related commands under tree
	TreeCmd.AddCommand(tree.AddCmd)
	TreeCmd.AddCommand(tree.RemoveCmd)
	TreeCmd.AddCommand(tree.ListCmd)
	TreeCmd.AddCommand(tree.WorkonCmd)
}