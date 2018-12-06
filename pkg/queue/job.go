package queue

import (
	"time"

	"github.com/isacikgoz/gitbatch/pkg/git"
)

// Job relates the type of the operation and the entity
type Job struct {
	JobType JobType
	Entity  *git.RepoEntity
}

// JobType is the a git operation supported
type JobType string

const (
	// Fetch is wrapper of git fetch command
	Fetch JobType = "fetch"
	// Pull is wrapper of git pull command
	Pull JobType = "pull"
	// Merge is wrapper of git merge command
	Merge JobType = "merge"
)

// starts the job
func (job *Job) start() error {
	job.Entity.State = git.Working
	// added for testing, TODO: remove
	time.Sleep(time.Second)
	// TODO: Handle errors?
	switch mode := job.JobType; mode {
	case Fetch:
		if err := git.Fetch(job.Entity, git.FetchOptions{
			RemoteName: job.Entity.Remote.Name,
		}); err != nil {
			job.Entity.State = git.Fail
			return nil
		}
	case Pull:
		if err := git.Fetch(job.Entity, git.FetchOptions{
			RemoteName: job.Entity.Remote.Name,
		}); err != nil {
			job.Entity.State = git.Fail
			return nil
		}
		if err := git.Merge(job.Entity, git.MergeOptions{
			BranchName: job.Entity.Remote.Branch.Name,
		}); err != nil {
			job.Entity.State = git.Fail
			return nil
		}
	case Merge:
		if err := git.Merge(job.Entity, git.MergeOptions{
			BranchName: job.Entity.Remote.Branch.Name,
		}); err != nil {
			job.Entity.State = git.Fail
			return nil
		}
	default:
		job.Entity.State = git.Available
		return nil
	}
	job.Entity.State = git.Success
	job.Entity.Refresh()
	return nil
}
