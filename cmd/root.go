package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"worktree-manager/cmd/root"
)

var rootCmd = &cobra.Command{
	Use:     "worktree-manager",
	Aliases: []string{"wt"},
	Short:   "A CLI tool for managing Git worktrees",
	Long:    `A command-line tool for managing Git worktrees efficiently.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
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
