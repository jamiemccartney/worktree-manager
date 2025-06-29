package root

import (
	"github.com/spf13/cobra"
	"worktree-manager/internal/output"
)

var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version information including build details.`,
	RunE:  runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	output.Info("Worktree Manager Version Information:")
	output.Item("Version: %s", version)
	output.Item("Commit: %s", commit)
	output.Item("Built: %s", date)
	return nil
}
