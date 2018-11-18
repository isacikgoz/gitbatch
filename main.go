package main

import (
	"github.com/isacikgoz/gitbatch/pkg/app"
	"github.com/isacikgoz/gitbatch/pkg/git"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	currentDir, err    = os.Getwd()
	dir                = kingpin.Flag("directory", "Directory to roam for git repositories.").Default(currentDir).Short('d').String()
	repoPattern        = kingpin.Flag("pattern", "Pattern to filter repositories").Short('p').String()
	repositories       []git.RepoEntity
)

func main() {
	kingpin.Parse()
	repositories = FindRepos(*dir)

	app, err := app.Setup(repositories) 
	if err != nil {
		log.Fatal(err)
	}

	defer app.Close()
}

func FindRepos(directory string) []git.RepoEntity {
	var gitRepositories []git.RepoEntity
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}

	filteredFiles := FilterRepos(files)
	for _, f := range filteredFiles {
		repo := directory + "/" + f.Name()

		entity, err := git.InitializeRepository(repo)
		if err != nil {
			continue
		}
		gitRepositories = append(gitRepositories, entity)
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


