package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/isacikgoz/gitbatch/internal/app"
)

func main() {
	kingpin.Version("gitbatch version 0.6.1")

	dirs := kingpin.Flag("directory", "Directory(s) to roam for git repositories.").Short('d').Strings()
	mode := kingpin.Flag("mode", "Application start mode, more sensible with quick run.").Short('m').String()
	recursionDepth := kingpin.Flag("recursive-depth", "Find directories recursively.").Default("0").Short('r').Int()
	logLevel := kingpin.Flag("log-level", "Logging level; trace,debug,info,warn,error").Default("error").Short('l').String()
	quick := kingpin.Flag("quick", "runs without gui and fetches/pull remote upstream.").Short('q').Bool()

	kingpin.Parse()

	if err := run(*dirs, *logLevel, *recursionDepth, *quick, *mode); err != nil {
		fmt.Fprintf(os.Stderr, "application quitted with an unhandled error: %v", err)
		os.Exit(1)
	}
}

func run(dirs []string, log string, depth int, quick bool, mode string) error {
	app, err := app.New(&app.Config{
		Directories: dirs,
		LogLevel:    log,
		Depth:       depth,
		QuickMode:   quick,
		Mode:        mode,
	})
	if err != nil {
		return err
	}

	return app.Run()
}
