package git

import (
	log "github.com/sirupsen/logrus"

	"context"
	"errors"
	"runtime"
	"sync"

	"golang.org/x/sync/semaphore"
)

type AsyncAdd func(e *RepoEntity)

// LoadRepositoryEntities initializes the go-git's repository obejcts with given
// slice of paths. since this job is done parallel, the order of the directories
// is not kept
func LoadRepositoryEntities(directories []string) (entities []*RepoEntity, err error) {
	entities = make([]*RepoEntity, 0)

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
			entity, err := InitializeRepo(d)
			if err != nil {
				log.WithFields(log.Fields{
					"directory": d,
				}).Trace("Cannot load git repository.")
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
		return entities, errors.New("There are no git repositories at given path(s)")
	}
	return entities, nil
}

func LoadRepositoryEntitiesAsync(directories []string, add AsyncAdd) error {
	ctx := context.TODO()

	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
	)

	var mx sync.Mutex
	for _, dir := range directories {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Errorf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(d string) {

			defer sem.Release(1)
			entity, err := InitializeRepo(d)
			if err != nil {
				log.WithFields(log.Fields{
					"directory": d,
				}).Trace("Cannot load git repository.")
				return
			}
			// lock so we don't get a race if multiple go routines try to add
			// to the same entities
			mx.Lock()
			add(entity)
			mx.Unlock()
		}(dir)
	}
	return nil
}
