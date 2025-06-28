package repo

import (
	"github.com/spf13/cobra"
	"worktree-manager/internal/state"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured repositories",
	Long:  `List all repositories configured in worktree-manager.`,
	RunE:  runRepoList,
}

func runRepoList(cmd *cobra.Command, args []string) error {
	appState := state.GetStateFromContext(cmd.Context())

	state.PrintRepoList(appState.Repos)
	return nil
}
