package autocomplete

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"worktree-manager/internal/output"
)

var BashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Install bash autocompletion",
	Long:  `Install bash autocompletion for worktree-manager.`,
	RunE:  runAutocompleteBash,
}

func runAutocompleteBash(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		output.Error("Failed to get home directory: %v", err)
		os.Exit(1)
	}

	completionScript := `#!/bin/bash
_wt_completion() {
    local cur prev words cword
    _init_completion || return

    case $prev in
        add|workon)
            # Complete with branch names from current repo worktrees
            if command -v wt >/dev/null 2>&1; then
                local branches=($(wt tree list 2>/dev/null | grep "🔸" | head -20))
                COMPREPLY=($(compgen -W "${branches[*]}" -- "$cur"))
            fi
            return
            ;;
        remove)
            # Complete with available worktrees that can be removed
            if command -v wt >/dev/null 2>&1 && command -v jq >/dev/null 2>&1; then
                local worktrees=($(wt tree list --json 2>/dev/null | jq -r '.[]' 2>/dev/null))
                COMPREPLY=($(compgen -W "${worktrees[*]}" -- "$cur"))
            fi
            return
            ;;
        clone)
            # No completion for URLs
            return
            ;;
        use)
            # Complete with repository aliases
            if command -v wt >/dev/null 2>&1; then
                local repos=($(wt repo list 2>/dev/null | grep "🔸" | cut -d' ' -f2))
                COMPREPLY=($(compgen -W "${repos[*]}" -- "$cur"))
            fi
            return
            ;;
    esac

    case $cword in
        1)
            # First level commands
            local commands="init doctor config repo tree autocomplete"
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
            ;;
        2)
            case ${words[1]} in
                config)
                    COMPREPLY=($(compgen -W "edit show" -- "$cur"))
                    ;;
                repo)
                    COMPREPLY=($(compgen -W "clone list remove use" -- "$cur"))
                    ;;
                tree)
                    COMPREPLY=($(compgen -W "add remove list workon" -- "$cur"))
                    ;;
                autocomplete)
                    COMPREPLY=($(compgen -W "bash zsh" -- "$cur"))
                    ;;
            esac
            ;;
    esac
}

complete -F _wt_completion wt
complete -F _wt_completion worktree-manager
`

	bashCompletionDir := filepath.Join(homeDir, ".bash_completion.d")
	if err := os.MkdirAll(bashCompletionDir, 0755); err != nil {
		output.Error("Failed to create bash completion directory: %v", err)
		os.Exit(1)
	}

	completionFile := filepath.Join(bashCompletionDir, "wt")
	if err := os.WriteFile(completionFile, []byte(completionScript), 0644); err != nil {
		output.Error("Failed to write completion script: %v", err)
		os.Exit(1)
	}

	output.Success("Bash completion installed to: %s", completionFile)
	output.Hint("To enable completion, add this to your ~/.bashrc:")
	output.Info("   source %s", completionFile)
	output.Hint("Or restart your terminal to load completions automatically.")

	return nil
}
