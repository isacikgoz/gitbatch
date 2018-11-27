package job

import (
	"errors"
	"fmt"
)

type Job struct {
	JobType JobType
	RepoID  string
	Name    string
	Args    []string
}

type JobQueue struct {
	series []*Job
}

type JobType string

const (
	Fetch JobType = "FETCH"
	Pull  JobType = "PULL"
)

func CreateJob() (j *Job, err error) {
	fmt.Println("Job created.")
	return j, nil
}

func (job *Job) start() error {
	switch mode := job.JobType; mode {
	case Fetch:
		fmt.Println("Fetch operation is started")
	case Pull:
		fmt.Println("Pull operation is started")
	default:
		return errors.New("Unknown job type")
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
		if job.RepoID == j.RepoID && job.JobType == j.JobType {
			return errors.New("Same job already is in the queue")
		}
	}
	jobQueue.series = append(jobQueue.series, j)
	return nil
}

func (jobQueue *JobQueue) StartNext() error {
	lastJob := jobQueue.series[len(jobQueue.series)-1]
	if err := lastJob.start(); err != nil {
		return err
	}
	return nil
}

func (jobQueue *JobQueue) RemoveFromQueue(repoID string) error {
	removed := false
	for i, job := range jobQueue.series {
		if job.RepoID == repoID {
			jobQueue.series = append(jobQueue.series[:i], jobQueue.series[i+1:]...)
			removed = true
		}
	}
	if !removed {
		return errors.New("There is no job with given repoID")
	}
	return nil
}
