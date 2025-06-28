package consts

import (
	"fmt"
	"strings"
)

// GetWorkOnScriptContent returns the content for the work-on script
func GetWorkOnScriptContent() string {
	envVars := GetEnvironmentVariables()

	// Build environment variable documentation
	var envDocs strings.Builder
	envDocs.WriteString("# Available environment variables:\n")
	envDocs.WriteString(fmt.Sprintf("# - %s: %s\n", envVars.RepoAlias.Name, envVars.RepoAlias.Description))
	envDocs.WriteString(fmt.Sprintf("# - %s: %s\n", envVars.RepoDir.Name, envVars.RepoDir.Description))
	envDocs.WriteString(fmt.Sprintf("# - %s: %s\n", envVars.WorktreePath.Name, envVars.WorktreePath.Description))

	return fmt.Sprintf(`#!/bin/bash
# Work-on script
# This script is executed when working on a worktree
%s
`, envDocs.String())
}

// GetPostWorktreeAddScriptContent returns the content for the post-worktree-add script
func GetPostWorktreeAddScriptContent(repoAlias string) string {
	envVars := GetEnvironmentVariables()

	// Build environment variable documentation
	var envDocs strings.Builder
	envDocs.WriteString("# Available environment variables:\n")
	envDocs.WriteString(fmt.Sprintf("# - %s: %s\n", envVars.RepoAlias.Name, envVars.RepoAlias.Description))
	envDocs.WriteString(fmt.Sprintf("# - %s: %s\n", envVars.RepoDir.Name, envVars.RepoDir.Description))
	envDocs.WriteString(fmt.Sprintf("# - %s: %s\n", envVars.WorktreePath.Name, envVars.WorktreePath.Description))

	return fmt.Sprintf(`#!/bin/bash
# Post worktree add script for %s
# This script runs after a new worktree is created
%s
`, repoAlias, envDocs.String())
}
