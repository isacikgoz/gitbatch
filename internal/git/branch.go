package git

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

// BranchState hold the ref commit
type BranchState struct {
	Commit *Commit
}

const (
	revlistCommand = "rev-list"
	hashLength     = 40
)

// search for branches in go-git way. It is useful to do so that checkout and
// checkout error handling can be handled by code rather than struggling with
// git command and its output
func (r *Repository) initBranches() error {
	lbs := make([]*Branch, 0)
	bs, err := r.Repo.Branches()
	if err != nil {
		return err
	}
	defer bs.Close()
	headRef, err := r.Repo.Head()
	if err != nil {
		return err
	}
	var branchFound bool
	var push, pull string
	_ = bs.ForEach(func(b *plumbing.Reference) error {
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
	if err == nil {
		r.State.Branch.Upstream = rb
	}

	r.Branches = lbs
	return nil
}

// Checkout to given branch. If any errors occur, the method returns it instead
// of returning nil
func (r *Repository) Checkout(b *Branch) error {
	if b.Name == r.State.Branch.Name {
		return nil
	}

	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Branch: b.Reference.Name(),
	}); err != nil {
		return err
	}
	r.State.Branch = b

	rb, err := getUpstream(r, r.State.Branch.Name)
	if err == nil {
		r.State.Branch.Upstream = rb
	}
	_ = b.initCommits(r)

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

// RevList is the legacy implementation of "git rev-list" command.
func RevList(r *Repository, options RevListOptions) ([]*object.Commit, error) {
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
	commits := make([]*object.Commit, 0)
	for _, hash := range hashes {
		if len(hash) == hashLength {
			c, err := r.Repo.CommitObject(plumbing.NewHash(hash))
			if err != nil {
				continue
			}
			commits = append(commits, c)
		}
	}
	sort.Sort(CommitTime(commits))
	return commits, nil
}

// SyncRemoteAndBranch synchronizes remote branch with current branch
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

	head := headRef.Hash().String()
	var push, pull string
	pushables, err := RevList(r, RevListOptions{
		Ref1: b.Upstream.Reference.Hash().String(),
		Ref2: head,
	})
	if err != nil {
		push = "?"
	} else {
		push = strconv.Itoa(len(pushables))
	}
	pullables, err := RevList(r, RevListOptions{
		Ref1: head,
		Ref2: b.Upstream.Reference.Hash().String(),
	})
	if err != nil {
		pull = "?"
	} else {
		pull = strconv.Itoa(len(pullables))
	}
	b.Pullables = pull
	b.Pushables = push
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
		return nil, fmt.Errorf("upstream not found")
	}

	args = []string{"config", "--get", "branch." + branchName + ".merge"}
	cmd = exec.Command("git", args...)
	cmd.Dir = r.AbsPath
	cm, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(cm), branchName) {
		return nil, fmt.Errorf("default merge branch found")
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
	return nil, fmt.Errorf("upstream not found")
}
