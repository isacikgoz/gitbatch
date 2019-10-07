package load

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/isacikgoz/gitbatch/internal/git"
	"golang.org/x/sync/semaphore"
)

// AsyncAdd is interface to caller
type AsyncAdd func(r *git.Repository)

// SyncLoad initializes the go-git's repository objects with given
// slice of paths. since this job is done parallel, the order of the directories
// is not kept
func SyncLoad(directories []string) (entities []*git.Repository, err error) {
	entities = make([]*git.Repository, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, dir := range directories {
		// increment wait counter by one because we run a single goroutine
		// below
		wg.Add(1)
		go func(d string) {
			// decrement the wait counter by one, we call it in a defer so it's
			// called at the end of this goroutine
			defer wg.Done()
			entity, err := git.InitializeRepo(d)
			if err != nil {
				return
			}
			// lock so we don't get a race if multiple go routines try to add
			// to the same entities
			mu.Lock()
			entities = append(entities, entity)
			mu.Unlock()
		}(dir)
	}
	// wait until the wait counter is zero, this happens if all goroutines have
	// finished
	wg.Wait()
	if len(entities) == 0 {
		return entities, fmt.Errorf("there are no git repositories at given path(s)")
	}
	return entities, nil
}

// AsyncLoad asynchronously adds to AsyncAdd function
func AsyncLoad(directories []string, add AsyncAdd, d chan bool) error {
	ctx := context.TODO()

	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
	)

	var mx sync.Mutex

	// Compute the output using up to maxWorkers goroutines at a time.
	for _, dir := range directories {
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}

		go func(d string) {

			defer sem.Release(1)
			entity, err := git.InitializeRepo(d)
			if err != nil {
				return
			}
			// lock so we don't get a race if multiple go routines try to add
			// to the same entities
			mx.Lock()
			add(entity)
			mx.Unlock()
		}(dir)
	}
	// Acquire all of the tokens to wait for any remaining workers to finish.
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		return err
	}
	d <- true
	sem = nil
	return nil
}
