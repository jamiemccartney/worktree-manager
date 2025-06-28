package root

import (
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/consts"
	"worktree-manager/internal/fileops"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
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

	// Create configuration
	if err := config.CreateDefault(force); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	output.Info("Config created at: %s", config.GetConfigPath())

	// Create state
	if err := state.CreateDefault(); err != nil {
		output.Error("Failed to create state: %v", err)
		os.Exit(1)
	}

	// Create scripts directory and work-on script
	if err := createScripts(); err != nil {
		output.Error("Failed to create scripts: %v", err)
		os.Exit(1)
	}

	output.Success("Init complete")
	return nil
}

func createScripts() error {
	// Create scripts directory
	scriptsDir := consts.GetDirectoryPaths().ScriptsDir
	if err := fileops.EnsureDir(scriptsDir); err != nil {
		return err
	}
	output.Info("Scripts Directory created at: %s", scriptsDir)

	// Create work-on script
	workOnScript := consts.GetFilePaths().WorkOnScript
	workOnContent := consts.GetWorkOnScriptContent()

	if err := fileops.CreateExecutableScript(workOnScript, workOnContent); err != nil {
		return err
	}
	output.Info("Work-on script created at: %s", workOnScript)

	return nil
}

func init() {
	InitCmd.Flags().BoolP("force", "f", false, "Force reinitialize even if config already exists")
}
