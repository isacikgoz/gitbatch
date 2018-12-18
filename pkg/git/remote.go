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
	RefSpecs []string
	Branch   *RemoteBranch
	Branches []*RemoteBranch
}

// NextRemote iterates over next branch of a remote
func (e *RepoEntity) NextRemote() error {
	e.Remote = e.Remotes[(e.currentRemoteIndex()+1)%len(e.Remotes)]
	e.Remote.SyncBranches(e.Branch.Name)
	return e.Publish(RepositoryUpdated, nil)
}

// PreviousRemote iterates over previous branch of a remote
func (e *RepoEntity) PreviousRemote() error {
	e.Remote = e.Remotes[(len(e.Remotes)+e.currentRemoteIndex()-1)%len(e.Remotes)]
	e.Remote.SyncBranches(e.Branch.Name)
	return e.Publish(RepositoryUpdated, nil)
}

// returns the active remote index
func (e *RepoEntity) currentRemoteIndex() int {
	cix := 0
	for i, remote := range e.Remotes {
		if remote.Name == e.Remote.Name {
			cix = i
		}
	}
	return cix
}

// search for remotes in go-git way. It is the short way to get remotes but it
// does not give any insght about remote branches
func (e *RepoEntity) loadRemotes() error {
	r := e.Repository
	e.Remotes = make([]*Remote, 0)

	rms, err := r.Remotes()
	for _, rm := range rms {
		rfs := make([]string, 0)
		for _, rf := range rm.Config().Fetch {
			rfs = append(rfs, string(rf))
		}
		remote := &Remote{
			Name:     rm.Config().Name,
			URL:      rm.Config().URLs,
			RefSpecs: rfs,
		}
		remote.loadRemoteBranches(e)
		if len(remote.Branches) > 0 {
			remote.Branch = remote.Branches[0]
		}
		e.Remotes = append(e.Remotes, remote)

	}
	if err != nil {
		log.Warn("Cannot find remotes " + err.Error())
		return err
	}
	return err
}

// SyncBranches sets the remote branch according to repository's active branch
func (r *Remote) SyncBranches(branchName string) error {
	if err := r.switchRemoteBranch(r.Name + "/" + branchName); err != nil {
		// probably couldn't find, but its ok.
	}
	return nil
}
