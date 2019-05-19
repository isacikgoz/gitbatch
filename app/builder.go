package app

import (
	"errors"
	"os"

	"github.com/isacikgoz/gitbatch/gui"
	log "github.com/sirupsen/logrus"
)

// The App struct is responsible to hold app-wide related entities. Currently
// it has only the gui.Gui pointer for interface entity.
type App struct {
	Gui    *gui.Gui
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

// Setup will handle pre-required operations. It is designed to be a wrapper for
// main method right now.
func Setup(argConfig *Config) (*App, error) {
	// initiate the app and give it initial values
	app := &App{}
	if len(argConfig.Directories) <= 0 {
		d, _ := os.Getwd()
		argConfig.Directories = []string{d}
	}
	presetConfig, err := LoadConfiguration()
	if err != nil {
		return nil, err
	}
	appConfig := overrideConfig(presetConfig, argConfig)

	setLogLevel(appConfig.LogLevel)

	// hopefull everything went smooth as butter
	log.Trace("App configuration completed")

	dirs := generateDirectories(appConfig.Directories, appConfig.Depth)

	if appConfig.QuickMode {
		if err := execQuickMode(dirs, appConfig); err != nil {
			return nil, err
		}
		// we are done here and no need for an app to be configured
		return nil, nil
	}

	// create a gui.Gui struct and set it as App's gui
	app.Gui, err = gui.NewGui(appConfig.Mode, dirs)
	if err != nil {
		// the error types and handling is not considered yet
		return nil, err
	}
	return app, nil
}

// set the level of logging it is fatal by default
func setLogLevel(logLevel string) {
	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.FatalLevel)
	}
	log.WithFields(log.Fields{
		"level": logLevel,
	}).Trace("logging level has been set")
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

func execQuickMode(dirs []string, cfg *Config) error {
	x := cfg.Mode == "fetch"
	y := cfg.Mode == "pull"
	if x == y {
		return errors.New("unrecognized quick mode: " + cfg.Mode)
	}
	quick(dirs, cfg.Mode)
	return nil
}
