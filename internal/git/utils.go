package git

import (
	"path"
	"strings"
)

// ExtractRepoNameFromURL extracts the repository name from a Git URL
// Supports both HTTPS and SSH formats:
// - https://github.com/user/repo.git -> repo
// - git@github.com:user/repo.git -> repo
// - https://github.com/user/repo -> repo
func ExtractRepoNameFromURL(url string) string {
	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")
	
	// Extract the last part of the path
	repoName := path.Base(url)
	
	// Handle SSH URLs like git@github.com:user/repo
	if strings.Contains(repoName, ":") {
		parts := strings.Split(repoName, ":")
		if len(parts) > 1 {
			repoName = path.Base(parts[len(parts)-1])
		}
	}
	
	return repoName
}