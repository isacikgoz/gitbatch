package app

import (
	"github.com/isacikgoz/gitbatch/pkg/gui"
	"github.com/isacikgoz/gitbatch/pkg/git"
	"io"
)

// App struct
type App struct {
	closers []io.Closer
}

// Setup bootstrap a new application
func Setup(repositories []git.RepoEntity) (*App, error) {
	app := &App{
		closers: []io.Closer{},
	}

	var err error

	err = gui.Run(repositories)
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
