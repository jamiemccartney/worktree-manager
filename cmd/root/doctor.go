package root

import (
	"fmt"
	"os"
	"path/filepath"
	"worktree-manager/internal/git"

	"github.com/spf13/cobra"
	"worktree-manager/internal/config"
	"worktree-manager/internal/consts"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
)

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the health of the worktree-manager configuration",
	Long:  `Validate the configuration and check for any issues with repositories or scripts.`,
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	output.Progress("Running worktree-manager health check...")

	checkConfig()
	appState := checkState()

	errors := validateConfigurationHealth(appState)

	checkFolders()
	checkWorkOnScript()
	checkRepos(appState)

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

// checkConfig verifies the config file exists and can be loaded
func checkConfig() *config.Config {
	if !config.CheckConfigExists() {
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

	if cfg.ConfigEditor != "" {
		output.Success("Config editor set to: %s", cfg.ConfigEditor)
	} else {
		output.Warning("No config editor specified, will use 'vi' as default")
	}

	return cfg
}

// checkState verifies the state file exists and can be loaded
func checkState() *state.State {
	if !state.StateExists() {
		output.Error("State file does not exist. Run 'wt init' to create it.")
		os.Exit(1)
	}
	output.Success("State file exists")

	appState, err := state.Load()
	if err != nil {
		output.Error("Failed to load state: %v", err)
		os.Exit(1)
	}
	output.Success("State file is valid JSON")

	return appState
}

// checkFolders verifies the existence of required directories
func checkFolders() {
	paths := consts.GetDirectoryPaths()
	gitReposDir := paths.DefaultGitReposDir
	if _, err := os.Stat(gitReposDir); os.IsNotExist(err) {
		output.Warning("Git repos directory does not exist: %s", gitReposDir)
		output.Hint("Run 'mkdir -p %s' to create it", gitReposDir)
	} else {
		output.Success("Git repos directory exists: %s", gitReposDir)
	}

	worktreesDir := paths.DefaultWorktreesDir
	if _, err := os.Stat(worktreesDir); os.IsNotExist(err) {
		output.Warning("Worktrees directory does not exist: %s", worktreesDir)
		output.Hint("Run 'mkdir -p %s' to create it", worktreesDir)
	} else {
		output.Success("Worktrees directory exists: %s", worktreesDir)
	}
}

// checkWorkOnScript verifies the work-on script status
func checkWorkOnScript() {
	workOnScript := consts.GetFilePaths().WorkOnScript
	if _, err := os.Stat(workOnScript); os.IsNotExist(err) {
		output.Warning("Work-on script does not exist: %s", workOnScript)
		output.Hint("It will be created when you first use 'wt tree workon'")
	} else {
		output.Success("Work-on script exists: %s", workOnScript)
	}
}

// checkRepos verifies all configured repositories and their scripts
func checkRepos(appState *state.State) {
	output.Info("Checking %d configured repositories:", len(appState.Repos))

	worktreesDir := consts.GetDirectoryPaths().DefaultWorktreesDir

	for _, repo := range appState.Repos {
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

			repoWorktreesDir := filepath.Join(worktreesDir, repo.Alias)
			if _, err := os.Stat(repoWorktreesDir); os.IsNotExist(err) {
				output.Warning("  Worktrees directory does not exist: %s", repoWorktreesDir)
				output.Hint("     It will be created when you add your first worktree")
			} else {
				output.Success("  Worktrees directory exists: %s", repoWorktreesDir)
			}
		}

		scriptPath := consts.GetFilePaths().PostWorktreeAddScript(repo.Alias)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			output.Warning("  Post-worktree-add script does not exist: %s", scriptPath)
			output.Hint("     It will be created when you first add a worktree")
		} else {
			output.Success("  Post-worktree-add script exists: %s", scriptPath)
		}
	}
}

func validateGitRepository(path string) error {
	if !git.IsGitRepository(path) {
		return fmt.Errorf("directory is not a git repository: %s", path)
	}
	return nil
}

func validateWorktreeStructure(repoPath string) error {
	if err := validateGitRepository(repoPath); err != nil {
		return err
	}

	worktreesDir := filepath.Join(repoPath, "worktrees")
	if stat, err := os.Stat(worktreesDir); err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("worktrees path exists but is not a directory: %s", worktreesDir)
		}
	}

	return nil
}

func validateScriptPath(scriptPath, repoDir string) error {
	if scriptPath == "" {
		return nil
	}

	if !filepath.IsAbs(scriptPath) {
		scriptPath = filepath.Join(repoDir, scriptPath)
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script does not exist: %s", scriptPath)
	}

	return nil
}

func validateRepositoryConfig(repo *state.Repo) error {
	if repo.Alias == "" {
		return fmt.Errorf("repository alias cannot be empty")
	}

	if repo.Dir == "" {
		return fmt.Errorf("repository directory cannot be empty")
	}

	if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
		return fmt.Errorf("repository directory does not exist: %s", repo.Dir)
	}

	if err := validateWorktreeStructure(repo.Dir); err != nil {
		return fmt.Errorf("invalid repository structure: %w", err)
	}

	// Validate post-worktree-add script exists
	scriptPath := consts.GetFilePaths().PostWorktreeAddScript(repo.Alias)
	if err := validateScriptPath(scriptPath, repo.Dir); err != nil {
		return fmt.Errorf("invalid post-worktree-add script: %w", err)
	}

	return nil
}

func validateConfigurationHealth(appState *state.State) []error {
	var errors []error

	defaultReposDir := consts.GetDirectoryPaths().DefaultGitReposDir
	if _, err := os.Stat(defaultReposDir); os.IsNotExist(err) {
		errors = append(errors, fmt.Errorf("git repos directory does not exist: %s", defaultReposDir))
	}

	workOnScriptPath := consts.GetFilePaths().WorkOnScript
	if _, err := os.Stat(workOnScriptPath); os.IsNotExist(err) {
		errors = append(errors, fmt.Errorf("work-on script does not exist: %s", workOnScriptPath))
	}

	for _, repo := range appState.Repos {
		if err := validateRepositoryConfig(&repo); err != nil {
			errors = append(errors, fmt.Errorf("repository '%s': %w", repo.Alias, err))
		}
	}

	return errors
}
