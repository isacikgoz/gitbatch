package main

import (
//	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	repo := "/home/isacikgoz/git/gitbatch/"
	argstr := []string{ "-C", repo, "status"}

	out, err := exec.Command("git", argstr...).Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println(string(out))
	
//	files, err := ioutil.ReadDir("/home/isacikgoz/git")
//	if err != nil {
//  		log.Fatal(err)
//	}
//
//	for _, f := range files {
//		log.Println(f.Name())
//	}
}

//func FindGitRepositories(directory string) []string, err {
//	[]string gitRepositories
//	files, err := ioutil.ReadDir(directory)
//		for _, f := range files {
//		log.Println(f.Name())
//	}
//	return files, err
//}

