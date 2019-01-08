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

// search for the remote branches of the remote. It takes the go-git's repo
// pointer in order to get storer struct
func (rm *Remote) loadRemoteBranches(r *Repository) error {
	rm.Branches = make([]*RemoteBranch, 0)
	bs, err := remoteBranchesIter(r.Repo.Storer)
	if err != nil {
		log.Warn("Cannot initiate iterator " + err.Error())
		return err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		if strings.Split(b.Name().Short(), "/")[0] == rm.Name {
			rm.Branches = append(rm.Branches, &RemoteBranch{
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
func (rm *Remote) switchRemoteBranch(remoteBranchName string) error {
	for _, rb := range rm.Branches {
		if rb.Name == remoteBranchName {
			rm.Branch = rb
			return nil
		}
	}
	return errors.New("Remote branch not found")
}
