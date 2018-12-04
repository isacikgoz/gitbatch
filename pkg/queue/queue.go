package queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/isacikgoz/gitbatch/pkg/git"
)

// Job relates the type of the operation and the entity
type Job struct {
	JobType JobType
	Entity  *git.RepoEntity
}

// JobQueue holds the slice of Jobs
type JobQueue struct {
	series []*Job
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

// CreateJob es its name implies creates a job struct and return its pointer
func CreateJob() (j *Job, err error) {
	fmt.Println("Job created.")
	return j, nil
}

// starts the job
func (job *Job) start() error {
	job.Entity.State = git.Working
	// added for testing, TODO: remove
	time.Sleep(time.Second)
	// TODO: Handle errors?
	switch mode := job.JobType; mode {
	case Fetch:
		if err := job.Entity.Fetch(); err != nil {
			job.Entity.State = git.Fail
			return nil
		}
		job.Entity.RefreshPushPull()
		job.Entity.State = git.Success
	case Pull:
		if err := job.Entity.Pull(); err != nil {
			job.Entity.State = git.Fail
			return nil
		}
		job.Entity.RefreshPushPull()
		job.Entity.State = git.Success
	case Merge:
		if err := job.Entity.Merge(); err != nil {
			job.Entity.State = git.Fail
			return nil
		}
		job.Entity.RefreshPushPull()
		job.Entity.State = git.Success
	default:
		job.Entity.State = git.Available
		return nil
	}
	return nil
}

// CreateJobQueue creates a jobqueue struct and initialize its slice then return
// its pointer
func CreateJobQueue() (jobQueue *JobQueue) {
	s := make([]*Job, 0)
	return &JobQueue{
		series: s,
	}
}

// AddJob adds a job to the queue
func (jobQueue *JobQueue) AddJob(j *Job) error {
	for _, job := range jobQueue.series {
		if job.Entity.RepoID == j.Entity.RepoID && job.JobType == j.JobType {
			return errors.New("Same job already is in the queue")
		}
	}
	jobQueue.series = append(jobQueue.series, nil)
	copy(jobQueue.series[1:], jobQueue.series[0:])
	jobQueue.series[0] = j
	return nil
}

// StartNext starts the next job in the queue
func (jobQueue *JobQueue) StartNext() (j *Job, finished bool, err error) {
	finished = false
	if len(jobQueue.series) < 1 {
		finished = true
		return nil, finished, nil
	}
	i := len(jobQueue.series) - 1
	lastJob := jobQueue.series[i]
	jobQueue.series = jobQueue.series[:i]
	if err = lastJob.start(); err != nil {
		return lastJob, finished, err
	}
	return lastJob, finished, nil
}

// RemoveFromQueue deletes the given entity and its job from the queue
// TODO: it is not safe if the job has been started
func (jobQueue *JobQueue) RemoveFromQueue(entity *git.RepoEntity) error {
	removed := false
	for i, job := range jobQueue.series {
		if job.Entity.RepoID == entity.RepoID {
			jobQueue.series = append(jobQueue.series[:i], jobQueue.series[i+1:]...)
			removed = true
		}
	}
	if !removed {
		return errors.New("There is no job with given repoID")
	}
	return nil
}

// IsInTheQueue function; since the job and entity is not tied with its own
// struct, this function returns true if that entity is in the queue along with
// the jobs type
func (jobQueue *JobQueue) IsInTheQueue(entity *git.RepoEntity) (inTheQueue bool, jt JobType) {
	inTheQueue = false
	for _, job := range jobQueue.series {
		if job.Entity.RepoID == entity.RepoID {
			inTheQueue = true
			jt = job.JobType
		}
	}
	return inTheQueue, jt
}
