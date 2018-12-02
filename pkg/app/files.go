package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
			log.WithFields(log.Fields{
				"file":      file,
				"directory": directory,
			}).Trace("Failed to open file in the directory")
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
