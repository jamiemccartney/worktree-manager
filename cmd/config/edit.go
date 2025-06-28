package config

import (
	"os"
	"os/exec"
	"worktree-manager/internal/consts"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Open the configuration file in the configured editor.`,
	RunE:  runConfigEdit,
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfigFromContext(cmd.Context())

	configPath := consts.GetFilePaths().Config
	editor := cfg.ConfigEditor
	if editor == "" {
		editor = "vi"
	}

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
