package git

import (
)

type Remote struct {
	Name      string
	URL       []string
	Branch    *RemoteBranch
	Branches  []*RemoteBranch
}

func (entity *RepoEntity) NextRemote() error {
	currentRemoteIndex := 0
	for i, remote := range entity.Remotes {
		if remote.Name == entity.Remote.Name {
			currentRemoteIndex = i
		}
	}
	// WARNING: DIDN'T CHECK THE LIFE CYCLE
	if currentRemoteIndex == len(entity.Remotes)-1 {
		entity.Remote = entity.Remotes[0]
	} else {
		entity.Remote = entity.Remotes[currentRemoteIndex+1]
	}
	
	return nil
}

func (entity *RepoEntity) loadRemotes() error {
	r := entity.Repository
	entity.Remotes = make([]*Remote, 0)

	remotes, err := r.Remotes()
	for _, rm := range remotes {

			remote := &Remote{
				Name: rm.Config().Name,
				URL: rm.Config().URLs,
				}
			remote.loadRemoteBranches(&r)
			if len(remote.Branches) > 0 {
				remote.Branch = remote.Branches[0]
			}
			entity.Remotes = append(entity.Remotes, remote)
		
	}
	if err != nil {
		return err
	}
	return err
}