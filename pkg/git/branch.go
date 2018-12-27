package git

import (
	"github.com/isacikgoz/gitbatch/pkg/helpers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"strconv"
	"strings"
)

// Branch is the wrapper of go-git's Reference struct. In addition to that, it
// also holds name of the branch, pullable and pushable commit count from the
// branchs' upstream. It also tracks if the repository has unstaged or uncommit-
// ed changes
type Branch struct {
	Name      string
	Reference *plumbing.Reference
	Pushables string
	Pullables string
	Clean     bool
}

// search for branches in go-git way. It is useful to do so that checkout and
// checkout error handling can be handled by code rather than struggling with
// git cammand and its output
func (e *RepoEntity) loadLocalBranches() error {
	lbs := make([]*Branch, 0)
	bs, err := e.Repository.Branches()
	if err != nil {
		log.Warn("Cannot load branches " + err.Error())
		return err
	}
	defer bs.Close()
	headRef, err := e.Repository.Head()
	if err != nil {
		return err
	}
	var branchFound bool
	bs.ForEach(func(b *plumbing.Reference) error {
		if b.Type() == plumbing.HashReference {
			var push, pull string
			pushables, err := RevList(e, RevListOptions{
				Ref1: "@{u}",
				Ref2: "HEAD",
			})
			if err != nil {
				push = pushables[0]
			} else {
				push = strconv.Itoa(len(pushables))
			}
			pullables, err := RevList(e, RevListOptions{
				Ref1: "HEAD",
				Ref2: "@{u}",
			})
			if err != nil {
				pull = pullables[0]
			} else {
				pull = strconv.Itoa(len(pullables))
			}
			clean := e.isClean()
			branch := &Branch{
				Name:      b.Name().Short(),
				Reference: b,
				Pushables: push,
				Pullables: pull,
				Clean:     clean,
			}
			if b.Name() == headRef.Name() {
				e.Branch = branch
				branchFound = true
			}
			lbs = append(lbs, branch)
		}
		return nil
	})
	if !branchFound {
		branch := &Branch{
			Name:      headRef.Hash().String(),
			Reference: headRef,
			Pushables: "?",
			Pullables: "?",
			Clean:     e.isClean(),
		}
		lbs = append(lbs, branch)
		e.Branch = branch
	}
	e.Branches = lbs
	return err
}

// NextBranch checkouts the next branch
func (e *RepoEntity) NextBranch() *Branch {
	return e.Branches[(e.currentBranchIndex()+1)%len(e.Branches)]
}

// PreviousBranch checkouts the previous branch
func (e *RepoEntity) PreviousBranch() *Branch {
	return e.Branches[(len(e.Branches)+e.currentBranchIndex()-1)%len(e.Branches)]
}

// returns the active branch index
func (e *RepoEntity) currentBranchIndex() int {
	bix := 0
	for i, lbs := range e.Branches {
		if lbs.Name == e.Branch.Name {
			bix = i
		}
	}
	return bix
}

// Checkout to given branch. If any errors occur, the method returns it instead
// of returning nil
func (e *RepoEntity) Checkout(branch *Branch) error {
	if branch.Name == e.Branch.Name {
		return nil
	}
	w, err := e.Repository.Worktree()
	if err != nil {
		log.Warn("Cannot get work tree " + err.Error())
		return err
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Branch: branch.Reference.Name(),
	}); err != nil {
		log.Warn("Cannot checkout " + err.Error())
		return err
	}

	// make this conditional on global scale
	err = e.Remote.SyncBranches(branch.Name)
	return e.Refresh()
}

// checking the branch if it has any changes from its head revision. Initially
// I implemented this with go-git but it was incredibly slow and there is also
// an issue about it: https://github.com/src-d/go-git/issues/844
func (e *RepoEntity) isClean() bool {
	s := e.StatusWithGit()
	s = helpers.TrimTrailingNewline(s)
	if s != "?" {
		vs := strings.Split(s, "\n")
		line := vs[len(vs)-1]
		// earlier versions of git returns "working directory clean" instead of
		//"working tree clean" message
		if strings.Contains(line, "working tree clean") ||
			strings.Contains(line, "working directory clean") {
			return true
		}
	}
	return false
}
