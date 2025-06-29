package tree

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
	"worktree-manager/internal/worktree"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List worktrees for the current repository",
	Long:  `List all worktrees for the current repository. Must be run from within a repository managed by worktree-manager.`,
	RunE:  runList,
}

func init() {
	ListCmd.Flags().Bool("json", false, "Output in JSON format for autocompletion")
}

func runList(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfigFromContext(cmd.Context())
	appState := state.GetStateFromContext(cmd.Context())

	// Check if JSON format is requested
	jsonFormat, _ := cmd.Flags().GetBool("json")
	
	if jsonFormat {
		if err := worktree.ListWorktreesJSON(cfg, appState); err != nil {
			output.Error("%v", err)
			os.Exit(1)
		}
	} else {
		if err := worktree.ListWorktrees(cfg, appState); err != nil {
			output.Error("%v", err)
			os.Exit(1)
		}
	}
	return nil
}
