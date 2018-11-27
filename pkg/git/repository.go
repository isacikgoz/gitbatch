package git

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
	"time"
	"strings"
	"errors"
)

type RepoEntity struct {
	Name       string
	AbsPath    string
	Repository git.Repository
	Branch     *Branch
	Branches   []*Branch
	Remote     *Remote
	Remotes    []*Remote
	Commit     *Commit
	Commits    []*Commit
	Marked     bool
}

func InitializeRepository(directory string) (entity *RepoEntity, err error) {
	file, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}
	entity = &RepoEntity{Name: fileInfo.Name(),
						AbsPath: directory,
						Repository: *r,
						Marked: false,
		}
	entity.loadLocalBranches()
	entity.loadCommits()
	if len(entity.Commits) > 0 {
		entity.Commit = entity.Commits[0]
	} else {
		return entity, errors.New("There is no commit for this repository: " + directory)
	}
	entity.loadRemoteBranches()
	entity.Branch = entity.GetActiveBranch()
	if len(entity.Remotes) > 0 {
		// TODO: tend to take origin/master as default
		entity.Remote = entity.Remotes[0]
	} else {
		return entity, errors.New("There is no remote for this repository: " + directory)
	}
	return entity, nil
}

func (entity *RepoEntity) Mark() {
	entity.Marked = true
}

func (entity *RepoEntity) Unmark() {
	entity.Marked = false
}

func (entity *RepoEntity) Pull() error {
	// TODO: Migrate this code to src-d/go-git
	// 2018-11-25: tried but it fails, will investigate.
	rm := entity.Remote.Reference.Name().Short()
	remote := strings.Split(rm, "/")[0]
	if err := entity.FetchWithGit(remote); err != nil {
		return err
	}
	if err := entity.MergeWithGit(rm); err != nil {
		return err
	}
	entity.Refresh()
	entity.Checkout(entity.Branch)
	return nil
}

func (entity *RepoEntity) PullTest() error {
	time.Sleep(5 * time.Second)
	return nil
}

func (entity *RepoEntity) Fetch() error {
	rm := entity.Remote.Reference.Name().Short()
	remote := strings.Split(rm, "/")[0]
	if err := entity.FetchWithGit(remote); err != nil {
		return err
	}
	entity.Refresh()
	entity.Checkout(entity.Branch)
	// err := entity.Repository.Fetch(&git.FetchOptions{
	// 	RemoteName: remote,
	// 	})
	// if err != nil {
	// 	return err
	// }
	return nil
}
func (entity *RepoEntity) Refresh() error {
	r, err := git.PlainOpen(entity.AbsPath)
	if err != nil {
		return err
	}
	entity.Repository = *r
	if err := entity.loadLocalBranches(); err != nil {
		return err
	}
	if err := entity.loadCommits(); err != nil {
		return err
	}
	if err := entity.loadRemoteBranches(); err != nil {
		return err
	}
	return nil
}