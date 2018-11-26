package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Branch struct {
	Name       string
	Reference  *plumbing.Reference
	Pushables  string
	Pullables  string
	Clean      bool
}

func (entity *RepoEntity) GetActiveBranch() (branch *Branch) {
	headRef, _ := entity.Repository.Head()
	localBranches, err := entity.LocalBranches()
	if err != nil {
		return nil
	}
	for _, lb := range localBranches {
		if lb.Name == headRef.Name().Short() {
			return lb
		}
	}
	return nil
}

func (entity *RepoEntity) LocalBranches() (lbs []*Branch, err error) {
	branches, err := entity.Repository.Branches()
	if err != nil {
		return nil, err
	}
	defer branches.Close()
	branches.ForEach(func(b *plumbing.Reference) error {
		if b.Type() == plumbing.HashReference {
			push, pull := UpstreamDifferenceCount(entity.AbsPath)
			clean := entity.isClean()
			branch := &Branch{Name: b.Name().Short(), Reference: b, Pushables: push, Pullables: pull, Clean: clean}
        	lbs = append(lbs, branch)
    	}
    	return nil
	})
	return lbs, err
}

func (entity *RepoEntity) NextBranch() *Branch {

	currentBranch := entity.GetActiveBranch()
	localBranches, err := entity.LocalBranches()
	if err != nil {
		return currentBranch
	}

	currentBranchIndex := 0
	for i, lbs := range localBranches {
		if lbs.Name == currentBranch.Name {
			currentBranchIndex = i
		}
	}

	if currentBranchIndex == len(localBranches)-1 {
		return localBranches[0]
	}
	return localBranches[currentBranchIndex+1]
}

func (entity *RepoEntity) Checkout(branch *Branch) error {
	if branch.Name == entity.Branch.Name {
		return nil
	}
	w, err := entity.Repository.Worktree()
	if err != nil {
		return err
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Branch: branch.Reference.Name(),
	}); err != nil {
		return err
	}
	entity.Branch = branch
	return nil
}

func (entity *RepoEntity) isClean() bool {
	worktree, err := entity.Repository.Worktree()
	if err != nil {
		return true
	}
	status, err := worktree.Status()
	if err != nil {
		return false
	}
	return status.IsClean()
}