package main

import (
	"github.com/isacikgoz/gitbatch/pkg/app"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	dirs         = kingpin.Flag("directory", "Directory(s) to roam for git repositories.").Short('d').Strings()
	mode         = kingpin.Flag("mode", "Application start mode, more sensible with quick run.").Short('m').String()
	recurseDepth = kingpin.Flag("recursive-depth", "Find directories recursively.").Default("0").Short('r').Int()
	logLevel     = kingpin.Flag("log-level", "Logging level; trace,debug,info,warn,error").Default("error").Short('l').String()
	quick        = kingpin.Flag("quick", "runs without gui and fetches/pull remote upstream.").Short('q').Bool()
)

func main() {
	kingpin.Version("gitbatch version 0.2.2")
	// parse the command line flag and options
	kingpin.Parse()

	// set the app
	app, err := app.Setup(&app.SetupConfig{
		Directories: *dirs,
		LogLevel:    *logLevel,
		Depth:       *recurseDepth,
		QuickMode:   *quick,
		Mode:        *mode,
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
