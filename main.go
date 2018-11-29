package main

import (
	"log"
	"os"

	"github.com/isacikgoz/gitbatch/pkg/app"
	"gopkg.in/alecthomas/kingpin.v2"
)


var (
	// take this as default directory if user does not start app with -d flag
	currentDir, err = os.Getwd()
	dir             = kingpin.Flag("directory", "Directory to roam for git repositories.").Default(currentDir).Short('d').String()
	repoPattern     = kingpin.Flag("pattern", "Pattern to filter repositories").Short('p').String()
)

func main() {
	// parse the command line flag and options
	kingpin.Parse()

	// set the app
	app, err := app.Setup(*dir, *repoPattern)
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
