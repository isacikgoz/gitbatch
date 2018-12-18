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
func (j *Job) start() error {
	j.Entity.SetState(Working)
	// TODO: Handle errors?
	// TOOD: Better implementation required
	switch mode := j.JobType; mode {
	case FetchJob:
		var opts FetchOptions
		if j.Options != nil {
			opts = j.Options.(FetchOptions)
		} else {
			opts = FetchOptions{
				RemoteName: j.Entity.Remote.Name,
			}
		}
		if err := Fetch(j.Entity, opts); err != nil {
			j.Entity.SetState(Fail)
			return err
		}
	case PullJob:
		var opts PullOptions
		if j.Options != nil {
			opts = j.Options.(PullOptions)
		} else {
			opts = PullOptions{
				RemoteName: j.Entity.Remote.Name,
			}
		}
		if err := Pull(j.Entity, opts); err != nil {
			j.Entity.SetState(Fail)
			return err
		}
	case MergeJob:
		if err := Merge(j.Entity, MergeOptions{
			BranchName: j.Entity.Remote.Branch.Name,
		}); err != nil {
			j.Entity.SetState(Fail)
			return nil
		}
	default:
		j.Entity.SetState(Available)
		return nil
	}
	return nil
}
