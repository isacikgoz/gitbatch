package app

import (
	"fmt"
	"os"

	"github.com/isacikgoz/gitbatch/internal/gui"
)

// The App struct is responsible to hold app-wide related entities. Currently
// it has only the gui.Gui pointer for interface entity.
type App struct {
	Config *Config
}

// Config is an assembler data to initiate a setup
type Config struct {
	Directories []string
	LogLevel    string
	Depth       int
	QuickMode   bool
	Mode        string
}

// New will handle pre-required operations. It is designed to be a wrapper for
// main method right now.
func New(argConfig *Config) (*App, error) {
	// initiate the app and give it initial values
	app := &App{}
	if len(argConfig.Directories) <= 0 {
		d, _ := os.Getwd()
		argConfig.Directories = []string{d}
	}
	presetConfig, err := loadConfiguration()
	if err != nil {
		return nil, err
	}
	app.Config = overrideConfig(presetConfig, argConfig)

	return app, nil
}

// Run starts the application.
func (a *App) Run() error {
	dirs := generateDirectories(a.Config.Directories, a.Config.Depth)
	if a.Config.QuickMode {
		return a.execQuickMode(dirs)
	}
	// create a gui.Gui struct and run the gui
	gui, err := gui.New(a.Config.Mode, dirs)
	if err != nil {
		return err
	}
	return gui.Run()
}

func overrideConfig(appConfig, setupConfig *Config) *Config {
	if len(setupConfig.Directories) > 0 {
		appConfig.Directories = setupConfig.Directories
	}
	if len(setupConfig.LogLevel) > 0 {
		appConfig.LogLevel = setupConfig.LogLevel
	}
	if setupConfig.Depth > 0 {
		appConfig.Depth = setupConfig.Depth
	}
	if setupConfig.QuickMode {
		appConfig.QuickMode = setupConfig.QuickMode
	}
	if len(setupConfig.Mode) > 0 {
		appConfig.Mode = setupConfig.Mode
	}
	return appConfig
}

func (a *App) execQuickMode(directories []string) error {
	if a.Config.Mode != "fetch" && a.Config.Mode != "pull" {
		return fmt.Errorf("unrecognized quick mode: " + a.Config.Mode)
	}

	return quick(directories, a.Config.Mode)
}
