package root

import (
	"github.com/spf13/cobra"
	"worktree-manager/cmd/config"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage worktree-manager configuration",
	Long:  `Commands for managing the worktree-manager configuration file.`,
}

func init() {
	ConfigCmd.AddCommand(config.EditCmd)
	ConfigCmd.AddCommand(config.ShowCmd)
}
