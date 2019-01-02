package git

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var (
	pullCmdMode  string
	pullTryCount int

	pullCommand       = "pull"
	pullCmdModeLegacy = "git"
	pullCmdModeNative = "go-git"
	pullMaxTry        = 1
)

// PullOptions defines the rules for pull operation
type PullOptions struct {
	// Name of the remote to fetch from. Defaults to origin.
	RemoteName string
	// ReferenceName Remote branch to clone. If empty, uses HEAD.
	ReferenceName string
	// Fetch only ReferenceName if true.
	SingleBranch bool
	// Credentials holds the user and pswd information
	Credentials Credentials
	// Process logs the output to stdout
	Progress bool
	// Force allows the pull to update a local branch even when the remote
	// branch does not descend from it.
	Force bool
}

// Pull ncorporates changes from a remote repository into the current branch.
func Pull(e *RepoEntity, options PullOptions) (err error) {
	// here we configure pull operation
	// default mode is go-git (this may be configured)
	pullCmdMode = pullCmdModeNative
	pullTryCount = 0

	switch pullCmdMode {
	case pullCmdModeLegacy:
		err = pullWithGit(e, options)
		return err
	case pullCmdModeNative:
		err = pullWithGoGit(e, options)
		return err
	}
	return nil
}

func pullWithGit(e *RepoEntity, options PullOptions) (err error) {
	args := make([]string, 0)
	args = append(args, pullCommand)
	// parse options to command line arguments
	if len(options.RemoteName) > 0 {
		args = append(args, options.RemoteName)
	}
	if options.Force {
		args = append(args, "-f")
	}
	if err := GenericGitCommand(e.AbsPath, args); err != nil {
		log.Warn("Error at git command (pull)")
		return err
	}
	e.SetState(Success)
	return e.Refresh()
}

func pullWithGoGit(e *RepoEntity, options PullOptions) (err error) {
	opt := &git.PullOptions{
		RemoteName:   options.RemoteName,
		SingleBranch: options.SingleBranch,
		Force:        options.Force,
	}
	if len(options.ReferenceName) > 0 {
		ref := plumbing.NewRemoteReferenceName(options.RemoteName, options.ReferenceName)
		opt.ReferenceName = ref
	}
	// if any credential is given, let's add it to the git.PullOptions
	if len(options.Credentials.User) > 0 {
		protocol, err := authProtocol(e.Remote)
		if err != nil {
			return err
		}
		if protocol == authProtocolHTTP || protocol == authProtocolHTTPS {
			opt.Auth = &http.BasicAuth{
				Username: options.Credentials.User,
				Password: options.Credentials.Password,
			}
		} else {
			return ErrInvalidAuthMethod
		}
	}
	if options.Progress {
		opt.Progress = os.Stdout
	}
	w, err := e.Repository.Worktree()
	if err != nil {
		return err
	}

	if err = w.Pull(opt); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			// log.Error("error: " + err.Error())
			// Already up-to-date
			log.Warn(err.Error())
			// TODO: submit a PR for this kind of error, this type of catch is lame
		} else if err == memory.ErrRefHasChanged && pullTryCount < pullMaxTry {
			pullTryCount++
			log.Error("trying to fetch")
			if err := Fetch(e, FetchOptions{
				RemoteName: options.RemoteName,
			}); err != nil {
				return err
			}
			return Pull(e, options)
		} else if strings.Contains(err.Error(), "SSH_AUTH_SOCK") {
			// The env variable SSH_AUTH_SOCK is not defined, maybe git can handle this
			return pullWithGit(e, options)
		} else if err == transport.ErrAuthenticationRequired {
			log.Warn(err.Error())
			return ErrAuthenticationRequired
		} else {
			log.Warn(err.Error())
			return pullWithGit(e, options)
		}
	}

	e.SetState(Success)
	return e.Refresh()
}
