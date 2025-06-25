package git

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	HEAD   string
	Branch string
	IsBare bool
}

// WorktreeCreateOptions contains options for creating a worktree
type WorktreeCreateOptions struct {
	Branch       string
	WorktreePath string
	SourceBranch string
	CreateBranch bool
}
