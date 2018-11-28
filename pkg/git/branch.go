package git

import (
	"github.com/isacikgoz/gitbatch/pkg/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"strings"
	"regexp"
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
	// this method is painfully slow
	// worktree, err := entity.Repository.Worktree()
	// if err != nil {
	// 	return true
	// }
	// status, err := worktree.Status()
	// if err != nil {
	// 	return false
	// }
	// return status.IsClean()
	status := entity.StatusWithGit()
	status = utils.TrimTrailingNewline(status)
	if status != "?" {
		verbose := strings.Split(status, "\n")
		lastLine := verbose[len(verbose)-1]
		if strings.Contains(lastLine, "working tree clean") {
			return true
		}
	}
	return false
}

func (entity *RepoEntity) RefreshPushPull() {
	entity.Branch.Pushables, entity.Branch.Pullables = UpstreamDifferenceCount(entity.AbsPath)
}

func (entity *RepoEntity) pullDiffsToUpstream() ([]*Commit, error) {
	remoteCommits := make([]*Commit, 0)
	hashes := UpstreamPullDiffs(entity.AbsPath)
	re := regexp.MustCompile(`\r?\n`)
	if hashes != "?" {
		sliced := strings.Split(hashes, "\n")
		for _, s := range sliced {
			if len(s) == 40 {
				commit := &Commit{
				Hash: s,
				Author: GitShowEmail(entity.AbsPath, s),
				Message: re.ReplaceAllString(GitShowBody(entity.AbsPath, s), " "),
				Time: GitShowDate(entity.AbsPath, s),
				CommitType: RemoteCommit,
			}
			remoteCommits = append(remoteCommits, commit)
			}
		}
	}
	return remoteCommits, nil
}
