package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"strings"
)

type Branch struct {
	Name      string
	Reference *plumbing.Reference
	Pushables string
	Pullables string
	Clean     bool
}

func (entity *RepoEntity) getActiveBranch() (branch *Branch) {
	headRef, _ := entity.Repository.Head()
	for _, lb := range entity.Branches {
		if lb.Name == headRef.Name().Short() {
			return lb
		}
	}
	return nil
}

func (entity *RepoEntity) loadLocalBranches() error {
	lbs := make([]*Branch, 0)
	branches, err := entity.Repository.Branches()
	if err != nil {
		return err
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
	entity.Branches = lbs
	return err
}

func (entity *RepoEntity) NextBranch() *Branch {
	currentBranch := entity.Branch
	currentBranchIndex := 0
	for i, lbs := range entity.Branches {
		if lbs.Name == currentBranch.Name {
			currentBranchIndex = i
		}
	}
	if currentBranchIndex == len(entity.Branches)-1 {
		return entity.Branches[0]
	}
	return entity.Branches[currentBranchIndex+1]
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
	entity.loadCommits()
	entity.Commit = entity.Commits[0]
	entity.Branch = branch
	entity.Branch.Pushables, entity.Branch.Pullables = UpstreamDifferenceCount(entity.AbsPath)
	// TODO: same code on 3 different occasion, maybe something wrong?
	// make this conditional on global scale
	if err = entity.Remote.switchRemoteBranch(entity.Remote.Name + "/" + entity.Branch.Name); err !=nil {
		// probably couldn't find, but its ok.
		return nil
	}
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

func (entity *RepoEntity) RefreshPushPull() {
	entity.Branch.Pushables, entity.Branch.Pullables = UpstreamDifferenceCount(entity.AbsPath)
}

func (entity *RepoEntity) PushDiffsToUpstream() error {
	hashes := UpstreamPushDiffs(entity.AbsPath)
	if hashes != "?" {
		sliced := strings.Split(hashes, "\n")
		for _, s := range sliced {
			GitShow(entity.AbsPath, s)
		}
	}
	return nil
}

func (entity *RepoEntity) PullDiffsToUpstream() error {
	hashes := UpstreamPullDiffs(entity.AbsPath)
	if hashes != "?" {
		sliced := strings.Split(hashes, "\n")
		for _, s := range sliced {
			GitShow(entity.AbsPath, s)
		}
	}
	return nil
}