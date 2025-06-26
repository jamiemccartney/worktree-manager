package root

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/output"
	"worktree-manager/internal/validation"
)

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the health of the worktree-manager configuration",
	Long:  `Validate the configuration and check for any issues with repositories or scripts.`,
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	output.Progress("Running worktree-manager health check...")

	if !config.ConfigExists() {
		output.Error("Config file does not exist. Run 'wt init' to create it.")
		os.Exit(1)
	}
	output.Success("Config file exists")

	cfg, err := config.Load()
	if err != nil {
		output.Error("Failed to load config: %v", err)
		os.Exit(1)
	}
	output.Success("Config file is valid JSON")

	errors := validation.ValidateConfigurationHealth(cfg)

	if _, err := os.Stat(cfg.GitReposDir); os.IsNotExist(err) {
		output.Warning("Git repos directory does not exist: %s", cfg.GitReposDir)
		output.Hint("Run 'mkdir -p %s' to create it", cfg.GitReposDir)
	} else {
		output.Success("Git repos directory exists: %s", cfg.GitReposDir)
	}

	if _, err := os.Stat(cfg.WorktreesDir); os.IsNotExist(err) {
		output.Warning("Worktrees directory does not exist: %s", cfg.WorktreesDir)
		output.Hint("Run 'mkdir -p %s' to create it", cfg.WorktreesDir)
	} else {
		output.Success("Worktrees directory exists: %s", cfg.WorktreesDir)
	}

	if cfg.ConfigEditor != "" {
		output.Success("Config editor set to: %s", cfg.ConfigEditor)
	} else {
		output.Warning("No config editor specified, will use 'vi' as default")
	}

	if cfg.WorkOnScript != "" {
		if _, err := os.Stat(cfg.WorkOnScript); os.IsNotExist(err) {
			output.Error("Work-on script does not exist: %s", cfg.WorkOnScript)
		} else {
			output.Success("Work-on script exists: %s", cfg.WorkOnScript)
		}
	}

	output.Info("Checking %d configured repositories:", len(cfg.Repos))

	for _, repo := range cfg.Repos {
		output.Info("🔍 Checking repository '%s':", repo.Alias)

		if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
			output.Error("  Repository directory does not exist: %s", repo.Dir)
			output.Hint("     Run 'wt repo clone <url>' to clone it again")
		} else {
			output.Success("  Repository directory exists: %s", repo.Dir)

			gitDir := filepath.Join(repo.Dir, ".git")
			if _, err := os.Stat(gitDir); os.IsNotExist(err) {
				output.Error("  Directory is not a git repository (no .git directory)")
			} else {
				output.Success("  Valid git repository")
			}

			repoWorktreesDir := filepath.Join(cfg.WorktreesDir, repo.Alias)
			if _, err := os.Stat(repoWorktreesDir); os.IsNotExist(err) {
				output.Warning("  Worktrees directory does not exist: %s", repoWorktreesDir)
				output.Hint("     It will be created when you add your first worktree")
			} else {
				output.Success("  Worktrees directory exists: %s", repoWorktreesDir)
			}
		}

		if repo.PostWorktreeAddScript != "" {
			scriptPath := repo.PostWorktreeAddScript
			if !filepath.IsAbs(scriptPath) {
				scriptPath = filepath.Join(repo.Dir, repo.PostWorktreeAddScript)
			}

			if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
				output.Error("  Post-worktree-add script does not exist: %s", scriptPath)
			} else {
				output.Success("  Post-worktree-add script exists: %s", scriptPath)
			}
		}
	}

	if len(errors) == 0 {
		output.Success("All checks passed! Your worktree-manager is ready to use.")
	} else {
		output.Warning("Some issues were found. Please address them before using worktree-manager.")
		for _, err := range errors {
			output.Error(err.Error())
		}
		os.Exit(1)
	}

	return nil
}
