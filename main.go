package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/isacikgoz/gitbatch/pkg/app"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	currentDir, err = os.Getwd()
	dir             = kingpin.Flag("directory", "Directory to roam for git repositories.").Default(currentDir).Short('d').String()
	repoPattern     = kingpin.Flag("pattern", "Pattern to filter repositories").Short('p').String()
)

func main() {
	kingpin.Parse()
	repositories := FindRepos(*dir)

	app, err := app.Setup(repositories)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Gui.Run()
	if err != nil {
		log.Fatal(err)
	}

	defer app.Close()
}

func FindRepos(directory string) (directories []string) {
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}

	filteredFiles := FilterRepos(files)
	for _, f := range filteredFiles {
		repo := directory + string(os.PathSeparator) + f.Name()
		file, err := os.Open(repo)
		if err != nil {
			continue
		}
		dir, err := filepath.Abs(file.Name())
		if err != nil {
			log.Fatal(err)
		}
		directories = append(directories, dir)
	}
	return directories
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
