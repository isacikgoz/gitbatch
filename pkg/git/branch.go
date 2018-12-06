package git

import (
	"github.com/isacikgoz/gitbatch/pkg/helpers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"regexp"
	"strings"
	"strconv"
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

// returns the active branch of the repository entity by simply getting the
// head reference and searching it from the entities branch slice
func (entity *RepoEntity) getActiveBranch() (branch *Branch) {
	headRef, _ := entity.Repository.Head()
	for _, lb := range entity.Branches {
		if lb.Name == headRef.Name().Short() {
			return lb
		}
	}
	return nil
}

// search for branches in go-git way. It is useful to do so that checkout and
// checkout error handling can be handled by code rather than struggling with
// git cammand and its output
func (entity *RepoEntity) loadLocalBranches() error {
	lbs := make([]*Branch, 0)
	branches, err := entity.Repository.Branches()
	if err != nil {
		log.Warn("Cannot load branches " + err.Error())
		return err
	}
	defer branches.Close()
	branches.ForEach(func(b *plumbing.Reference) error {
		if b.Type() == plumbing.HashReference {
			var push, pull string
			pushables, err := RevList(entity, RevListOptions{
				Ref1: "@{u}",
				Ref2: "HEAD",
				})
			if err != nil {
				push = pushables[0]
			} else {
				push = strconv.Itoa(len(pushables))
			}
			pullables, err := RevList(entity, RevListOptions{
				Ref1: "HEAD",
				Ref2: "@{u}",
				})
			if err != nil {
				pull = pullables[0]
			} else {
				pull = strconv.Itoa(len(pullables))
			}
			clean := entity.isClean()
			branch := &Branch{Name: b.Name().Short(), Reference: b, Pushables: push, Pullables: pull, Clean: clean}
			lbs = append(lbs, branch)
		}
		return nil
	})
	entity.Branches = lbs
	return err
}

// NextBranch checkouts the next branch
func (entity *RepoEntity) NextBranch() *Branch {
	currentBranchIndex := entity.findCurrentBranchIndex()
	if currentBranchIndex == len(entity.Branches)-1 {
		return entity.Branches[0]
	}
	return entity.Branches[currentBranchIndex+1]
}

// PreviousBranch checkouts the previous branch
func (entity *RepoEntity) PreviousBranch() *Branch {
	currentBranchIndex := entity.findCurrentBranchIndex()
	if currentBranchIndex == 0 {
		return entity.Branches[len(entity.Branches)-1]
	}
	return entity.Branches[currentBranchIndex-1]
}

// returns the active branch index
func (entity *RepoEntity) findCurrentBranchIndex() int {
	currentBranch := entity.Branch
	currentBranchIndex := 0
	for i, lbs := range entity.Branches {
		if lbs.Name == currentBranch.Name {
			currentBranchIndex = i
		}
	}
	return currentBranchIndex
}

// Checkout to given branch. If any errors occur, the method returns it instead
// of returning nil
func (entity *RepoEntity) Checkout(branch *Branch) error {
	if branch.Name == entity.Branch.Name {
		return nil
	}
	w, err := entity.Repository.Worktree()
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

	// after checking out we need to refresh some values such as;
	entity.loadCommits()
	entity.Commit = entity.Commits[0]
	entity.Branch = branch
	entity.RefreshPushPull()
	// TODO: same code on 3 different occasion, maybe something wrong?
	// make this conditional on global scale
	if err = entity.Remote.switchRemoteBranch(entity.Remote.Name + "/" + entity.Branch.Name); err != nil {
		// probably couldn't find, but its ok.
		log.Trace("Cannot find proper remote branch " + err.Error())
		return nil
	}
	return nil
}

// checking the branch if it has any changes from its head revision. Initially
// I implemented this with go-git but it was incredibly slow and there is also
// an issue about it: https://github.com/src-d/go-git/issues/844
func (entity *RepoEntity) isClean() bool {
	status := entity.StatusWithGit()
	status = helpers.TrimTrailingNewline(status)
	if status != "?" {
		verbose := strings.Split(status, "\n")
		lastLine := verbose[len(verbose)-1]
		// earlier versions of git returns "working directory clean" instead of
		//"working tree clean" message
		if strings.Contains(lastLine, "working tree clean") ||
			strings.Contains(lastLine, "working directory clean") {
			return true
		}
	}
	return false
}

// RefreshPushPull refreshes the active branchs pushable and pullable count
func (entity *RepoEntity) RefreshPushPull() {
	pushables, err := RevList(entity, RevListOptions{
		Ref1: "@{u}",
		Ref2: "HEAD",
		})
	if err != nil {
		entity.Branch.Pushables = pushables[0]
	} else {
		entity.Branch.Pushables = strconv.Itoa(len(pushables))
	}
	pullables, err := RevList(entity, RevListOptions{
		Ref1: "HEAD",
		Ref2: "@{u}",
		})
	if err != nil {
		entity.Branch.Pullables = pullables[0]
	} else {
		entity.Branch.Pullables = strconv.Itoa(len(pullables))
	}
}

// this function creates the commit entities according to active branchs diffs
// to *its* configured upstream
func (entity *RepoEntity) pullDiffsToUpstream() ([]*Commit, error) {
	remoteCommits := make([]*Commit, 0)
	pullables, err := RevList(entity, RevListOptions{
		Ref1: "HEAD",
		Ref2: "@{u}",
		})
	if err != nil {
		// possibly found nothing or no upstream set
	} else {
		re := regexp.MustCompile(`\r?\n`)
		for _, s := range pullables {
			commit := &Commit{
				Hash:       s,
				Author:     GitShowEmail(entity.AbsPath, s),
				Message:    re.ReplaceAllString(GitShowBody(entity.AbsPath, s), " "),
				Time:       GitShowDate(entity.AbsPath, s),
				CommitType: RemoteCommit,
			}
			remoteCommits = append(remoteCommits, commit)
		}
	}
	return remoteCommits, nil
}

func (entity *RepoEntity) pushDiffsToUpstream() ([]string, error) {
	pushables, err := RevList(entity, RevListOptions{
		Ref1: "@{u}",
		Ref2: "HEAD",
		})
	if err != nil {
		return make([]string, 0), nil
	}
	return pushables, nil
}
