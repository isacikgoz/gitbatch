package load

import (
	"context"
	"runtime"
	"sync"

	"github.com/isacikgoz/gitbatch/internal/git"
	"golang.org/x/sync/semaphore"
)

// AsyncAdd is interface to caller
type AsyncAdd func(r *git.Repository)

// AsyncLoad initializes repositories in parallel, passing each one to add as it
// is loaded. It returns once every directory has been tried.
func AsyncLoad(directories []string, add AsyncAdd) error {
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
	return sem.Acquire(ctx, int64(maxWorkers))
}
