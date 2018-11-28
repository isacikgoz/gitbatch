package git

import (
	"errors"
	"os"
	"time"

	"github.com/isacikgoz/gitbatch/pkg/utils"
	"gopkg.in/src-d/go-git.v4"
)

type RepoEntity struct {
	RepoID     string
	Name       string
	AbsPath    string
	Repository git.Repository
	Branch     *Branch
	Branches   []*Branch
	Remote     *Remote
	Remotes    []*Remote
	Commit     *Commit
	Commits    []*Commit
	State      RepoState
}

type RepoState uint8

const (
	Available RepoState = 0 
	Queued    RepoState = 1
	Working   RepoState = 2
	Success   RepoState = 3
	Fail      RepoState = 4
)

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
	entity = &RepoEntity{RepoID: utils.RandomString(8),
		Name:       fileInfo.Name(),
		AbsPath:    directory,
		Repository: *r,
		State:     Available,
	}
	entity.loadLocalBranches()
	entity.loadCommits()
	if len(entity.Commits) > 0 {
		entity.Commit = entity.Commits[0]
	} else {
		return entity, errors.New("There is no commit for this repository: " + directory)
	}
	entity.loadRemotes()
	entity.Branch = entity.getActiveBranch()
	if len(entity.Remotes) > 0 {
		// TODO: tend to take origin/master as default
		entity.Remote = entity.Remotes[0]
		// TODO: same code on 3 different occasion, maybe something wrong?
		if err = entity.Remote.switchRemoteBranch(entity.Remote.Name + "/" + entity.Branch.Name); err !=nil {
		// probably couldn't find, but its ok.
		}
	} else {
		return entity, errors.New("There is no remote for this repository: " + directory)
	}
	return entity, nil
}

func (entity *RepoEntity) Pull() error {
	// TODO: Migrate this code to src-d/go-git
	// 2018-11-25: tried but it fails, will investigate.
	rm := entity.Remote.Name
	if err := entity.FetchWithGit(rm); err != nil {
		return err
	}
	entity.Checkout(entity.Branch)
	if err := entity.MergeWithGit(entity.Remote.Branch.Name); err != nil {
		entity.Refresh()
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
	rm := entity.Remote.Name
	if err := entity.FetchWithGit(rm); err != nil {
		return err
	}
	entity.Refresh()
	entity.Checkout(entity.Branch)
	return nil
}

func (entity *RepoEntity) Merge() error {
	entity.Checkout(entity.Branch)
	if err := entity.MergeWithGit(entity.Remote.Branch.Name); err != nil {
		entity.Refresh()
		return err
	}
	entity.Refresh()
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
	if err := entity.loadRemotes(); err != nil {
		return err
	}
	return nil
}
