package config

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	configpkg "worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Open the configuration file in the configured editor.`,
	RunE:  runConfigEdit,
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	cfg, err := configpkg.Load()
	if err != nil {
		output.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	configPath := configpkg.GetConfigPath()
	editor := cfg.ConfigEditor
	if editor == "" {
		editor = "vi" // fallback
	}

	// Open the config file in the editor
	editorCmd := exec.Command(editor, configPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		output.Error("Failed to run editor: %v", err)
		os.Exit(1)
	}

	output.Success("Configuration edited: %s", configPath)
	return nil
}