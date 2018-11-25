package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"time"
)

type RepoEntity struct {
	Name       string
	AbsPath    string
	Repository git.Repository
	Pushables  string
	Pullables  string
	Branch     string
	Remote     string
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
	w, err := entity.Repository.Worktree()
	if err != nil {
		return err
	}
	rf := plumbing.NewBranchReferenceName(entity.Branch)
	rm := entity.Remote
	err = w.Pull(&git.PullOptions{
		RemoteName: rm,
		ReferenceName: rf,
		})
	if err != nil {
		return err
	}

	return nil
}

func (entity *RepoEntity) PullTest() error {
	time.Sleep(5 * time.Second)

	return nil
}

func (entity *RepoEntity) Fetch() error {
	err := entity.Repository.Fetch(&git.FetchOptions{
		RemoteName: entity.Remote,
		})
	if err != nil {
		return err
	}

	return nil
}

func (entity *RepoEntity) GetActiveRemote() string {
	if list, err := entity.Repository.Remotes(); err != nil {
        return ""
    } else {
        for _, r := range list {
        	return r.Config().Name
        }
    }
    return ""
}