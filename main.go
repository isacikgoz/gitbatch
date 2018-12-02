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
	dir             = kingpin.Flag("directory", "Directory to roam for git repositories.").Default(currentDir).Short('d').String()
	repoPattern     = kingpin.Flag("pattern", "Pattern to filter repositories").Short('p').String()
	logLevel        = kingpin.Flag("log-level", "Logging level; trace,debug,info,warn,error").Default("error").Short('l').String()
)

func main() {
	// parse the command line flag and options
	kingpin.Parse()

	// set the app
	app, err := app.Setup(*dir, *repoPattern, *logLevel)
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
