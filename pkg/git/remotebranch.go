package git

import (
	"errors"
	"regexp"
	"strings"

	"github.com/isacikgoz/gitbatch/pkg/helpers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

// RemoteBranch is the wrapper of go-git's Reference struct. In addition to
// that, it also holds name of the remote branch
type RemoteBranch struct {
	Name      string
	Reference *plumbing.Reference
	Deleted   bool
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
func (remote *Remote) loadRemoteBranches(entity *RepoEntity) error {
	remote.Branches = make([]*RemoteBranch, 0)
	bs, err := remoteBranchesIter(entity.Repository.Storer)
	if err != nil {
		log.Warn("Cannot initiate iterator " + err.Error())
		return err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		deleted := false
		if strings.Split(b.Name().Short(), "/")[0] == remote.Name {
			remote.Branches = append(remote.Branches, &RemoteBranch{
				Name:      b.Name().Short(),
				Reference: b,
				Deleted:   deleted,
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
func (remote *Remote) switchRemoteBranch(remoteBranchName string) error {
	for _, rb := range remote.Branches {
		if rb.Name == remoteBranchName {
			remote.Branch = rb
			return nil
		}
	}
	return errors.New("Remote branch not found.")
}

func deletedRemoteBranches(entity *RepoEntity, remote string) ([]string, error) {
	deletedRemoteBranches := make([]string, 0)
	output := entity.DryFetchAndPruneWithGit(remote)
	output = helpers.TrimTrailingNewline(output)
	re := regexp.MustCompile(` - \[deleted\].+-> `)
	if output != "?" {
		sliced := strings.Split(output, "\n")
		for _, s := range sliced {
			if re.MatchString(s) {
				ss := re.ReplaceAllString(s, "")
				deletedRemoteBranches = append(deletedRemoteBranches, ss)
			}
		}
	}
	return deletedRemoteBranches, nil
}
