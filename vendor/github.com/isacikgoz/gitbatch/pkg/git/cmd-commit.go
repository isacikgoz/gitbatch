package git

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
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
func CommitCommand(e *RepoEntity, options CommitOptions) (err error) {
	// here we configure commit operation
	// default mode is go-git (this may be configured)
	commitCmdMode = commitCmdModeNative

	switch commitCmdMode {
	case commitCmdModeLegacy:
		return commitWithGit(e, options)
	case commitCmdModeNative:
		return commitWithGoGit(e, options)
	}
	return errors.New("Unhandled commit operation")
}

// commitWithGit is simply a bare git commit -m <msg> command which is flexible
func commitWithGit(e *RepoEntity, options CommitOptions) (err error) {
	args := make([]string, 0)
	args = append(args, commitCommand)
	args = append(args, "-m")
	// parse options to command line arguments
	if len(options.CommitMsg) > 0 {
		args = append(args, options.CommitMsg)
	}
	if err := GenericGitCommand(e.AbsPath, args); err != nil {
		log.Warn("Error at git command (commit)")
		e.Refresh()
		return err
	}
	// till this step everything should be ok
	return e.Refresh()
}

// commitWithGoGit is the primary commit method
func commitWithGoGit(e *RepoEntity, options CommitOptions) (err error) {
	config, err := e.Repository.Config()
	if err != nil {
		return err
	}
	name := config.Raw.Section("user").Option("name")
	email := config.Raw.Section("user").Option("email")
	opt := &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	}

	w, err := e.Repository.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Commit(options.CommitMsg, opt)
	if err != nil {
		e.Refresh()
		return err
	}
	// till this step everything should be ok
	return e.Refresh()
}
