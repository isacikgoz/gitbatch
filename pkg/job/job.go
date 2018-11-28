package job

import (
	"errors"
	"fmt"
	"time"

	"github.com/isacikgoz/gitbatch/pkg/git"
)

type Job struct {
	JobType JobType
	Entity  *git.RepoEntity
}

type JobQueue struct {
	series []*Job
}

type JobType string

const (
	Fetch JobType = "fetch"
	Pull  JobType = "pull"
	Merge JobType = "merge"
)

func CreateJob() (j *Job, err error) {
	fmt.Println("Job created.")
	return j, nil
}

func (job *Job) start() error {
	time.Sleep(2*time.Second)
	job.Entity.State = git.Working
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

func CreateJobQueue() (jobQueue *JobQueue) {
	s := make([]*Job, 0)
	return &JobQueue{
		series: s,
	}
}

func (jobQueue *JobQueue) AddJob(j *Job) error {
	for _, job := range jobQueue.series {
		if job.Entity.RepoID == j.Entity.RepoID && job.JobType == j.JobType {
			return errors.New("Same job already is in the queue")
		}
	}
	jobQueue.series = append(jobQueue.series, j)
	return nil
}

func (jobQueue *JobQueue) StartNext() (j *Job, finished bool, err error) {
	finished = false
	if len(jobQueue.series) < 1 {
		finished = true
		return nil, finished, nil
	}
	i := len(jobQueue.series)-1
	lastJob := jobQueue.series[i]
	jobQueue.series = jobQueue.series[:i]
	if err = lastJob.start(); err != nil {
		return lastJob, finished, err
	}
	return lastJob, finished, nil
}

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
