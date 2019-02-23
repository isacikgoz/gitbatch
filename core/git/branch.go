package git

import (
	"errors"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
)

// Branch is the wrapper of go-git's Reference struct. In addition to that, it
// also holds name of the branch, pullable and pushable commit count from the
// branchs' upstream. It also tracks if the repository has unstaged or uncommit-
// ed changes
type Branch struct {
	Name      string
	Reference *plumbing.Reference
	Upstream  *RemoteBranch
	Commits   []*Commit
	State     *BranchState
	Pushables string
	Pullables string
	Clean     bool
}

type BranchState struct {
	Commit *Commit
}

var (
	revlistCommand = "rev-list"
	hashLength     = 40
)

// search for branches in go-git way. It is useful to do so that checkout and
// checkout error handling can be handled by code rather than struggling with
// git cammand and its output
func (r *Repository) initBranches() error {
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
	var push, pull string
	bs.ForEach(func(b *plumbing.Reference) error {
		if b.Type() != plumbing.HashReference {
			return nil
		}
		clean := r.isClean()
		branch := &Branch{
			Name:      b.Name().Short(),
			Reference: b,
			State:     &BranchState{},
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
			State:     &BranchState{},
			Pushables: "?",
			Pullables: "?",
			Clean:     r.isClean(),
		}
		lbs = append(lbs, branch)
		r.State.Branch = branch
	}
	rb, err := getUpstream(r, r.State.Branch.Name)
	if err != nil {
		log.Warn("Upstream not set " + r.Name)
	} else {
		r.State.Branch.Upstream = rb
	}

	r.Branches = lbs
	return nil
}

// Checkout to given branch. If any errors occur, the method returns it instead
// of returning nil
func (r *Repository) Checkout(b *Branch) error {
	// var reinit bool
	if b.Name == r.State.Branch.Name {
		return nil
	}
	// if it already loaded its commits, consider reload again
	// if len(b.Commits) > 0 {
	// 	reinit = true
	// }
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
	r.State.Branch = b

	rb, err := getUpstream(r, r.State.Branch.Name)
	if err != nil {
		log.Warn("Upstream not set")
	} else {
		r.State.Branch.Upstream = rb
	}
	// if reinit {
	b.initCommits(r)
	// }
	// if err := r.Refresh(); err != nil {
	// 	return err
	// }
	if err := r.Publish(BranchUpdated, nil); err != nil {
		return err
	}
	if err := r.SyncRemoteAndBranch(b); err != nil {
		return err
	}
	return r.Publish(RepositoryUpdated, nil)
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

// RevList is native implemetation of git rev-list command
func RevList(r *Repository, opts RevListOptions) ([]*object.Commit, error) {

	commits := make([]*object.Commit, 0)
	ref1hist, err := revlist.Objects(r.Repo.Storer, []plumbing.Hash{plumbing.NewHash(opts.Ref1)}, nil)
	if err != nil {
		return nil, err
	}
	ref2hist, err := revlist.Objects(r.Repo.Storer, []plumbing.Hash{plumbing.NewHash(opts.Ref2)}, ref1hist)
	if err != nil {
		return nil, err
	}

	for _, h := range ref2hist {
		c, err := r.Repo.CommitObject(h)
		if err != nil {
			continue
		}
		commits = append(commits, c)
	}
	sort.Sort(CommitTime(commits))
	return commits, err
}

// SyncRemoteAndBranch is essegin ziki
func (r *Repository) SyncRemoteAndBranch(b *Branch) error {
	headRef, err := r.Repo.Head()
	if err != nil {
		return err
	}
	if b.Upstream == nil {
		b.Pullables = "?"
		b.Pushables = "?"
		return nil
	}

	headd := headRef.Hash().String()
	var push, pull string
	pushables, err := RevList(r, RevListOptions{
		Ref1: b.Upstream.Reference.Hash().String(),
		Ref2: headd,
	})
	if err != nil {
		push = "?"
	} else {
		push = strconv.Itoa(len(pushables))
	}
	pullables, err := RevList(r, RevListOptions{
		Ref1: headd,
		Ref2: b.Upstream.Reference.Hash().String(),
	})
	if err != nil {
		pull = "?"
	} else {
		pull = strconv.Itoa(len(pullables))
	}
	b.Pullables = pull
	b.Pushables = push
	// return b.initCommits(r)
	return nil
}

// InitializeCommits loads the commits
func (b *Branch) InitializeCommits(r *Repository) error {
	return b.initCommits(r)
}

func getUpstream(r *Repository, branchName string) (*RemoteBranch, error) {
	args := []string{"config", "--get", "branch." + branchName + ".remote"}
	cmd := exec.Command("git", args...)
	cmd.Dir = r.AbsPath
	cr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New("upstream not found")
	}

	args = []string{"config", "--get", "branch." + branchName + ".merge"}
	cmd = exec.Command("git", args...)
	cmd.Dir = r.AbsPath
	cm, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(cm), branchName) {
		return nil, errors.New("default merge branch found")
	}

	for _, rm := range r.Remotes {
		if rm.Name == strings.TrimSpace(string(cr)) {
			r.State.Remote = rm
		}
	}

	for _, rb := range r.State.Remote.Branches {
		if rb.Name == r.State.Remote.Name+"/"+branchName {
			return rb, nil
		}
	}
	return nil, errors.New("upstream not found")
}

// trimTrailingNewline removes the trailing new line form a string. this method
// is used mostly on outputs of a command
func trimTrailingNewline(s string) string {
	if strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\r") {
		return s[:len(s)-1]
	}
	return s
}
