package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// generateDirectories returns poosible git repositories to pipe into git pkg's
// load function
func generateDirectories(directories []string, depth int) (gitDirectories []string) {
	for i := 0; i <= depth; i++ {
		nonrepos, repos := walkRecursive(directories, gitDirectories)
		directories = nonrepos
		gitDirectories = repos
	}
	return gitDirectories
}

// returns given values, first search directories and second stands for possible
// git repositories. Call this func from a "for i := 0; i<depth; i++" loop
func walkRecursive(search, appendant []string) ([]string, []string) {
	max := len(search)
	for i := 0; i < max; i++ {
		if i >= len(search) {
			continue
		}
		// find possible repositories and remaining ones, b slice is possible ones
		a, b, err := seperateDirectories(search[i])
		if err != nil {
			log.WithFields(log.Fields{
				"directory": search[i],
			}).Trace("Can't read directory")
			continue
		}
		// since we started to search let's get rid of it and remove from search
		// array
		search[i] = search[len(search)-1]
		search = search[:len(search)-1]
		// lets append what we have found to continue recursion
		search = append(search, a...)
		appendant = append(appendant, b...)
	}
	return search, appendant
}

// seperateDirectories is to find all the files in given path. This method
// does not check if the given file is a valid git repositories
func seperateDirectories(directory string) (directories, gitDirectories []string, err error) {
	files, err := ioutil.ReadDir(directory)
	// can we read the directory?
	if err != nil {
		log.WithFields(log.Fields{
			"directory": directory,
		}).Trace("Can't read directory")
		return directories, gitDirectories, nil
	}
	for _, f := range files {
		repo := directory + string(os.PathSeparator) + f.Name()
		file, err := os.Open(repo)
		// if we cannot open it, simply continue to iteration and don't consider
		if err != nil {
			log.WithFields(log.Fields{
				"file":      file,
				"directory": repo,
			}).Trace("Failed to open file in the directory")
			continue
		}
		dir, err := filepath.Abs(file.Name())
		if err != nil {
			return nil, nil, err
		}
		// with this approach, we ignore submodule or sub repositoreis in a git repository
		_, err = os.Open(dir + string(os.PathSeparator) + ".git")
		if err != nil {
			directories = append(directories, dir)
		} else {
			gitDirectories = append(gitDirectories, dir)
		}
	}
	return directories, gitDirectories, nil
}

// takes a fileInfo slice and returns it with the ones matches with the
// pattern string
func filterDirectories(files []os.FileInfo, pattern string) []os.FileInfo {
	var filteredRepos []os.FileInfo
	for _, f := range files {
		// it is just a simple filter
		if strings.Contains(f.Name(), pattern) && f.Name() != ".git" {
			filteredRepos = append(filteredRepos, f)
		} else {
			continue
		}
	}
	return filteredRepos
}
