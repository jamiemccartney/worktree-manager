package consts

import "path/filepath"

type FilePathConstants struct {
	Config          string
	State           string
	WorkOnScript    string
	PostWorktreeAdd func(string) string
}

func GetFilePaths() FilePathConstants {
	fileNames := GetFileNames()
	directoryPaths := GetDirectoryPaths()

	return FilePathConstants{
		Config:       filepath.Join(directoryPaths.WorktreeManagerDir, fileNames.Config),
		State:        filepath.Join(directoryPaths.WorktreeManagerDir, fileNames.State),
		WorkOnScript: filepath.Join(directoryPaths.ScriptsDir, fileNames.WorkOnScript),
		PostWorktreeAdd: func(repo string) string {
			return filepath.Join(directoryPaths.ScriptsDir, repo, "post-worktree-add.sh")
		},
	}
}
