package git

import "testing"

func TestExtractRepoNameFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "HTTPS URL with .git suffix",
			url:      "https://github.com/user/repo.git",
			expected: "repo",
		},
		{
			name:     "HTTPS URL without .git suffix",
			url:      "https://github.com/user/repo",
			expected: "repo",
		},
		{
			name:     "SSH URL with .git suffix",
			url:      "git@github.com:user/repo.git",
			expected: "repo",
		},
		{
			name:     "SSH URL without .git suffix",
			url:      "git@github.com:user/repo",
			expected: "repo",
		},
		{
			name:     "GitLab HTTPS URL",
			url:      "https://gitlab.com/user/my-project.git",
			expected: "my-project",
		},
		{
			name:     "BitBucket SSH URL",
			url:      "git@bitbucket.org:user/my-repo.git",
			expected: "my-repo",
		},
		{
			name:     "Complex repo name with dashes",
			url:      "https://github.com/org/worktree-manager.git",
			expected: "worktree-manager",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractRepoNameFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("ExtractRepoNameFromURL(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}
