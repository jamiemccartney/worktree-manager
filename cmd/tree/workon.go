package tree

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
	"worktree-manager/internal/worktree"
)

var WorkonCmd = &cobra.Command{
	Use:   "workon <branch>",
	Short: "Work on a specific worktree",
	Long:  `Change to a worktree directory and run the work-on script if configured. Must be run from within a repository managed by worktree-manager.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkon,
}

func runWorkon(cmd *cobra.Command, args []string) error {
	branch := args[0]
	appState := state.GetStateFromContext(cmd.Context())

	if err := worktree.WorkOnWorktree(appState, branch); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return nil
}
