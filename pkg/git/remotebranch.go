package git

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
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
func (r *Remote) NextRemoteBranch(e *RepoEntity) error {
	r.Branch = r.Branches[(r.currentRemoteBranchIndex()+1)%len(r.Branches)]
	return e.Publish(RepositoryUpdated, nil)
}

// PreviousRemoteBranch iterates to the previous remote branch
func (r *Remote) PreviousRemoteBranch(e *RepoEntity) error {
	r.Branch = r.Branches[(len(r.Branches)+r.currentRemoteBranchIndex()-1)%len(r.Branches)]
	return e.Publish(RepositoryUpdated, nil)
}

// returns the active remote branch index
func (r *Remote) currentRemoteBranchIndex() int {
	cix := 0
	for i, rb := range r.Branches {
		if rb.Reference.Hash() == r.Branch.Reference.Hash() {
			cix = i
		}
	}
	return cix
}

// search for the remote branches of the remote. It takes the go-git's repo
// pointer in order to get storer struct
func (r *Remote) loadRemoteBranches(e *RepoEntity) error {
	r.Branches = make([]*RemoteBranch, 0)
	bs, err := remoteBranchesIter(e.Repository.Storer)
	if err != nil {
		log.Warn("Cannot initiate iterator " + err.Error())
		return err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		if strings.Split(b.Name().Short(), "/")[0] == r.Name {
			r.Branches = append(r.Branches, &RemoteBranch{
				Name:      b.Name().Short(),
				Reference: b,
			})
		}
		return nil
	})
	return err
}

// create an iterator for the references. it checks if the reference is a hash
// reference
func remoteBranchesIter(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		log.Warn("Cannot find references " + err.Error())
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
func (r *Remote) switchRemoteBranch(remoteBranchName string) error {
	for _, rb := range r.Branches {
		if rb.Name == remoteBranchName {
			r.Branch = rb
			return nil
		}
	}
	return errors.New("Remote branch not found")
}
