package git

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
	"time"
	"strings"
)

type RepoEntity struct {
	Name       string
	AbsPath    string
	Repository git.Repository
	Pushables  string
	Pullables  string
	Branch     string
	Remote     *Remote
	Commit     *Commit
	Marked     bool
	Clean      bool
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
	pushable, pullable := UpstreamDifferenceCount(directory)
	headRef, err := r.Head()
	if err != nil {
		return nil, err
	}
	branch := headRef.Name().Short()
	remotes, err := remoteBranches(r)
	commit, _ := lastCommit(r)
	entity = &RepoEntity{fileInfo.Name(), directory, *r, pushable, pullable, branch, remotes[0], commit, false, isClean(r, fileInfo.Name())}
	return entity, nil
}

func isClean(r *git.Repository, name string) bool {
	w, err := r.Worktree()
	if err != nil {
		return false
	}
	// TODO: This function is incredibly slow
	s, err := w.Status()
	if err != nil {
		return false
	}
	return s.IsClean()
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
	if err := entity.MergeWithGit(entity.Branch, remote); err != nil {
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