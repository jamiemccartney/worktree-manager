package worktree

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"worktree-manager/internal/config"
	"worktree-manager/internal/consts"
	"worktree-manager/internal/executors"
	"worktree-manager/internal/fileops"
	"worktree-manager/internal/git"
	"worktree-manager/internal/output"
	"worktree-manager/internal/state"
)

var scriptExecutor executors.ScriptExecutor

func init() {
	scriptExecutor = executors.NewBashScriptExecutor()
}

// getWorktreePath computes the path for a worktree branch for a given repo
func getWorktreePath(repo *state.Repo, branch string) string {
	return filepath.Join(consts.GetDirectoryPaths().DefaultWorktreesDir, repo.Alias, branch)
}

// getWorktreesDir computes the worktrees directory for a given repo
func getWorktreesDir(repo *state.Repo) string {
	return filepath.Join(consts.GetDirectoryPaths().DefaultWorktreesDir, repo.Alias)
}

func AddWorktree(cfg *config.Config, appState *state.State, branch string) error {
	activeRepo, err := appState.GetActiveRepo()
	if err != nil {
		return fmt.Errorf("❌ %v", err)
	}

	worktreePath := getWorktreePath(activeRepo, branch)

	if err := validateWorktreeDoesNotExist(worktreePath, branch); err != nil {
		return err
	}

	if err := fileops.EnsureDir(getWorktreesDir(activeRepo)); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	err = fileops.WithDir(activeRepo.Dir, func() error {
		output.Progress("Fetching from origin...")

		if err := git.FetchFromOrigin(activeRepo.Dir); err != nil {
			return fmt.Errorf("failed to fetch from origin: %w", err)
		}

		var sourceBranch string
		var message string

		if git.RemoteBranchExists(activeRepo.Dir, branch) {
			sourceBranch = fmt.Sprintf("origin/%s", branch)
			message = fmt.Sprintf("Worktree tracking remote branch '%s' created at %s", branch, worktreePath)
		} else {
			baseBranch, err := git.GetBaseBranch(activeRepo.Dir)
			if err != nil {
				return fmt.Errorf("failed to determine base branch: %w", err)
			}
			sourceBranch = baseBranch
			message = fmt.Sprintf("New branch '%s' worktree created at %s from %s", branch, worktreePath, baseBranch)
		}

		opts := git.WorktreeCreateOptions{
			Branch:       branch,
			WorktreePath: worktreePath,
			SourceBranch: sourceBranch,
			CreateBranch: true,
		}

		if err := git.CreateWorktree(activeRepo.Dir, opts); err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}

		output.Success(message)
		return nil
	})

	if err != nil {
		return err
	}

	if err := scriptExecutor.Execute(&executors.ScriptExecutionContext{
		ScriptPath:   fileops.GetPostWorktreeAddScriptPath(activeRepo.Alias),
		Repo:         activeRepo,
		WorktreePath: worktreePath,
		WorkingDir:   consts.GetDirectoryPaths().RepoScriptsDir(activeRepo.Alias),
		ProgressMsg:  "Executing post-worktree-add script: %s",
	}); err != nil {
		output.Warning("Post-worktree-add script failed: %v", err)
	}

	if cfg.AutomaticWorkOnAfterAdd {
		output.Progress("Running work-on logic...")
		if err := scriptExecutor.Execute(&executors.ScriptExecutionContext{
			ScriptPath:   consts.GetFilePaths().WorkOnScript,
			Repo:         activeRepo,
			WorktreePath: worktreePath,
			WorkingDir:   worktreePath,
			ProgressMsg:  "Executing work-on script: %s",
		}); err != nil {
			output.Warning("Work-on script failed: %v", err)
		}
	}

	return nil
}

func RemoveWorktree(cfg *config.Config, appState *state.State, branch string) error {
	activeRepo, err := appState.GetActiveRepo()
	if err != nil {
		return fmt.Errorf("❌ %v", err)
	}

	worktreePath := getWorktreePath(activeRepo, branch)

	if err := validateWorktreeExists(worktreePath, branch); err != nil {
		return err
	}

	err = fileops.WithDir(activeRepo.Dir, func() error {
		if err := git.RemoveWorktree(activeRepo.Dir, worktreePath); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		if err := os.RemoveAll(worktreePath); err != nil {
			return fmt.Errorf("failed to delete folder: %w", err)
		}
		output.Cleanup("Deleted folder %s", worktreePath)
	}

	output.Success("Worktree '%s' removed", branch)
	return nil
}

func ListWorktrees(cfg *config.Config, appState *state.State) error {
	activeRepo, err := appState.GetActiveRepo()
	if err != nil {
		return fmt.Errorf("❌ %v", err)
	}

	var worktrees []git.Worktree
	err = fileops.WithDir(activeRepo.Dir, func() error {
		var err error
		worktrees, err = git.ListWorktrees(activeRepo.Dir)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	PrintWorktreeList(activeRepo.Alias, worktrees)
	return nil
}

func ListWorktreesJSON(cfg *config.Config, appState *state.State) error {
	activeRepo, err := appState.GetActiveRepo()
	if err != nil {
		return fmt.Errorf("❌ %v", err)
	}

	var worktrees []git.Worktree
	err = fileops.WithDir(activeRepo.Dir, func() error {
		var err error
		worktrees, err = git.ListWorktrees(activeRepo.Dir)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Extract just the branch names for autocompletion
	var branches []string
	for _, wt := range worktrees {
		if wt.Branch != "" {
			branches = append(branches, wt.Branch)
		}
	}

	jsonOutput, err := json.Marshal(branches)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonOutput))
	return nil
}

func WorkOnWorktree(cfg *config.Config, appState *state.State, branch string) error {
	activeRepo, err := appState.GetActiveRepo()
	if err != nil {
		return fmt.Errorf("❌ %v", err)
	}

	worktreePath := getWorktreePath(activeRepo, branch)

	if err := validateWorktreeExists(worktreePath, branch); err != nil {
		return fmt.Errorf("%w\n\n💡 Use 'wt tree add %s' to create it first", err, branch)
	}

	output.Progress("Working on branch '%s'...", branch)
	output.Info("Worktree path: %s", worktreePath)

	if err := scriptExecutor.Execute(&executors.ScriptExecutionContext{
		ScriptPath:   consts.GetFilePaths().WorkOnScript,
		Repo:         activeRepo,
		WorktreePath: worktreePath,
		WorkingDir:   worktreePath,
		ProgressMsg:  "Executing work-on script: %s",
	}); err != nil {
		output.Warning("Work-on script failed: %v", err)
	}

	return nil
}

func validateWorktreeExists(worktreePath, branch string) error {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree for branch '%s' does not exist at %s", branch, worktreePath)
	}
	return nil
}

func validateWorktreeDoesNotExist(worktreePath, branch string) error {
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		return fmt.Errorf("worktree path '%s' already exists", worktreePath)
	}
	return nil
}

func FormatWorktreeInfo(wt git.Worktree) string {
	info := fmt.Sprintf("%s\n   Path: %s", filepath.Base(wt.Path), wt.Path)

	if wt.Branch != "" {
		info += fmt.Sprintf("\n   Branch: %s", wt.Branch)
	}

	return info
}

func PrintWorktreeList(repoAlias string, worktrees []git.Worktree) {
	output.Info("Worktrees for repository '%s':", repoAlias)

	if len(worktrees) == 0 {
		output.Hint("No worktrees found. Use 'wt tree add <branch>' to create one.")
		return
	}

	for _, wt := range worktrees {
		output.Item(FormatWorktreeInfo(wt))
	}
}
