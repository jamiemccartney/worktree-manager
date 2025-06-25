package root

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize worktree-manager configuration",
	Long:  `Create the initial configuration file for worktree-manager. Use --force to reinitialize if config already exists.`,
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		output.Error("Failed to read force flag: %v", err)
		os.Exit(1)
	}
	
	if err := config.CreateDefault(force); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return nil
}

func init() {
	InitCmd.Flags().BoolP("force", "f", false, "Force reinitialize even if config already exists")
}