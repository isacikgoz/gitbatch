package job

import (
	"github.com/isacikgoz/gitbatch/internal/command"
	"github.com/isacikgoz/gitbatch/internal/git"
)

// Job relates the type of the operation and the entity
type Job struct {
	// JobType is to select operation type that will be applied to repository
	JobType Type
	// Repository points to the repository that will be used for operation
	Repository *git.Repository
	// Options is a placeholder for operation options
	Options interface{}
}

// Type is the a git operation supported
type Type string

const (
	// FetchJob is wrapper of git fetch command
	FetchJob Type = "fetch"

	// PullJob is wrapper of git pull command
	PullJob Type = "pull"

	// MergeJob is wrapper of git merge command
	MergeJob Type = "merge"

	// CheckoutJob is wrapper of git merge command
	CheckoutJob Type = "checkout"
)

// starts the job
func (j *Job) start() error {
	j.Repository.SetWorkStatus(git.Working)
	// TODO: Better implementation required
	switch mode := j.JobType; mode {
	case FetchJob:
		j.Repository.State.Message = "fetching.."
		var opts *command.FetchOptions
		if j.Options != nil {
			opts = j.Options.(*command.FetchOptions)
		} else {
			opts = &command.FetchOptions{
				RemoteName:  j.Repository.State.Remote.Name,
				CommandMode: command.ModeNative,
			}
		}
		if err := command.Fetch(j.Repository, opts); err != nil {
			j.Repository.SetWorkStatus(git.Fail)
			j.Repository.State.Message = err.Error()
			return err
		}
	case PullJob:
		j.Repository.State.Message = "pulling.."
		var opts *command.PullOptions
		if j.Repository.State.Branch.Upstream == nil {
			j.Repository.SetWorkStatus(git.Fail)
			j.Repository.State.Message = "upstream not set"
			return nil
		}
		if j.Options != nil {
			opts = j.Options.(*command.PullOptions)
		} else {
			opts = &command.PullOptions{
				RemoteName:  j.Repository.State.Remote.Name,
				CommandMode: command.ModeNative,
			}
		}
		if err := command.Pull(j.Repository, opts); err != nil {
			j.Repository.SetWorkStatus(git.Fail)
			j.Repository.State.Message = err.Error()
			return err
		}
	case MergeJob:
		j.Repository.State.Message = "merging.."
		if j.Repository.State.Branch.Upstream == nil {
			j.Repository.SetWorkStatus(git.Fail)
			j.Repository.State.Message = "upstream not set"
			return nil
		}
		if err := command.Merge(j.Repository, &command.MergeOptions{
			BranchName: j.Repository.State.Branch.Upstream.Name,
		}); err != nil {
			j.Repository.SetWorkStatus(git.Fail)
			j.Repository.State.Message = err.Error()
			return err
		}
	case CheckoutJob:
		j.Repository.State.Message = "switching to.."
		var opts *command.CheckoutOptions
		if j.Options != nil {
			opts = j.Options.(*command.CheckoutOptions)
		} else {
			opts = &command.CheckoutOptions{
				TargetRef:   "master",
				CommandMode: command.ModeNative,
			}
		}
		if err := command.Checkout(j.Repository, opts); err != nil {
			j.Repository.SetWorkStatus(git.Fail)
			j.Repository.State.Message = err.Error()
			return err
		}
	default:
		j.Repository.SetWorkStatus(git.Available)
		return nil
	}
	return nil
}
