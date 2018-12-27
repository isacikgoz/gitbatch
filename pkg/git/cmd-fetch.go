package git

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var (
	fetchCmdMode  string
	fetchTryCount int

	fetchCommand       = "fetch"
	fetchCmdModeLegacy = "git"
	fetchCmdModeNative = "go-git"
	fetchMaxTry        = 1
)

// FetchOptions defines the rules for fetch operation
type FetchOptions struct {
	// Name of the remote to fetch from. Defaults to origin.
	RemoteName string
	// Credentials holds the user and pswd information
	Credentials Credentials
	// Before fetching, remove any remote-tracking references that no longer
	// exist on the remote.
	Prune bool
	// Show what would be done, without making any changes.
	DryRun bool
	// Force allows the fetch to update a local branch even when the remote
	// branch does not descend from it.
	Force bool
	// There should be more room for authentication, tags and progress
}

// Fetch branches refs from one or more other repositories, along with the
// objects necessary to complete their histories
func Fetch(e *RepoEntity, options FetchOptions) (err error) {
	// here we configure fetch operation
	// default mode is go-git (this may be configured)
	fetchCmdMode = fetchCmdModeNative
	fetchTryCount = 0
	// prune and dry run is not supported from go-git yet, rely on old friend
	if options.Prune || options.DryRun {
		fetchCmdMode = fetchCmdModeLegacy
	}
	switch fetchCmdMode {
	case fetchCmdModeLegacy:
		err = fetchWithGit(e, options)
		return err
	case fetchCmdModeNative:
		// this should be the refspec as default, let's give it a try
		// TODO: Fix for quick mode, maybe better read config file
		var refspec string
		if e.Branch == nil {
			refspec = "+refs/heads/*:refs/remotes/origin/*"
		} else {
			refspec = "+" + "refs/heads/" + e.Branch.Name + ":" + "/refs/remotes/" + e.Remote.Branch.Name
		}
		err = fetchWithGoGit(e, options, refspec)
		return err
	}
	return nil
}

// fetchWithGit is simply a bare git fetch <remote> command which is flexible
// for complex operations, but on the other hand, it ties the app to another
// tool. To avoid that, using native implementation is preferred.
func fetchWithGit(e *RepoEntity, options FetchOptions) (err error) {
	args := make([]string, 0)
	args = append(args, fetchCommand)
	// parse options to command line arguments
	if len(options.RemoteName) > 0 {
		args = append(args, options.RemoteName)
	}
	if options.Prune {
		args = append(args, "-p")
	}
	if options.Force {
		args = append(args, "-f")
	}
	if options.DryRun {
		args = append(args, "--dry-run")
	}
	if err := GenericGitCommand(e.AbsPath, args); err != nil {
		log.Warn("Error at git command (fetch)")
		return err
	}
	e.SetState(Success)
	// till this step everything should be ok
	return e.Refresh()
}

// fetchWithGoGit is the primary fetch method and refspec is the main feature.
// RefSpec is a mapping from local branches to remote references The format of
// the refspec is an optional +, followed by <src>:<dst>, where <src> is the
// pattern for references on the remote side and <dst> is where those references
// will be written locally. The + tells Git to update the reference even if it
// isn’t a fast-forward.
func fetchWithGoGit(e *RepoEntity, options FetchOptions, refspec string) (err error) {
	opt := &git.FetchOptions{
		RemoteName: options.RemoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec(refspec)},
		Force:      options.Force,
	}
	// if any credential is given, let's add it to the git.FetchOptions
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

	err = e.Repository.Fetch(opt)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			// Already up-to-date
			log.Warn(err.Error())
			// TODO: submit a PR for this kind of error, this type of catch is lame
		} else if strings.Contains(err.Error(), "couldn't find remote ref") {
			// we dont have remote ref, so lets pull other things.. maybe it'd be useful
			rp := e.Remote.RefSpecs[0]
			if fetchTryCount < fetchMaxTry {
				fetchTryCount++
				fetchWithGoGit(e, options, rp)
			} else {
				return err
			}
		} else if err == transport.ErrAuthenticationRequired {
			log.Warn(err.Error())
			return ErrAuthenticationRequired
		} else {
			log.Warn(err.Error())
			return err
		}
	}

	e.SetState(Success)
	// till this step everything should be ok
	return e.Refresh()
}
