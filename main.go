package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
//	"gopkg.in/alecthomas/kingpin.v2"
)

//var (
//	currentDirectory, err = os.Getwd()
//	directory = kingpin.Flag("directory", "Directory to roam for git repositories.").Default(currentDirectory).Short('d').String()
//)

func main() {

	args := os.Args[1:]
	if len(args) > 0 {
		repo := args[0]
		FindRepos(repo)
	} else {
		repo, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		FindRepos(repo)
	}
}

func FindRepos(directory string) []string {
	var gitRepositories []string
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if IsRepo((directory + "/" + f.Name())) == true {
			log.Println(f.Name() + " is a git repository")
		}
	}
	return gitRepositories
}

func IsRepo(directory string) bool {

	argstr := []string{ "-C", directory, "status"}

	cmd := exec.Command("git", argstr...)

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// Do nothing
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf(directory + " is not a git repository, git command return code: %d", status.ExitStatus())
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	} else {
		return true
	}
	return false
}

