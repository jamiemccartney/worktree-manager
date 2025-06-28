package consts

// EnvironmentVariable represents a single environment variable with name and description
type EnvironmentVariable struct {
	Name        string
	Description string
}

// EnvironmentVariables represents all environment variables used by the application
type EnvironmentVariables struct {
	RepoAlias    EnvironmentVariable
	RepoDir      EnvironmentVariable
	WorktreePath EnvironmentVariable
}

// GetEnvironmentVariables returns all environment variables with names and descriptions
func GetEnvironmentVariables() EnvironmentVariables {
	return EnvironmentVariables{
		RepoAlias: EnvironmentVariable{
			Name:        "WT_REPO_ALIAS",
			Description: "The repository alias",
		},
		RepoDir: EnvironmentVariable{
			Name:        "WT_REPO_DIR",
			Description: "The repository directory",
		},
		WorktreePath: EnvironmentVariable{
			Name:        "WT_WORKTREE_PATH",
			Description: "The path to the worktree",
		},
	}
}
