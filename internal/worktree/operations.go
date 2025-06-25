package worktree

import (
	"fmt"
	"os"

	"worktree-manager/internal/context"
	"worktree-manager/internal/git"
	"worktree-manager/internal/output"
	"worktree-manager/internal/script"
	"worktree-manager/internal/validation"
)

// AddWorktree creates a new worktree for the specified branch
func AddWorktree(branch string) error {
	// Set up repository context
	ctx, err := context.NewRepoContext()
	if err != nil {
		return err
	}

	worktreePath := ctx.GetWorktreePath(branch)

	// Validate that worktree doesn't already exist
	if err := validation.ValidateWorktreeDoesNotExist(worktreePath, branch); err != nil {
		return err
	}

	// Ensure worktrees directory exists
	if err := ctx.EnsureWorktreesDir(); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	// Perform git operations in repository directory
	err = ctx.WithRepoDir(func() error {
		// Fetch from origin
		output.Progress("Fetching from origin...")
		if err := git.FetchFromOrigin(ctx.CurrentRepo.Dir); err != nil {
			return fmt.Errorf("failed to fetch from origin: %w", err)
		}

		// Determine source branch and create worktree
		var sourceBranch string
		var message string

		if git.RemoteBranchExists(ctx.CurrentRepo.Dir, branch) {
			// Track existing remote branch
			sourceBranch = fmt.Sprintf("origin/%s", branch)
			message = fmt.Sprintf("Worktree tracking remote branch '%s' created at %s", branch, worktreePath)
		} else {
			// Create new branch from base branch
			baseBranch, err := git.GetBaseBranch(ctx.CurrentRepo.Dir)
			if err != nil {
				return fmt.Errorf("failed to determine base branch: %w", err)
			}
			sourceBranch = baseBranch
			message = fmt.Sprintf("New branch '%s' worktree created at %s from %s", branch, worktreePath, baseBranch)
		}

		// Create the worktree
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

	// Execute post-worktree-add script if configured
	if ctx.CurrentRepo.PostWorktreeAddScript != "" {
		if err := script.ExecuteScript(ctx.CurrentRepo.PostWorktreeAddScript, ctx.CurrentRepo, worktreePath); err != nil {
			output.Warning("Post-worktree-add script failed: %v", err)
		}
	}

	// Automatically run work-on logic if configured
	if ctx.Config.AutomaticWorkOnAfterAdd {
		output.Progress("Running work-on logic...")
		if err := script.ExecuteWorkOnScript(ctx.Config, ctx.CurrentRepo, worktreePath); err != nil {
			output.Warning("Work-on script failed: %v", err)
		}
	}

	return nil
}

// RemoveWorktree removes a worktree for the specified branch
func RemoveWorktree(branch string) error {
	// Set up repository context
	ctx, err := context.NewRepoContext()
	if err != nil {
		return err
	}

	worktreePath := ctx.GetWorktreePath(branch)

	// Validate that worktree exists
	if err := validation.ValidateWorktreeExists(worktreePath, branch); err != nil {
		return err
	}

	// Perform git operations in repository directory
	err = ctx.WithRepoDir(func() error {
		// Remove the worktree
		if err := git.RemoveWorktree(ctx.CurrentRepo.Dir, worktreePath); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Check if directory still exists and remove it
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		if err := os.RemoveAll(worktreePath); err != nil {
			return fmt.Errorf("failed to delete folder: %w", err)
		}
		output.Cleanup("Deleted folder %s", worktreePath)
	}

	output.Success("Worktree '%s' removed", branch)
	return nil
}

// ListWorktrees lists all worktrees for the current repository
func ListWorktrees() error {
	// Set up repository context
	ctx, err := context.NewRepoContext()
	if err != nil {
		return err
	}

	// Get worktrees list
	var worktrees []git.Worktree
	err = ctx.WithRepoDir(func() error {
		var err error
		worktrees, err = git.ListWorktrees(ctx.CurrentRepo.Dir)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Display the list
	output.PrintWorktreeList(ctx.CurrentRepo.Alias, worktrees)
	return nil
}

// WorkOnWorktree changes to a worktree and runs work-on script
func WorkOnWorktree(branch string) error {
	// Set up repository context
	ctx, err := context.NewRepoContext()
	if err != nil {
		return err
	}

	worktreePath := ctx.GetWorktreePath(branch)

	// Validate that worktree exists
	if err := validation.ValidateWorktreeExists(worktreePath, branch); err != nil {
		return fmt.Errorf("%w\n\n💡 Use 'wt tree add %s' to create it first.", err, branch)
	}

	output.Progress("Working on branch '%s'...", branch)
	output.Info("Worktree path: %s", worktreePath)

	// Execute work-on script if configured
	if err := script.ExecuteWorkOnScript(ctx.Config, ctx.CurrentRepo, worktreePath); err != nil {
		output.Warning("Work-on script failed: %v", err)
	}

	return nil
}
