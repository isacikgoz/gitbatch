package app

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/isacikgoz/gitbatch/pkg/gui"
)

// The App struct is responsible to hold app-wide related entities. Currently
// it has only the gui.Gui pointer for interface entity.
type App struct {
	Gui     *gui.Gui
}

// If any pre-required operation is needed, setup will handle that task. It is
// designed to be a wrapper for main method right now.
func Setup(directory string, repoPattern string) (*App, error) {
	// initiate the app and give it initial values
	app := &App{
	}
	var err error
	directories := generateDirectories(directory, repoPattern)

	// create a gui.Gui struct and set it as App's gui
	app.Gui, err = gui.NewGui(directories)
	if err != nil {
		// the error types and handling is not considered yer
		return app, err
	}
	// hopefull everything went smooth as butter
	return app, nil
}

// If any cleanup is required Close method with handle it. e.g. closing streams
// or cleaning temproray files so on and so forth
func (app *App) Close() error {
	return nil
}

// generateDirectories is to find all the files in given path. This method 
// does not check if the given file is a valid git repositories
func generateDirectories(directory string, repoPattern string) (directories []string) {
	files, err := ioutil.ReadDir(directory)

	// can we read the directory?
	if err != nil {
		log.Fatal(err)
	}

	// filter according to a pattern
	filteredFiles := filterDirectories(files, repoPattern)

	// now let's iterate over the our desired git directories
	for _, f := range filteredFiles {
		repo := directory + string(os.PathSeparator) + f.Name()
		file, err := os.Open(repo)

		// if we cannot open it, simply continue to iteration and don't consider
		if err != nil {
			continue
		}
		dir, err := filepath.Abs(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		// shaping our directory slice
		directories = append(directories, dir)
	}
	return directories
}

// takes a fileInfo slice and returns it with the ones matches with the 
// repoPattern string
func filterDirectories(files []os.FileInfo, repoPattern string) []os.FileInfo {
	var filteredRepos []os.FileInfo
	for _, f := range files {
		// it is just a simple filter
		if strings.Contains(f.Name(), repoPattern) {
			filteredRepos = append(filteredRepos, f)
		} else {
			continue
		}
	}
	return filteredRepos
}
