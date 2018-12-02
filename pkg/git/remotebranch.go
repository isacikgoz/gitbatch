package git

import (
	"errors"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

// RemoteBranch is the wrapper of go-git's Reference struct. In addition to
// that, it also holds name of the remote branch
type RemoteBranch struct {
	Name      string
	Reference *plumbing.Reference
}

// NextRemoteBranch iterates to the next remote branch
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

// search for the remote branches of the remote. It takes the go-git's repo
// pointer in order to get storer struct
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
				Name:      b.Name().Short(),
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

// create an iterator for the references. it checks if the reference is a hash
// reference
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

// switches to the given remote branch
func (remote *Remote) switchRemoteBranch(remoteBranchName string) error {
	for _, rb := range remote.Branches {
		if rb.Name == remoteBranchName {
			remote.Branch = rb
			return nil
		}
	}
	return errors.New("Remote branch not found.")
}
