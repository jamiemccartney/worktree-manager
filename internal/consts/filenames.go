package consts

type FileNameConstants struct {
	Config          string
	State           string
	WorkOnScript    string
	PostWorktreeAdd string
}

func GetFileNames() FileNameConstants {
	return FileNameConstants{
		Config:          "config.json",
		State:           "state.json",
		WorkOnScript:    "work-on.sh",
		PostWorktreeAdd: "post-worktree-add.sh",
	}
}
