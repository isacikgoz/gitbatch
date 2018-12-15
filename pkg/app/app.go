package app

import (
	"os"

	"github.com/isacikgoz/gitbatch/pkg/gui"
	log "github.com/sirupsen/logrus"
)

// The App struct is responsible to hold app-wide related entities. Currently
// it has only the gui.Gui pointer for interface entity.
type App struct {
	Gui    *gui.Gui
	Config *SetupConfig
}

// SetupConfig is an assembler data to initiate a setup
type SetupConfig struct {
	Directories []string
	LogLevel    string
	Depth       int
	QuickMode   bool
	Mode        string
}

// Setup will handle pre-required operations. It is designed to be a wrapper for
// main method right now.
func Setup(setupConfig *SetupConfig) (*App, error) {
	// initiate the app and give it initial values
	app := &App{}
	if len(setupConfig.Directories) <= 0 {
		d, _ := os.Getwd()
		setupConfig.Directories = []string{d}
	}

	appConfig, err := overrideDefaults(setupConfig)
	if err != nil {
		return nil, err
	}

	setLogLevel(appConfig.LogLevel)
	directories := generateDirectories(appConfig.Directories, appConfig.Depth)

	if appConfig.QuickMode {
		x := appConfig.Mode == "fetch"
		y := appConfig.Mode == "pull"
		if x == y {
			log.Fatal("Unrecognized quick mode: " + appConfig.Mode)
		}
		quick(directories, appConfig.Depth, appConfig.Mode)
		log.Fatal("Finished")
	}

	// create a gui.Gui struct and set it as App's gui
	app.Gui, err = gui.NewGui(appConfig.Mode, directories)
	if err != nil {
		// the error types and handling is not considered yer
		log.Error(err)
		return app, err
	}
	// hopefull everything went smooth as butter
	log.Trace("App configuration completed")
	return app, nil
}

// Close function will handle if any cleanup is required. e.g. closing streams
// or cleaning temproray files so on and so forth
func (app *App) Close() error {
	return nil
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

func overrideDefaults(setupConfig *SetupConfig) (appConfig *SetupConfig, err error) {
	appConfig, err = LoadConfiguration()
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
	return appConfig, err
}
