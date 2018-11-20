package app

import (
	"github.com/isacikgoz/gitbatch/pkg/gui"
	"github.com/isacikgoz/gitbatch/pkg/git"
	"io"
	"sync"
)

// App struct
type App struct {
	closers []io.Closer
	Gui *gui.Gui
}

// Setup bootstrap a new application
func Setup(directories []string) (*App, error) {
	app := &App{
		closers: []io.Closer{},
	}

	entities, err := createRepositoryEntities(directories)
	if err != nil {
		return app, err
	}

	app.Gui, err = gui.NewGui(entities)
	if err != nil {
		return app, err
	}
	return app, nil
}

// Close closes any resources
func (app *App) Close() error {
	for _, closer := range app.closers {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func createRepositoryEntities(directories []string) (entities []*git.RepoEntity, err error) {
	entities = make([]*git.RepoEntity, 0)

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
			entity, err := git.InitializeRepository(d)
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

	return entities, nil
}