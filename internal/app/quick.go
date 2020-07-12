package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/isacikgoz/gitbatch/internal/command"
	"github.com/isacikgoz/gitbatch/internal/git"
)

func quick(directories []string, mode string) error {
	var wg sync.WaitGroup
	start := time.Now()
	for _, dir := range directories {
		wg.Add(1)
		go func(d string, mode string) {
			defer wg.Done()
			if err := operate(d, mode); err != nil {
				fmt.Printf("could not perform %s on %s: %s", mode, d, err)
			}
			fmt.Printf("%s: successful\n", d)
		}(dir, mode)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("%d repositories finished in: %s\n", len(directories), elapsed)
	return nil
}

func operate(directory, mode string) error {
	r, err := git.FastInitializeRepo(directory)
	if err != nil {
		return err
	}
	switch mode {
	case "fetch":
		return command.Fetch(r, &command.FetchOptions{
			RemoteName: "origin",
			Progress:   true,
		})
	case "pull":
		return command.Pull(r, &command.PullOptions{
			RemoteName: "origin",
			Progress:   true,
		})
	}
	return nil
}
