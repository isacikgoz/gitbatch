package command

import (
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	giterr "github.com/isacikgoz/gitbatch/internal/errors"
	"github.com/isacikgoz/gitbatch/internal/git"
)

// CommitOptions defines the rules for commit operation
type CommitOptions struct {
	// CommitMsg
	CommitMsg string
	// User
	User string
	// Email
	Email string
	// Mode is the command mode
	CommandMode Mode
}

// Commit defines which commit command to use.
func Commit(r *git.Repository, o *CommitOptions) (err error) {
	// here we configure commit operation

	switch o.CommandMode {
	case ModeLegacy:
		return commitWithGit(r, o)
	case ModeNative:
		return commitWithGoGit(r, o)
	}
	return fmt.Errorf("unhandled commit operation")
}

// commitWithGit is simply a bare git commit -m <msg> command which is flexible
func commitWithGit(r *git.Repository, opt *CommitOptions) (err error) {
	args := make([]string, 0)
	args = append(args, "commit")
	args = append(args, "-m")
	// parse options to command line arguments
	if len(opt.CommitMsg) > 0 {
		args = append(args, opt.CommitMsg)
	}
	if out, err := Run(r.AbsPath, "git", args); err != nil {
		_ = r.Refresh()
		return giterr.ParseGitError(out, err)
	}
	// till this step everything should be ok
	return r.Refresh()
}

// commitWithGoGit is the primary commit method
func commitWithGoGit(r *git.Repository, options *CommitOptions) (err error) {
	opt := &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  options.User,
			Email: options.Email,
			When:  time.Now(),
		},
		Committer: &object.Signature{
			Name:  options.User,
			Email: options.Email,
			When:  time.Now(),
		},
	}

	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Commit(options.CommitMsg, opt)
	if err != nil {
		_ = r.Refresh()
		return err
	}
	// till this step everything should be ok
	return r.Refresh()
}
