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
}

func InitializeRepository(directory string) (RepoEntity, error) {
	var entity RepoEntity
	file, err := os.Open(directory)
	if err != nil {
		return entity, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return entity, err
	}
	r, err := git.PlainOpen(directory)
	if err != nil {
		return entity, err
	}
	pushable, pullable := UpstreamDifferenceCount(directory)
	branch, err := CurrentBranchName(directory)
	entity = RepoEntity{fileInfo.Name(), directory, *r, pushable, pullable, branch}
	
	return entity, nil
}



