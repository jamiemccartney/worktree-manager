package autocomplete

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"worktree-manager/internal/output"
)

var ZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Install zsh autocompletion",
	Long:  `Install zsh autocompletion for worktree-manager.`,
	RunE:  runAutocompleteZsh,
}

func runAutocompleteZsh(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		output.Error("Failed to get home directory: %v", err)
		os.Exit(1)
	}

	completionScript := `#compdef wt worktree-manager

_wt() {
    local context state line
    local -a commands

    _arguments -C \
        "1: :_wt_commands" \
        "*::arg:->args"

    case $state in
        args)
            case $words[1] in
                config)
                    _values "config commands" \
                        "edit[Edit configuration]" \
                        "show[Show configuration]"
                    ;;
                repo)
                    case $words[2] in
                        clone)
                            _message "repository URL"
                            ;;
                        use|remove)
                            _wt_repos
                            ;;
                        *)
                            _values "repo commands" \
                                "clone[Clone repository]" \
                                "list[List repositories]" \
                                "remove[Remove repository]" \
                                "use[Use repository]"
                            ;;
                    esac
                    ;;
                tree)
                    case $words[2] in
                        add|remove|workon)
                            _wt_branches
                            ;;
                        *)
                            _values "tree commands" \
                                "add[Add worktree]" \
                                "remove[Remove worktree]" \
                                "list[List worktrees]" \
                                "workon[Work on worktree]"
                            ;;
                    esac
                    ;;
                autocomplete)
                    _values "shell types" \
                        "bash[Install bash completion]" \
                        "zsh[Install zsh completion]"
                    ;;
            esac
            ;;
    esac
}

_wt_commands() {
    local commands=(
        "init:Initialize configuration"
        "doctor:Check configuration health"
        "config:Manage configuration"
        "repo:Manage repositories"
        "tree:Manage worktrees"
        "autocomplete:Install shell completion"
    )
    _describe "commands" commands
}

_wt_branches() {
    local branches=(${(f)"$(wt tree list 2>/dev/null | grep "🔸" | head -20)"})
    _describe "branches" branches
}

_wt_repos() {
    local repos=(${(f)"$(wt repo list 2>/dev/null | grep "🔸" | cut -d' ' -f2)"})
    _describe "repositories" repos
}

_wt
`

	zshCompletionDir := filepath.Join(homeDir, ".zsh", "completions")
	if err := os.MkdirAll(zshCompletionDir, 0755); err != nil {
		output.Error("Failed to create zsh completion directory: %v", err)
		os.Exit(1)
	}

	completionFile := filepath.Join(zshCompletionDir, "_wt")
	if err := os.WriteFile(completionFile, []byte(completionScript), 0644); err != nil {
		output.Error("Failed to write completion script: %v", err)
		os.Exit(1)
	}

	output.Success("Zsh completion installed to: %s", completionFile)
	output.Hint("To enable completion, add this to your ~/.zshrc:")
	output.Info("   fpath=(~/.zsh/completions $fpath)")
	output.Info("   autoload -U compinit && compinit")
	output.Hint("Or restart your terminal to load completions automatically.")

	return nil
}