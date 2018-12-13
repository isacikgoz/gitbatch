package git

import ()

// Job relates the type of the operation and the entity
type Job struct {
	// JobType is to select operation type that will be applied to repository
	JobType JobType
	// Entity points to the repository that will be used for operation
	Entity *RepoEntity
	// Options is a placeholder for operation options
	Options interface{}
}

// JobType is the a git operation supported
type JobType string

const (
	// Fetch is wrapper of git fetch command
	FetchJob JobType = "fetch"
	// Pull is wrapper of git pull command
	PullJob JobType = "pull"
	// Merge is wrapper of git merge command
	MergeJob JobType = "merge"
)

// starts the job
func (job *Job) start() error {
	job.Entity.State = Working
	// TODO: Handle errors?
	// TOOD: Better implementation required
	switch mode := job.JobType; mode {
	case FetchJob:
		var opts FetchOptions
		if job.Options != nil {
			opts = job.Options.(FetchOptions)
		} else {
			opts = FetchOptions{
				RemoteName: job.Entity.Remote.Name,
			}
		}
		if err := Fetch(job.Entity, opts); err != nil {
			job.Entity.State = Fail
			return err
		}
	case PullJob:
		var opts FetchOptions
		if job.Options != nil {
			opts = job.Options.(FetchOptions)
		} else {
			opts = FetchOptions{
				RemoteName: job.Entity.Remote.Name,
			}
		}
		if err := Fetch(job.Entity, opts); err != nil {
			job.Entity.State = Fail
			return err
		}
		if err := Merge(job.Entity, MergeOptions{
			BranchName: job.Entity.Remote.Branch.Name,
		}); err != nil {
			job.Entity.State = Fail
			return nil
		}
	case MergeJob:
		if err := Merge(job.Entity, MergeOptions{
			BranchName: job.Entity.Remote.Branch.Name,
		}); err != nil {
			job.Entity.State = Fail
			return nil
		}
	default:
		job.Entity.State = Available
		return nil
	}
	job.Entity.State = Success
	return nil
}
