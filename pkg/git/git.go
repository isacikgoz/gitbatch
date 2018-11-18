package git

import (
	"gopkg.in/src-d/go-git.v4"
)

type RepoEntity struct {
	Name       string
	Repository git.Repository
	Pushables  string
	Pullables  string
	Branch     string
}

func InitializeRepository(directory string) (RepoEntity, error) {
	var entity RepoEntity

	r, err := git.PlainOpen(directory)
	if err != nil {
		return entity, err
	}
	entity = RepoEntity{directory, *r, "", "", ""}
	
	return entity, nil
}

func InitializeRepositories(directories []string) []RepoEntity {
	var gitRepositories []RepoEntity
	for _, f := range directories {
		r, err := git.PlainOpen(f)
		if err != nil {
			continue
		}
		entity := RepoEntity{f, *r, "", "", ""}
		gitRepositories = append(gitRepositories, entity)
	}
	return gitRepositories
}



