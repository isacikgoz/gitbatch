package git

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
)

type RepoEntity struct {
	Name       string
	AbsPath    string
	Repository git.Repository
	Pushables  string
	Pullables  string
	Branch     string
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
	branch, err := CurrentBranchName(directory)

	entity = &RepoEntity{fileInfo.Name(), directory, *r, pushable, pullable, branch, false, isClean(r, fileInfo.Name())}
	
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

func (entity *RepoEntity) UnMark() {
	entity.Marked = false
}

func (entity *RepoEntity) Pull() error {
	return nil
}


