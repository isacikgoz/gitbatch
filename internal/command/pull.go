package command

import (
	"os"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage"
	gerr "github.com/isacikgoz/gitbatch/internal/errors"
	"github.com/isacikgoz/gitbatch/internal/git"
)

var (
	pullTryCount int
	pullMaxTry   = 1
)

// PullOptions defines the rules for pull operation
type PullOptions struct {
	// Name of the remote to fetch from. Defaults to origin.
	RemoteName string
	// ReferenceName Remote branch to clone. If empty, uses HEAD.
	ReferenceName string
	// Fetch only ReferenceName if true.
	SingleBranch bool
	// Credentials holds the user and password information
	Credentials *git.Credentials
	// Process logs the output to stdout
	Progress bool
	// Force allows the pull to update a local branch even when the remote
	// branch does not descend from it.
	Force bool
	// Mode is the command mode
	CommandMode Mode
}

// Pull incorporates changes from a remote repository into the current branch.
func Pull(r *git.Repository, o *PullOptions) (err error) {
	pullTryCount = 0

	// here we configure pull operation
	switch o.CommandMode {
	case ModeLegacy:
		err = pullWithGit(r, o)
		return err
	case ModeNative:
		err = pullWithGoGit(r, o)
		return err
	}
	return nil
}

func pullWithGit(r *git.Repository, options *PullOptions) (err error) {
	args := make([]string, 0)
	args = append(args, "pull")
	// parse options to command line arguments
	if len(options.RemoteName) > 0 {
		args = append(args, options.RemoteName)
	}
	if options.Force {
		args = append(args, "-f")
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

func pullWithGoGit(r *git.Repository, options *PullOptions) (err error) {
	opt := &gogit.PullOptions{
		RemoteName:   options.RemoteName,
		SingleBranch: options.SingleBranch,
		Force:        options.Force,
	}
	if len(options.ReferenceName) > 0 {
		ref := plumbing.NewRemoteReferenceName(options.RemoteName, options.ReferenceName)
		opt.ReferenceName = ref
	}
	// if any credential is given, let's add it to the git.PullOptions
	if options.Credentials != nil {
		protocol, err := git.AuthProtocol(r.State.Remote)
		if err != nil {
			return err
		}
		if protocol == git.AuthProtocolHTTP || protocol == git.AuthProtocolHTTPS {
			opt.Auth = &http.BasicAuth{
				Username: options.Credentials.User,
				Password: options.Credentials.Password,
			}
		} else {
			return gerr.ErrInvalidAuthMethod
		}
	}
	if options.Progress {
		opt.Progress = os.Stdout
	}
	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}
	ref, _ := r.Repo.Head()
	if err = w.Pull(opt); err != nil {
		if err == gogit.NoErrAlreadyUpToDate {
			// log.Error("error: " + err.Error())
			// Already up-to-date
			// TODO: submit a PR for this kind of error, this type of catch is lame
		} else if err == storage.ErrReferenceHasChanged && pullTryCount < pullMaxTry {
			pullTryCount++
			if err := Fetch(r, &FetchOptions{
				RemoteName: options.RemoteName,
			}); err != nil {
				return err
			}
			return Pull(r, options)
		} else if strings.Contains(err.Error(), "SSH_AUTH_SOCK") {
			// The env variable SSH_AUTH_SOCK is not defined, maybe git can handle this
			return pullWithGit(r, options)
		} else if err == transport.ErrAuthenticationRequired {
			return gerr.ErrAuthenticationRequired
		} else {
			return pullWithGit(r, options)
		}
	}
	newref, _ := r.Repo.Head()

	msg, err := getMergeMessage(r, ref.Hash().String(), newref.Hash().String())
	if err != nil {
		msg = "couldn't get stat"
	}
	r.SetWorkStatus(git.Success)
	r.State.Message = msg
	return r.Refresh()
}
