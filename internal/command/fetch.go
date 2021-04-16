package command

import (
	"os"
	"regexp"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gerr "github.com/isacikgoz/gitbatch/internal/errors"
	"github.com/isacikgoz/gitbatch/internal/git"
)

var (
	fetchTryCount int
	fetchMaxTry   = 1
)

// FetchOptions defines the rules for fetch operation
type FetchOptions struct {
	// Name of the remote to fetch from. Defaults to origin.
	RemoteName string
	// Credentials holds the user and password information
	Credentials *git.Credentials
	// Before fetching, remove any remote-tracking references that no longer
	// exist on the remote.
	Prune bool
	// Show what would be done, without making any changes.
	DryRun bool
	// Process logs the output to stdout
	Progress bool
	// Force allows the fetch to update a local branch even when the remote
	// branch does not descend from it.
	Force bool
	// Mode is the command mode
	CommandMode Mode
	// There should be more room for authentication, tags and progress
}

// Fetch branches refs from one or more other repositories, along with the
// objects necessary to complete their histories
func Fetch(r *git.Repository, o *FetchOptions) (err error) {
	// here we configure fetch operation
	// default mode is go-git (this may be configured)
	mode := o.CommandMode
	fetchTryCount = 0
	// prune and dry run is not supported from go-git yet, rely on old friend
	if o.Prune || o.DryRun {
		mode = ModeLegacy
	}
	switch mode {
	case ModeLegacy:
		err = fetchWithGit(r, o)
		return err
	case ModeNative:
		// this should be the refspec as default, let's give it a try
		// TODO: Fix for quick mode, maybe better read config file
		var refspec string
		if r.State.Branch == nil {
			refspec = "+refs/heads/*:refs/remotes/origin/*"
		} else {
			refspec = "+" + "refs/heads/" + r.State.Branch.Name + ":" + "/refs/remotes/" + r.State.Remote.Name + "/" + r.State.Branch.Name
		}
		err = fetchWithGoGit(r, o, refspec)
		return err
	}
	return nil
}

// fetchWithGit is simply a bare git fetch <remote> command which is flexible
// for complex operations, but on the other hand, it ties the app to another
// tool. To avoid that, using native implementation is preferred.
func fetchWithGit(r *git.Repository, options *FetchOptions) (err error) {
	args := make([]string, 0)
	args = append(args, "fetch")
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
	if out, err := Run(r.AbsPath, "git", args); err != nil {
		return gerr.ParseGitError(out, err)
	}
	r.SetWorkStatus(git.Success)
	r.State.Message = ""
	// till this step everything should be ok
	return r.Refresh()
}

// fetchWithGoGit is the primary fetch method and refspec is the main feature.
// RefSpec is a mapping from local branches to remote references The format of
// the refspec is an optional +, followed by <src>:<dst>, where <src> is the
// pattern for references on the remote side and <dst> is where those references
// will be written locally. The + tells Git to update the reference even if it
// isn’t a fast-forward.
func fetchWithGoGit(r *git.Repository, options *FetchOptions, refspec string) (err error) {
	opt := &gogit.FetchOptions{
		RemoteName: options.RemoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec(refspec)},
		Force:      options.Force,
	}
	// if any credential is given, let's add it to the git.FetchOptions
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
	if err := r.Repo.Fetch(opt); err != nil {
		if err == gogit.NoErrAlreadyUpToDate {
			// Already up-to-date
			// TODO: submit a PR for this kind of error, this type of catch is lame
		} else if strings.Contains(err.Error(), "couldn't find remote ref") {
			// we don't have remote ref, so lets pull other things.. maybe it'd be useful
			rp := r.State.Remote.RefSpecs[0]
			if fetchTryCount < fetchMaxTry {
				fetchTryCount++
				_ = fetchWithGoGit(r, options, rp)
			} else {
				return err
			}
			// TODO: submit a PR for this kind of error, this type of catch is lame
		} else if strings.Contains(err.Error(), "SSH_AUTH_SOCK") {
			// The env variable SSH_AUTH_SOCK is not defined, maybe git can handle this
			return fetchWithGit(r, options)
		} else if err == transport.ErrAuthenticationRequired {
			return gerr.ErrAuthenticationRequired
		} else {
			return fetchWithGit(r, options)
		}
	}
	r.SetWorkStatus(git.Success)

	ref, _ := r.Repo.Head()
	// TODO: fix this, refresh two times not cool
	_ = r.Refresh()
	uRef := "origin/HEAD"
	if r.State.Branch != nil && r.State.Branch.Upstream != nil {
		uRef = r.State.Branch.Upstream.Reference.Hash().String()[:7]
	}

	msg, err := getFetchMessage(r, ref.Hash().String()[:7], uRef)
	if err != nil {
		msg = "couldn't get stat"
	}
	r.State.Message = msg
	// till this step everything should be ok
	return r.Refresh()
}

func getFetchMessage(r *git.Repository, ref1, ref2 string) (string, error) {
	msg := ref1 + ".." + ref2 + " "
	if ref1 == ref2 {
		msg = msg + "already up-to-date"
	} else {
		out, err := DiffStatRefs(r, ref1, ref2)
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(`\r?\n`)
		lines := re.Split(out, -1)
		last := lines[len(lines)-1]
		if len(last) > 0 {
			changes := strings.Split(last, ",")
			msg = msg + changes[0][1:]
		}
	}
	return msg, nil
}
