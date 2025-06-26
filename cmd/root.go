package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"os"
	"worktree-manager/cmd/root"
	"worktree-manager/internal/config"
	"worktree-manager/internal/contextkeys"
	"worktree-manager/internal/output"
)

var rootCmd = &cobra.Command{
	Use:     "worktree-manager",
	Aliases: []string{"wt"},
	Short:   "A CLI tool for managing Git worktrees",
	Long:    `A command-line tool for managing Git worktrees efficiently.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		if cmd.Name() == "init" {
			return nil
		}

		cfg, err := config.Load()

		if err != nil {
			output.Error("Failed to load config: %v", err)
			output.Prompt("Have you ran wt init")
			os.Exit(1)
		}

		cmd.SetContext(context.WithValue(cmd.Context(), contextkeys.ConfigKey, cfg))
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(root.InitCmd)
	rootCmd.AddCommand(root.DoctorCmd)
	rootCmd.AddCommand(root.TreeCmd)
	rootCmd.AddCommand(root.ConfigCmd)
	rootCmd.AddCommand(root.RepoCmd)
	rootCmd.AddCommand(root.AutocompleteCmd)
}
