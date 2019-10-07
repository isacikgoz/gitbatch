package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// generateDirectories returns poosible git repositories to pipe into git pkg's
// load function
func generateDirectories(dirs []string, depth int) []string {
	gitDirs := make([]string, 0)
	for i := 0; i <= depth; i++ {
		nonrepos, repos := walkRecursive(dirs, gitDirs)
		dirs = nonrepos
		gitDirs = repos
	}
	return gitDirs
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
func seperateDirectories(directory string) ([]string, []string, error) {
	dirs := make([]string, 0)
	gitDirs := make([]string, 0)
	files, err := ioutil.ReadDir(directory)
	// can we read the directory?
	if err != nil {
		return nil, nil, nil
	}
	for _, f := range files {
		repo := directory + string(os.PathSeparator) + f.Name()
		file, err := os.Open(repo)
		// if we cannot open it, simply continue to iteration and don't consider
		if err != nil {
			file.Close()
			continue
		}
		dir, err := filepath.Abs(file.Name())
		if err != nil {
			file.Close()
			continue
		}
		// with this approach, we ignore submodule or sub repositoreis in a git repository
		ff, err := os.Open(dir + string(os.PathSeparator) + ".git")
		if err != nil {
			dirs = append(dirs, dir)
		} else {
			gitDirs = append(gitDirs, dir)
		}
		ff.Close()
		file.Close()

	}
	return dirs, gitDirs, nil
}
