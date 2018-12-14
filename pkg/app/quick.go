package app

import (
	"sync"
	"time"

	"github.com/isacikgoz/gitbatch/pkg/git"
	log "github.com/sirupsen/logrus"
)

func quick(directories []string, depth int, mode string) {

	var wg sync.WaitGroup
	start := time.Now()
	for _, dir := range directories {
		wg.Add(1)
		go func(d string, mode string) {
			defer wg.Done()
			err := operate(d, mode)
			if err != nil {
				log.Errorf("%s: %s", d, err.Error())
			}
		}(dir, mode)
	}
	wg.Wait()
	elapsed := time.Since(start)
	log.Infof("%d repositories finished in: %s\n", len(directories), elapsed)
}

func operate(directory, mode string) error {
	r, err := git.FastInitializeRepo(directory)
	if err != nil {
		return err
	}
	switch mode {
	case "fetch":
		return git.Fetch(r, git.FetchOptions{
			RemoteName: "origin",
		})
	case "pull":
		return git.Pull(r, git.PullOptions{
			RemoteName: "origin",
		})
	}
	return nil
}
