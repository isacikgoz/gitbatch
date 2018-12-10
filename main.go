package main

import (
	"os"

	"github.com/isacikgoz/gitbatch/pkg/app"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// take this as default directory if user does not start app with -d flag
	currentDir, err = os.Getwd()
	dirs            = kingpin.Flag("directory", "Directory to roam for git repositories").Default(currentDir).Short('d').Strings()
	ignoreConfig    = kingpin.Flag("ignore-config", "Ignore config file").Short('i').Bool()
	recurseDepth    = kingpin.Flag("recursive-depth", "Find directories recursively").Default("1").Short('r').Int()
	logLevel        = kingpin.Flag("log-level", "Logging level; trace,debug,info,warn,error").Default("error").Short('l').String()
)

func main() {
	kingpin.Version("gitbatch version 0.1.0 (alpha)")
	// parse the command line flag and options
	kingpin.Parse()

	// set the app
	app, err := app.Setup(app.SetupConfig{
		Directories:  *dirs,
		LogLevel:     *logLevel,
		IgnoreConfig: *ignoreConfig,
		Depth:        *recurseDepth,
	})
	if err != nil {
		log.Fatal(err)
	}

	// execute the app and wait its routine
	err = app.Gui.Run()
	if err != nil {
		log.Fatal(err)
	}

	// good citizens always clean up their mess
	defer app.Close()
}
