package git

import (
	"gopkg.in/src-d/go-git.v4"
)

func getRemotes(r *git.Repository) (remotes []string, err error) {

    if list, err := r.Remotes(); err != nil {
        return remotes, err
    } else {
        for _, r := range list {
        	remoteString := r.Config().Name
            remotes = append(remotes, remoteString)
        }
    }
    return remotes, nil
}

func (entity *RepoEntity) NextRemote() (string, error) {

	remotes, err := getRemotes(&entity.Repository)
	if err != nil {
		return entity.Remote, err
	}

	currentRemoteIndex := 0
	for i, remote := range remotes {
		if remote == entity.Remote {
			currentRemoteIndex = i
		}
	}

	if currentRemoteIndex == len(remotes)-1 {
		entity.Remote = remotes[0]
	} else {
		entity.Remote = remotes[currentRemoteIndex+1]
	}
	
	return entity.Remote, nil
}