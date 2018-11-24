package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
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

	remotes, err := remoteBranches(&entity.Repository)
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

func remoteBranches(r *git.Repository) (remotes []string, err error) {
	bs, err := remoteBranchesIter(r.Storer)
	if err != nil {
		return remotes, err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		remotes = append(remotes, b.Name().Short())
		return nil
	})
	if err != nil {
		return remotes, err
	}
	return remotes, err
}

func remoteBranchesIter(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		if ref.Type() == plumbing.HashReference {
			return ref.Name().IsRemote()
		}
		return false
	}, refs), nil
}
