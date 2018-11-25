package app

import (
	"github.com/isacikgoz/gitbatch/pkg/gui"
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

	var err error

	app.Gui, err = gui.NewGui(directories)
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