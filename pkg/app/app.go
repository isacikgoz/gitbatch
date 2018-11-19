package app

import (
	"github.com/isacikgoz/gitbatch/pkg/gui"
	"github.com/isacikgoz/gitbatch/pkg/git"
	"io"
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

func createRepositoryEntities(directories []string) (entities []git.RepoEntity, err error) {
	for _, dir := range directories {
		entity, err := git.InitializeRepository(dir)
		if err != nil {
			continue
		}
		entities = append(entities, entity)
	}
	return entities, nil
}
