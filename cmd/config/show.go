package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	cfg, err := configpkg.Load()
	if err != nil {
		output.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Pretty print the resolved configuration
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		output.Error("Failed to marshal config: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration (resolved):\n%s\n", string(data))
	fmt.Printf("\nConfig file location: %s\n", configpkg.GetConfigPath())
	return nil
}
