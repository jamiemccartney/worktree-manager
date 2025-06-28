package executors

import (
	"testing"
	"worktree-manager/internal/consts"
	"worktree-manager/internal/state"
)

func TestBuildScriptEnvironment(t *testing.T) {
	repo := &state.Repo{
		Alias: "test-repo",
		Dir:   "/repo/dir",
	}

	env := buildScriptEnvironment(repo, "/worktree/path")

	// Check that our custom environment variables are present
	envVars := consts.GetEnvironmentVariables()
	found := map[string]bool{
		envVars.RepoAlias.Name + "=test-repo":         false,
		envVars.RepoDir.Name + "=/repo/dir":           false,
		envVars.WorktreePath.Name + "=/worktree/path": false,
	}

	for _, envVar := range env {
		if _, exists := found[envVar]; exists {
			found[envVar] = true
		}
	}

	for envVar, wasFound := range found {
		if !wasFound {
			t.Errorf("Expected environment variable %s was not found", envVar)
		}
	}
}
