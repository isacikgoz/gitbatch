package job

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/thorstenhirsch/gitbatch/internal/git"
	"golang.org/x/sync/semaphore"
)

// Queue holds the slice of Jobs
type Queue struct {
	series []*Job
}

// CreateJobQueue creates a jobqueue struct and initialize its slice then return
// its pointer
func CreateJobQueue() (jq *Queue) {
	s := make([]*Job, 0)
	return &Queue{
		series: s,
	}
}

// AddJob adds a job to the queue
func (jq *Queue) AddJob(j *Job) error {
	for _, job := range jq.series {
		if job.Repository.RepoID == j.Repository.RepoID && job.JobType == j.JobType {
			return fmt.Errorf("same job already is in the queue")
		}
	}
	jq.series = append(jq.series, nil)
	copy(jq.series[1:], jq.series[0:])
	jq.series[0] = j
	return nil
}

// StartNext starts the next job in the queue
func (jq *Queue) StartNext() (j *Job, finished bool, err error) {
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
func (jq *Queue) RemoveFromQueue(r *git.Repository) error {
	removed := false
	for i, job := range jq.series {
		if job.Repository.RepoID == r.RepoID {
			jq.series = append(jq.series[:i], jq.series[i+1:]...)
			removed = true
		}
	}
	if !removed {
		return fmt.Errorf("there is no job with given repoID")
	}
	return nil
}

// IsInTheQueue function; since the job and entity is not tied with its own
// struct, this function returns true if that entity is in the queue along with
// the jobs type
func (jq *Queue) IsInTheQueue(r *git.Repository) (inTheQueue bool, j *Job) {
	inTheQueue = false
	for _, job := range jq.series {
		if job.Repository.RepoID == r.RepoID {
			inTheQueue = true
			j = job
		}
	}
	return inTheQueue, j
}

// StartJobsAsync start the jobs in the queue asynchronously
func (jq *Queue) StartJobsAsync() map[*Job]error {
	if len(jq.series) == 0 {
		return make(map[*Job]error)
	}

	// Use a context with better lifecycle management
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
		fails      = make(map[*Job]error)
	)

	var wg sync.WaitGroup
	var mx sync.Mutex

	// Process jobs with proper synchronization
	jobCount := len(jq.series)
	for i := 0; i < jobCount; i++ {
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}

		wg.Add(1)
		go func() {
			defer func() {
				sem.Release(1)
				wg.Done()
			}()

			j, finished, err := jq.StartNext()
			if finished {
				return
			}
			if err != nil {
				mx.Lock()
				fails[j] = err
				mx.Unlock()
			}
		}()
	}

	// Wait for all jobs to complete
	wg.Wait()
	return fails
}
