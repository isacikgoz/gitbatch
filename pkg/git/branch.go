package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (entity *RepoEntity) GetActiveBranch() string{
	headRef, _ := entity.Repository.Head()
	return headRef.Name().Short()
}

func (entity *RepoEntity) LocalBranches() (lbs []string, err error){
	branches, err := entity.Repository.Branches()
	if err != nil {
		return nil, err
	}
	defer branches.Close()
	branches.ForEach(func(b *plumbing.Reference) error {
		if b.Type() == plumbing.HashReference {
        	lbs = append(lbs, b.Name().Short())
    	}
    	return nil
	})
	return lbs, err
}

func (entity *RepoEntity) NextBranch() string{

	currentBranch := entity.GetActiveBranch()
	localBranches, err := entity.LocalBranches()
	if err != nil {
		return currentBranch
	}

	currentBranchIndex := 0
	for i, lbs := range localBranches {
		if lbs == currentBranch {
			currentBranchIndex = i
		}
	}

	if currentBranchIndex == len(localBranches)-1 {
		return localBranches[0]
	}
	return localBranches[currentBranchIndex+1]
}

func (entity *RepoEntity) Checkout(branchName string) error {
	if branchName == entity.Branch {
		return nil
	}
	w, err := entity.Repository.Worktree()
	if err != nil {
		return err
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	}); err != nil {
		return err
	}
	entity.Branch = branchName
	entity.Pushables, entity.Pullables = UpstreamDifferenceCount(entity.AbsPath)
	return nil
}

func (entity *RepoEntity) IsClean() (bool, error) {
	worktree, err := entity.Repository.Worktree()
	if err != nil {
		return true, nil
	}
	status, err := worktree.Status()
	if err != nil {
		return status.IsClean(), nil
	}
	return false, nil
}
