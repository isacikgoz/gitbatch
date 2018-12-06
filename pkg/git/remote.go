package git

import (
	log "github.com/sirupsen/logrus"
)

// Remote struct is simply a collection of remote branches and wraps it with the
// name of the remote and fetch/push urls. It also holds the *selected* remote
// branch
type Remote struct {
	Name     string
	URL      []string
	Branch   *RemoteBranch
	Branches []*RemoteBranch
}

// NextRemote iterates over next branch of a remote
func (entity *RepoEntity) NextRemote() error {
	currentRemoteIndex := entity.findCurrentRemoteIndex()
	if currentRemoteIndex == len(entity.Remotes)-1 {
		entity.Remote = entity.Remotes[0]
	} else {
		entity.Remote = entity.Remotes[currentRemoteIndex+1]
	}
	// TODO: same code on 3 different occasion, maybe something wrong?
	if err := entity.Remote.switchRemoteBranch(entity.Remote.Name + "/" + entity.Branch.Name); err != nil {
		// probably couldn't find, but its ok.
	}
	return nil
}

// PreviousRemote iterates over previous branch of a remote
func (entity *RepoEntity) PreviousRemote() error {
	currentRemoteIndex := entity.findCurrentRemoteIndex()
	if currentRemoteIndex == 0 {
		entity.Remote = entity.Remotes[len(entity.Remotes)-1]
	} else {
		entity.Remote = entity.Remotes[currentRemoteIndex-1]
	}
	// TODO: same code on 3 different occasion, maybe something wrong?
	if err := entity.Remote.switchRemoteBranch(entity.Remote.Name + "/" + entity.Branch.Name); err != nil {
		// probably couldn't find, but its ok.
	}
	return nil
}

// returns the active remote index
func (entity *RepoEntity) findCurrentRemoteIndex() int {
	currentRemoteIndex := 0
	for i, remote := range entity.Remotes {
		if remote.Name == entity.Remote.Name {
			currentRemoteIndex = i
		}
	}
	return currentRemoteIndex
}

// search for remotes in go-git way. It is the short way to get remotes but it
// does not give any insght about remote branches
func (entity *RepoEntity) loadRemotes() error {
	r := entity.Repository
	entity.Remotes = make([]*Remote, 0)

	remotes, err := r.Remotes()
	for _, rm := range remotes {

		remote := &Remote{
			Name: rm.Config().Name,
			URL:  rm.Config().URLs,
		}
		remote.loadRemoteBranches(entity)
		if len(remote.Branches) > 0 {
			remote.Branch = remote.Branches[0]
		}
		entity.Remotes = append(entity.Remotes, remote)

	}
	if err != nil {
		log.Warn("Cannot find remotes " + err.Error())
		return err
	}
	return err
}
