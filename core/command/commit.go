package command

import (
	"errors"
	"time"

	"github.com/isacikgoz/gitbatch/core/git"
	log "github.com/sirupsen/logrus"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	commitCmdMode string

	commitCommand       = "commit"
	commitCmdModeLegacy = "git"
	commitCmdModeNative = "go-git"
)

// CommitOptions defines the rules for commit operation
type CommitOptions struct {
	// CommitMsg
	CommitMsg string
	// User
	User string
	// Email
	Email string
}

// CommitCommand defines which commit command to use.
func CommitCommand(r *git.Repository, options CommitOptions) (err error) {
	// here we configure commit operation
	// default mode is go-git (this may be configured)
	commitCmdMode = commitCmdModeNative

	switch commitCmdMode {
	case commitCmdModeLegacy:
		return commitWithGit(r, options)
	case commitCmdModeNative:
		return commitWithGoGit(r, options)
	}
	return errors.New("Unhandled commit operation")
}

// commitWithGit is simply a bare git commit -m <msg> command which is flexible
func commitWithGit(r *git.Repository, options CommitOptions) (err error) {
	args := make([]string, 0)
	args = append(args, commitCommand)
	args = append(args, "-m")
	// parse options to command line arguments
	if len(options.CommitMsg) > 0 {
		args = append(args, options.CommitMsg)
	}
	if err := GenericGitCommand(r.AbsPath, args); err != nil {
		log.Warn("Error at git command (commit)")
		r.Refresh()
		return err
	}
	// till this step everything should be ok
	return r.Refresh()
}

// commitWithGoGit is the primary commit method
func commitWithGoGit(r *git.Repository, options CommitOptions) (err error) {
	config, err := r.Repo.Config()
	if err != nil {
		return err
	}
	name := config.Raw.Section("user").Option("name")
	email := config.Raw.Section("user").Option("email")
	opt := &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	}

	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Commit(options.CommitMsg, opt)
	if err != nil {
		r.Refresh()
		return err
	}
	// till this step everything should be ok
	return r.Refresh()
}
