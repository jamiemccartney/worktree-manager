package git

type Worktree struct {
	Path   string
	Branch string
}

type WorktreeCreateOptions struct {
	Branch       string
	WorktreePath string
	SourceBranch string
	CreateBranch bool
}
