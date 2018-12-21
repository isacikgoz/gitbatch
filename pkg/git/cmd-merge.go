package git

import (
	log "github.com/sirupsen/logrus"
)

var mergeCommand = "merge"

// MergeOptions defines the rules of a merge operation
type MergeOptions struct {
	// Name of the branch to merge with.
	BranchName string
	// Be verbose.
	Verbose bool
	// With true do not show a diffstat at the end of the merge.
	NoStat bool
}

// Merge incorporates changes from the named commits or branches into the
// current branch
func Merge(e *RepoEntity, options MergeOptions) error {
	args := make([]string, 0)
	args = append(args, mergeCommand)
	if len(options.BranchName) > 0 {
		args = append(args, options.BranchName)
	}
	if options.Verbose {
		args = append(args, "-v")
	}
	if options.NoStat {
		args = append(args, "-n")
	}
	if err := GenericGitCommand(e.AbsPath, args); err != nil {
		log.Warn("Error while merging")
		return err
	}
	e.SetState(Success)
	return e.Refresh()
}
