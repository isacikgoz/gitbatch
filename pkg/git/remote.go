package git

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

type Remote struct {
	Name      string
	Reference *plumbing.Reference
}

func (entity *RepoEntity) NextRemote() error {
	currentRemoteIndex := 0
	for i, remote := range entity.Remotes {
		if remote.Reference.Hash() == entity.Remote.Reference.Hash() {
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

func (entity *RepoEntity) loadRemoteBranches() error {
	r := entity.Repository
	entity.Remotes = make([]*Remote, 0)
	bs, err := remoteBranchesIter(r.Storer)
	if err != nil {
		return err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		entity.Remotes = append(entity.Remotes, &Remote{
			Name: b.Name().Short(),
			Reference: b,
			})
		return nil
	})
	if err != nil {
		return err
	}
	return err
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