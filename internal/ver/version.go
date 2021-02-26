package ver

var (
	// Build runVersion
	version   = "n/a"
	branch    = "n/a"
	commit    = "n/a"
	buildTime = "n/a"
)

func GetVersion() string {
	return version
}

func GetBranch() string {
	return branch
}

func GetCommit() string {
	return commit
}

func GetBuildTime() string {
	return buildTime
}
