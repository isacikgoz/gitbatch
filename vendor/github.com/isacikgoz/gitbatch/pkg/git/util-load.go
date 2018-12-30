package git

import (
	log "github.com/sirupsen/logrus"

	"errors"
	"sync"
)

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
