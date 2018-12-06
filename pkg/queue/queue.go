package queue

import (
	"errors"

	"github.com/isacikgoz/gitbatch/pkg/git"
)

// JobQueue holds the slice of Jobs
type JobQueue struct {
	series []*Job
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
