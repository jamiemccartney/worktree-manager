package output

import (
	"fmt"
	"os"
	"path/filepath"

	"worktree-manager/internal/config"
	"worktree-manager/internal/git"
)

// Success prints a success message with ✅ emoji
func Success(format string, args ...interface{}) {
	fmt.Printf("✅  "+format+"\n", args...)
}

// Error prints an error message with ❌ emoji
func Error(format string, args ...interface{}) {
	fmt.Printf("❌  "+format+"\n", args...)
}

// Progress prints a progress message with 🔄 emoji
func Progress(format string, args ...interface{}) {
	fmt.Printf("🔄  "+format+"\n", args...)
}

// Info prints an info message with 📁 emoji
func Info(format string, args ...interface{}) {
	fmt.Printf("📁  "+format+"\n", args...)
}

// Hint prints a helpful hint with 💡 emoji
func Hint(format string, args ...interface{}) {
	fmt.Printf("💡  "+format+"\n", args...)
}

// Warning prints a warning message with ⚠️ emoji
func Warning(format string, args ...interface{}) {
	fmt.Printf("⚠️  "+format+"\n", args...)
}

// Item prints a list item with 🔸 emoji
func Item(format string, args ...interface{}) {
	fmt.Printf("🔸  "+format+"\n", args...)
}

// Cleanup prints a cleanup message with 🗑️ emoji
func Cleanup(format string, args ...interface{}) {
	fmt.Printf("🗑️  "+format+"\n", args...)
}

// FormatRepoStatus formats repository status information
func FormatRepoStatus(repo *config.Repo) string {
	status := ""

	// Check if directory exists
	if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
		status = "❌  Directory does not exist"
	} else {
		status = "✅  Available"

		// Count worktrees if possible
		worktreesDir := filepath.Join(repo.Dir, "worktrees")
		if entries, err := os.ReadDir(worktreesDir); err == nil {
			status += fmt.Sprintf(" (%d worktrees)", len(entries))
		}
	}

	return status
}

// FormatWorktreeInfo formats worktree information for display
func FormatWorktreeInfo(wt git.Worktree) string {
	if wt.IsBare {
		return fmt.Sprintf("(bare) %s", wt.Path)
	}

	info := fmt.Sprintf("%s\n   Path: %s", filepath.Base(wt.Path), wt.Path)

	if wt.Branch != "" {
		info += fmt.Sprintf("\n   Branch: %s", wt.Branch)
	}

	if wt.HEAD != "" && len(wt.HEAD) >= 8 {
		info += fmt.Sprintf("\n   HEAD: %s", wt.HEAD[:8])
	}

	return info
}

// PrintRepoList prints a formatted list of repositories
func PrintRepoList(repos []config.Repo) {
	if len(repos) == 0 {
		fmt.Println("No repositories configured. Use 'wt repo clone <url>' to add one.")
		return
	}

	fmt.Printf("📁  Configured repositories (%d):\n\n", len(repos))

	for _, repo := range repos {
		Item(repo.Alias)
		fmt.Printf("   Directory: %s\n", repo.Dir)
		fmt.Printf("   Status: %s\n", FormatRepoStatus(&repo))

		if repo.PostWorktreeAddScript != "" {
			fmt.Printf("   Post-add script: %s\n", repo.PostWorktreeAddScript)
		}

		fmt.Println()
	}
}

// PrintWorktreeList prints a formatted list of worktrees
func PrintWorktreeList(repoAlias string, worktrees []git.Worktree) {
	fmt.Printf("📁 Worktrees for repository '%s':\n\n", repoAlias)

	if len(worktrees) == 0 {
		fmt.Println("No worktrees found. Use 'wt tree add <branch>' to create one.")
		return
	}

	for _, wt := range worktrees {
		Item(FormatWorktreeInfo(wt))
		fmt.Println()
	}
}
