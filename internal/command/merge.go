package command

import (
	"regexp"

	gerr "github.com/isacikgoz/gitbatch/internal/errors"
	"github.com/isacikgoz/gitbatch/internal/git"
)

// MergeOptions defines the rules of a merge operation
type MergeOptions struct {
	// Name of the branch to merge with.
	BranchName string
	// Be verbose.
	Verbose bool
	// With true do not show a diffstat at the end of the merge.
	NoStat bool
	// Mode is the command mode
	CommandMode Mode
}

// Merge incorporates changes from the named commits or branches into the
// current branch
func Merge(r *git.Repository, options *MergeOptions) error {

	args := make([]string, 0)
	args = append(args, "merge")
	if len(options.BranchName) > 0 {
		args = append(args, options.BranchName)
	}
	if options.Verbose {
		args = append(args, "-v")
	}
	if options.NoStat {
		args = append(args, "-n")
	}

	ref, _ := r.Repo.Head()
	if out, err := Run(r.AbsPath, "git", args); err != nil {
		return gerr.ParseGitError(out, err)
	}

	newref, _ := r.Repo.Head()
	r.SetWorkStatus(git.Success)
	msg, err := getMergeMessage(r, ref.Hash().String(), newref.Hash().String())
	if err != nil {
		msg = "couldn't get stat"
	}
	r.State.Message = msg
	return r.Refresh()
}

func getMergeMessage(r *git.Repository, ref1, ref2 string) (string, error) {
	var msg string
	if ref1 == ref2 {
		msg = "already up-to-date"
	} else {
		out, err := DiffStatRefs(r, ref1, ref2)
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(`\r?\n`)
		lines := re.Split(out, -1)
		last := lines[len(lines)-1]
		if len(last) > 0 {
			msg = lines[len(lines)-1][1:]
		}
	}
	return msg, nil
}
