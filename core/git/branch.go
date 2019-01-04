package git

import (
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

var (
	revlistCommand = "rev-list"
	hashLength     = 40
)

// search for branches in go-git way. It is useful to do so that checkout and
// checkout error handling can be handled by code rather than struggling with
// git cammand and its output
func (r *Repository) loadLocalBranches() error {
	lbs := make([]*Branch, 0)
	bs, err := r.Repo.Branches()
	if err != nil {
		log.Warn("Cannot load branches " + err.Error())
		return err
	}
	defer bs.Close()
	headRef, err := r.Repo.Head()
	if err != nil {
		return err
	}
	var branchFound bool
	bs.ForEach(func(b *plumbing.Reference) error {
		if b.Type() != plumbing.HashReference {
			return nil
		}

		var push, pull string
		pushables, err := RevList(r, RevListOptions{
			Ref1: "@{u}",
			Ref2: "HEAD",
		})
		if err != nil {
			push = "?"
		} else {
			push = strconv.Itoa(len(pushables))
		}
		pullables, err := RevList(r, RevListOptions{
			Ref1: "HEAD",
			Ref2: "@{u}",
		})
		if err != nil {
			pull = "?"
		} else {
			pull = strconv.Itoa(len(pullables))
		}
		clean := r.isClean()
		branch := &Branch{
			Name:      b.Name().Short(),
			Reference: b,
			Pushables: push,
			Pullables: pull,
			Clean:     clean,
		}
		if b.Name() == headRef.Name() {
			r.State.Branch = branch
			branchFound = true
		}
		lbs = append(lbs, branch)

		return nil
	})
	if !branchFound {
		branch := &Branch{
			Name:      headRef.Hash().String(),
			Reference: headRef,
			Pushables: "?",
			Pullables: "?",
			Clean:     r.isClean(),
		}
		lbs = append(lbs, branch)
		r.State.Branch = branch
	}
	r.Branches = lbs
	return err
}

// NextBranch checkouts the next branch
func (r *Repository) NextBranch() *Branch {
	return r.Branches[(r.currentBranchIndex()+1)%len(r.Branches)]
}

// PreviousBranch checkouts the previous branch
func (r *Repository) PreviousBranch() *Branch {
	return r.Branches[(len(r.Branches)+r.currentBranchIndex()-1)%len(r.Branches)]
}

// returns the active branch index
func (r *Repository) currentBranchIndex() int {
	bix := 0
	for i, lbs := range r.Branches {
		if lbs.Name == r.State.Branch.Name {
			bix = i
		}
	}
	return bix
}

// Checkout to given branch. If any errors occur, the method returns it instead
// of returning nil
func (r *Repository) Checkout(b *Branch) error {
	if b.Name == r.State.Branch.Name {
		return nil
	}

	w, err := r.Repo.Worktree()
	if err != nil {
		log.Warn("Cannot get work tree " + err.Error())
		return err
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Branch: b.Reference.Name(),
	}); err != nil {
		log.Warn("Cannot checkout " + err.Error())
		return err
	}

	// make this conditional on global scale
	// we don't care if this function returns an error
	r.State.Remote.SyncBranches(b.Name)

	return r.Refresh()
}

// checking the branch if it has any changes from its head revision. Initially
// I implemented this with go-git but it was incredibly slow and there is also
// an issue about it: https://github.com/src-d/go-git/issues/844
func (r *Repository) isClean() bool {
	args := []string{"status"}
	cmd := exec.Command("git", args...)
	cmd.Dir = r.AbsPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	s := string(out)
	if strings.HasSuffix(s, "\n") {
		s = s[:len(s)-1]
	}
	if len(s) >= 0 {
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

// RevListOptions defines the rules of rev-list func
type RevListOptions struct {
	// Ref1 is the first reference hash to link
	Ref1 string
	// Ref2 is the second reference hash to link
	Ref2 string
}

// RevList returns the commit hashes that are links from the given commit(s).
// The output is given in reverse chronological order by default.
func RevList(r *Repository, options RevListOptions) ([]string, error) {
	args := make([]string, 0)
	args = append(args, revlistCommand)
	if len(options.Ref1) > 0 && len(options.Ref2) > 0 {
		arg1 := options.Ref1 + ".." + options.Ref2
		args = append(args, arg1)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = r.AbsPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	s := string(out)

	hashes := strings.Split(s, "\n")
	for _, hash := range hashes {
		if len(hash) != hashLength {
			return make([]string, 0), nil
		}
		break
	}
	return hashes, nil
}
