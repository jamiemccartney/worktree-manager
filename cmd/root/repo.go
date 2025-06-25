package root

import (
	"github.com/spf13/cobra"
	"worktree-manager/cmd/repo"
)

var RepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories",
	Long:  `Commands for managing repositories used with worktree-manager.`,
}

func init() {
	RepoCmd.AddCommand(repo.CloneCmd)
	RepoCmd.AddCommand(repo.ListCmd)
	RepoCmd.AddCommand(repo.RemoveCmd)
	RepoCmd.AddCommand(repo.UseCmd)
	RepoCmd.AddCommand(repo.CurrentCmd)
}