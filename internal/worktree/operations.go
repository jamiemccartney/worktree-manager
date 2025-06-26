package worktree

import (
	"fmt"
	"os"
	"path/filepath"

	"worktree-manager/internal/config"
	"worktree-manager/internal/context"
	"worktree-manager/internal/git"
	"worktree-manager/internal/output"
	"worktree-manager/internal/script"
	"worktree-manager/internal/validation"
)

func AddWorktree(cfg *config.Config, branch string) error {
	ctx, err := context.NewRepoContext(cfg)
	if err != nil {
		return err
	}

	worktreePath := ctx.GetWorktreePath(branch)

	if err := validation.ValidateWorktreeDoesNotExist(worktreePath, branch); err != nil {
		return err
	}

	if err := ctx.EnsureWorktreesDir(); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	err = ctx.WithRepoDir(func() error {
		output.Progress("Fetching from origin...")
		if err := git.FetchFromOrigin(ctx.CurrentRepo.Dir); err != nil {
			return fmt.Errorf("failed to fetch from origin: %w", err)
		}

		var sourceBranch string
		var message string

		if git.RemoteBranchExists(ctx.CurrentRepo.Dir, branch) {
			sourceBranch = fmt.Sprintf("origin/%s", branch)
			message = fmt.Sprintf("Worktree tracking remote branch '%s' created at %s", branch, worktreePath)
		} else {
			baseBranch, err := git.GetBaseBranch(ctx.CurrentRepo.Dir)
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

		if err := git.CreateWorktree(ctx.CurrentRepo.Dir, opts); err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}

		output.Success(message)
		return nil
	})

	if err != nil {
		return err
	}

	if ctx.CurrentRepo.PostWorktreeAddScript != "" {
		if err := script.ExecuteScript(ctx.CurrentRepo.PostWorktreeAddScript, ctx.CurrentRepo, worktreePath); err != nil {
			output.Warning("Post-worktree-add script failed: %v", err)
		}
	}

	if ctx.Config.AutomaticWorkOnAfterAdd {
		output.Progress("Running work-on logic...")
		if err := script.ExecuteWorkOnScript(ctx.Config, ctx.CurrentRepo, worktreePath); err != nil {
			output.Warning("Work-on script failed: %v", err)
		}
	}

	return nil
}

func RemoveWorktree(cfg *config.Config, branch string) error {
	ctx, err := context.NewRepoContext(cfg)
	if err != nil {
		return err
	}

	worktreePath := ctx.GetWorktreePath(branch)

	if err := validation.ValidateWorktreeExists(worktreePath, branch); err != nil {
		return err
	}

	err = ctx.WithRepoDir(func() error {
		if err := git.RemoveWorktree(ctx.CurrentRepo.Dir, worktreePath); err != nil {
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

func ListWorktrees(cfg *config.Config) error {
	ctx, err := context.NewRepoContext(cfg)
	if err != nil {
		return err
	}

	var worktrees []git.Worktree
	err = ctx.WithRepoDir(func() error {
		var err error
		worktrees, err = git.ListWorktrees(ctx.CurrentRepo.Dir)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	PrintWorktreeList(ctx.CurrentRepo.Alias, worktrees)
	return nil
}

func WorkOnWorktree(cfg *config.Config, branch string) error {
	ctx, err := context.NewRepoContext(cfg)
	if err != nil {
		return err
	}

	worktreePath := ctx.GetWorktreePath(branch)

	if err := validation.ValidateWorktreeExists(worktreePath, branch); err != nil {
		return fmt.Errorf("%w\n\n💡 Use 'wt tree add %s' to create it first", err, branch)
	}

	output.Progress("Working on branch '%s'...", branch)
	output.Info("Worktree path: %s", worktreePath)

	if err := script.ExecuteWorkOnScript(ctx.Config, ctx.CurrentRepo, worktreePath); err != nil {
		output.Warning("Work-on script failed: %v", err)
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
