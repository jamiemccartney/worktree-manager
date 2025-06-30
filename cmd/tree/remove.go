package tree

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
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
	appState := state.GetStateFromContext(cmd.Context())

	if err := worktree.RemoveWorktree(appState, branch); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return nil
}
