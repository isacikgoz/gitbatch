package git

import (
	"errors"
	"sync"
)

// JobQueue holds the slice of Jobs
type JobQueue struct {
	series []*Job
}

// CreateJobQueue creates a jobqueue struct and initialize its slice then return
// its pointer
func CreateJobQueue() (jq *JobQueue) {
	s := make([]*Job, 0)
	return &JobQueue{
		series: s,
	}
}

// AddJob adds a job to the queue
func (jq *JobQueue) AddJob(j *Job) error {
	for _, job := range jq.series {
		if job.Entity.RepoID == j.Entity.RepoID && job.JobType == j.JobType {
			return errors.New("Same job already is in the queue")
		}
	}
	jq.series = append(jq.series, nil)
	copy(jq.series[1:], jq.series[0:])
	jq.series[0] = j
	return nil
}

// StartNext starts the next job in the queue
func (jq *JobQueue) StartNext() (j *Job, finished bool, err error) {
	finished = false
	if len(jq.series) < 1 {
		finished = true
		return nil, finished, nil
	}
	i := len(jq.series) - 1
	lastJob := jq.series[i]
	jq.series = jq.series[:i]
	if err = lastJob.start(); err != nil {
		return lastJob, finished, err
	}
	return lastJob, finished, nil
}

// RemoveFromQueue deletes the given entity and its job from the queue
// TODO: it is not safe if the job has been started
func (jq *JobQueue) RemoveFromQueue(entity *RepoEntity) error {
	removed := false
	for i, job := range jq.series {
		if job.Entity.RepoID == entity.RepoID {
			jq.series = append(jq.series[:i], jq.series[i+1:]...)
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
func (jq *JobQueue) IsInTheQueue(entity *RepoEntity) (inTheQueue bool, j *Job) {
	inTheQueue = false
	for _, job := range jq.series {
		if job.Entity.RepoID == entity.RepoID {
			inTheQueue = true
			j = job
		}
	}
	return inTheQueue, j
}

// StartJobsAsync start he jobs in the queue asynchronously
func (jq *JobQueue) StartJobsAsync() map[*Job]error {
	fails := make(map[*Job]error)
	var wg sync.WaitGroup
	var mx sync.Mutex
	for range jq.series {
		wg.Add(1)
		go func() {
			defer wg.Done()
			j, _, err := jq.StartNext()
			if err != nil {
				mx.Lock()
				fails[j] = err
				mx.Unlock()
			}
		}()
	}
	wg.Wait()
	return fails
}
