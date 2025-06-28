package config

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	configpkg "worktree-manager/internal/config"
	"worktree-manager/internal/output"
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration",
	Long:  `Display the current configuration with resolved environment variables.`,
	RunE:  runConfigShow,
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfigFromContext(cmd.Context())

	data, err := json.MarshalIndent(cfg, "", "    ")

	if err != nil {
		output.Error("Failed to marshal config: %v", err)
		os.Exit(1)
	}

	output.Success("Configuration (resolved):\n%s", string(data))
	output.Info("Config file location: %s", configpkg.GetConfigPath())
	return nil
}
