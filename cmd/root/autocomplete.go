package root

import (
	"github.com/spf13/cobra"
	"worktree-manager/cmd/autocomplete"
)

var AutocompleteCmd = &cobra.Command{
	Use:   "completion",
	Short: "Install shell autocompletion",
	Long:  `Install shell autocompletion for worktree-manager.`,
}

func init() {
	AutocompleteCmd.AddCommand(autocomplete.BashCmd)
	AutocompleteCmd.AddCommand(autocomplete.ZshCmd)
}
