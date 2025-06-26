package repo

import (
	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/contextkeys"
	"worktree-manager/internal/output"
)

var CurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently active repository",
	Long:  `Display the currently active repository for tree commands.`,
	RunE:  runRepoCurrent,
}

func runRepoCurrent(cmd *cobra.Command, args []string) error {
	cfg := cmd.Context().Value(contextkeys.ConfigKey).(*config.Config)

	activeRepo, err := cfg.GetActiveRepo()
	if err != nil {
		output.Warning("No active repository set")
		output.Hint("Use 'wt repo use <alias>' to set an active repository")
		return nil
	}

	output.Success("Active repository: %s", activeRepo.Alias)
	output.Info("Location: %s", activeRepo.Dir)

	return nil
}
