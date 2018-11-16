package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"os"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	currentDir, err = os.Getwd()
	dir = kingpin.Flag("directory", "Directory to roam for git repositories.").Default(currentDir).Short('d').String()
	repoPattern = kingpin.Flag("pattern", "Pattern to filter repositories").Short('p').String()
	branch = kingpin.Flag("branch", "branch to be pulled").Default("master").Short('b').String()
	remote = kingpin.Flag("remote", "remote name te be pulled").Default("origin").Short('r').String()
)

func main() {
	kingpin.Parse()
	log.Printf("%s is your repo pattern", *repoPattern)
	FindRepos(*dir)
}

func FindRepos(directory string) []string {
	var gitRepositories []string
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}
	filteredFiles := FilterRepos(files)
	for _, f := range filteredFiles {
		repo := directory + "/" + f.Name()
		r, err := git.PlainOpen(repo)
		if err !=nil {
			continue
		}
		// Get the working directory for the repository
		w, err := r.Worktree()
		CheckIfError(err)

		ref := plumbing.ReferenceName("refs/heads/" + *branch)
		err = w.Pull(&git.PullOptions{
			RemoteName: *remote,
			Progress: os.Stdout,
			ReferenceName: ref,
		})
		CheckIfError(err)
	}
	return gitRepositories
}

func FilterRepos(files []os.FileInfo) []os.FileInfo {
	var filteredRepos []os.FileInfo
	for _, f := range files {
		if strings.Contains(f.Name(), *repoPattern) {
			filteredRepos = append(filteredRepos, f)
		} else {
			continue
		}
	}
	return filteredRepos
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
