package command

import (
	gerr "github.com/isacikgoz/gitbatch/core/errors"
	"github.com/isacikgoz/gitbatch/core/git"
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
func Merge(r *git.Repository, options MergeOptions) error {
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
	if out, err := GenericGitCommandWithOutput(r.AbsPath, args); err != nil {
		return gerr.ParseGitError(out, err)
	}
	r.SetState(git.Success)
	return r.Refresh()
}
