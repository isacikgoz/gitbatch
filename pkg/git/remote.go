package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

type Remote struct {
	Name      string
	Reference *plumbing.Reference
}

func (entity *RepoEntity) GetRemotes() (remotes []*Remote, err error) {

	r := entity.Repository
    if list, err := remoteBranches(&r); err != nil {
        return remotes, err
    } else {
        for _, r := range list {
            remotes = append(remotes, r)
        }
    }

    return remotes, nil
}

func (entity *RepoEntity) NextRemote() error {

	remotes, err := remoteBranches(&entity.Repository)
	if err != nil {
		return err
	}

	currentRemoteIndex := 0
	for i, remote := range remotes {
		if remote.Reference.Hash() == entity.Remote.Reference.Hash() {
			currentRemoteIndex = i
		}
	}
	// WARNING: DIDN'T CHECK THE LIFE CYCLE
	if currentRemoteIndex == len(remotes)-1 {
		entity.Remote = remotes[0]
	} else {
		entity.Remote = remotes[currentRemoteIndex+1]
	}
	
	return nil
}

func remoteBranches(r *git.Repository) (remotes []*Remote, err error) {
	bs, err := remoteBranchesIter(r.Storer)
	if err != nil {
		return remotes, err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		remotes = append(remotes, &Remote{
			Name: b.Name().Short(),
			Reference: b,
			})
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