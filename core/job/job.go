package job

import (
	"github.com/isacikgoz/gitbatch/core/command"
	"github.com/isacikgoz/gitbatch/core/git"
)

// Job relates the type of the operation and the entity
type Job struct {
	// JobType is to select operation type that will be applied to repository
	JobType JobType
	// Entity points to the repository that will be used for operation
	Entity *git.RepoEntity
	// Options is a placeholder for operation options
	Options interface{}
}

// JobType is the a git operation supported
type JobType string

const (
	// FetchJob is wrapper of git fetch command
	FetchJob JobType = "fetch"

	// PullJob is wrapper of git pull command
	PullJob JobType = "pull"

	// MergeJob is wrapper of git merge command
	MergeJob JobType = "merge"
)

// starts the job
func (j *Job) start() error {
	j.Entity.SetState(git.Working)
	// TODO: Handle errors?
	// TOOD: Better implementation required
	switch mode := j.JobType; mode {
	case FetchJob:
		var opts command.FetchOptions
		if j.Options != nil {
			opts = j.Options.(command.FetchOptions)
		} else {
			opts = command.FetchOptions{
				RemoteName: j.Entity.Remote.Name,
			}
		}
		if err := command.Fetch(j.Entity, opts); err != nil {
			j.Entity.SetState(git.Fail)
			j.Entity.SetStateMessage(err.Error())
			return err
		}
	case PullJob:
		var opts command.PullOptions
		if j.Options != nil {
			opts = j.Options.(command.PullOptions)
		} else {
			opts = command.PullOptions{
				RemoteName: j.Entity.Remote.Name,
			}
		}
		if err := command.Pull(j.Entity, opts); err != nil {
			j.Entity.SetState(git.Fail)
			j.Entity.SetStateMessage(err.Error())
			return err
		}
	case MergeJob:
		if err := command.Merge(j.Entity, command.MergeOptions{
			BranchName: j.Entity.Remote.Branch.Name,
		}); err != nil {
			j.Entity.SetState(git.Fail)
			j.Entity.SetStateMessage(err.Error())
			return nil
		}
	default:
		j.Entity.SetState(git.Available)
		return nil
	}
	return nil
}
