package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"strings"
)

type RemoteBranch struct {
	Name      string
	Reference *plumbing.Reference
}

func (remote *Remote) NextRemoteBranch() error {
	currentRemoteIndex := 0
	for i, rb := range remote.Branches {
		if rb.Reference.Hash() == remote.Branch.Reference.Hash() {
			currentRemoteIndex = i
		}
	}
	// WARNING: DIDN'T CHECK THE LIFE CYCLE
	if currentRemoteIndex == len(remote.Branches)-1 {
		remote.Branch = remote.Branches[0]
	} else {
		remote.Branch = remote.Branches[currentRemoteIndex+1]
	}
	
	return nil
}

func (remote *Remote) loadRemoteBranches(r *git.Repository) error {
	remote.Branches = make([]*RemoteBranch, 0)
	bs, err := remoteBranchesIter(r.Storer)
	if err != nil {
		return err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		if strings.Split(b.Name().Short(), "/")[0] == remote.Name {
			remote.Branches = append(remote.Branches, &RemoteBranch{
				Name: b.Name().Short(),
				Reference: b,
				})
		}
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