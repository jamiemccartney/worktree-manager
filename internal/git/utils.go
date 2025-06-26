package git

import (
	"path"
	"strings"
)

func ExtractRepoNameFromURL(url string) string {
	url = strings.TrimSuffix(url, ".git")
	repoName := path.Base(url)
	if strings.Contains(repoName, ":") {
		parts := strings.Split(repoName, ":")
		if len(parts) > 1 {
			repoName = path.Base(parts[len(parts)-1])
		}
	}
	
	return repoName
}