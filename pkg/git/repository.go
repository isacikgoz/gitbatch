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
	Remote     *Remote
	Commit     *Commit
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
	
	commit, _ := lastCommit(r)
	entity = &RepoEntity{Name: fileInfo.Name(),
						AbsPath: directory,
						Repository: *r,
						Commit: commit,
						Marked: false,
		}
	entity.Branch = entity.GetActiveBranch()
	remotes, err := remoteBranches(r)
	if len(remotes) > 0 {
		entity.Remote = remotes[0]
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
	if err := entity.MergeWithGit(remote); err != nil {
		return err
	}
	return nil
}

func (entity *RepoEntity) PullTest() error {
	time.Sleep(5 * time.Second)
	return nil
}

func (entity *RepoEntity) Fetch() error {
	rm := entity.Remote.Reference.Name().Short()
	remote := strings.Split(rm, "/")[0]
	err := entity.Repository.Fetch(&git.FetchOptions{
		RemoteName: remote,
		})
	if err != nil {
		return err
	}
	return nil
}