package git

import (
	"regexp"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit is the lightweight version of go-git's Reference struct. it holds
// hash of the commit, author's e-mail address, Message (subject and body
// combined) commit date and commit type whether it is local commit or a remote
type Commit struct {
	Hash       string
	Author     *Contributor
	Commiter   *Contributor
	Message    string
	Time       string
	CommitType CommitType
	C          *object.Commit
}

// Contributor is the person
type Contributor struct {
	Name  string
	Email string
	When  time.Time
}

// CommitType is the Type of the commit; it can be local or remote (upstream diff)
type CommitType string

const (
	// LocalCommit is the commit that not pushed to remote branch
	LocalCommit CommitType = "local"
	// EvenCommit is the commit that recorded locally
	EvenCommit CommitType = "even"
	// RemoteCommit is the commit that not merged to local branch
	RemoteCommit CommitType = "remote"
)

// loads the local commits by simply using git log way. ALso, gets the upstream
// diff commits
func (b *Branch) initCommits(r *Repository) error {
	b.Commits = make([]*Commit, 0)
	ref := b.Reference

	// git log first
	cIter, err := r.Repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return err
	}
	defer cIter.Close()
	// find commits that fetched from upstream but not merged commits
	rmcs, _ := b.pullDiffsToUpstream(r)
	b.Commits = append(b.Commits, rmcs...)

	// find commits that not pushed to upstream
	lcs, _ := b.pushDiffsToUpstream(r)

	// ... just iterates over the commits
	err = cIter.ForEach(func(c *object.Commit) error {

		cmType := EvenCommit
		for _, lc := range lcs {
			if lc.Hash == c.Hash.String() {
				cmType = LocalCommit
			}
		}

		commit := commit(c, cmType)
		b.Commits = append(b.Commits, commit)
		return nil
	})
	if err != nil {
		return err
	}
	if b.State.Commit == nil {
		b.State.Commit = b.Commits[0]
	}
	return nil
}

// this function creates the commit entities according to active branchs diffs
// to *its* configured upstream
func (b *Branch) pullDiffsToUpstream(r *Repository) ([]*Commit, error) {
	remoteCommits := make([]*Commit, 0)
	if r.State.Branch.Upstream == nil {
		return remoteCommits, nil
	}
	head := b.Reference.Hash().String()

	pullables, err := RevList(r, RevListOptions{
		Ref1: head,
		Ref2: b.Upstream.Reference.Hash().String(),
	})
	if err != nil {
		// possibly found nothing or no upstream set
	} else {
		for _, c := range pullables {
			commit := commit(c, RemoteCommit)
			remoteCommits = append(remoteCommits, commit)
		}
	}
	return remoteCommits, nil
}

// this function returns the hashes of the commits that are not pushed to the
// upstream of the specific branch
func (b *Branch) pushDiffsToUpstream(r *Repository) ([]*Commit, error) {
	notPushedCommits := make([]*Commit, 0)
	if r.State.Branch.Upstream == nil {
		return notPushedCommits, nil
	}
	head := b.Reference.Hash().String()

	pushables, err := RevList(r, RevListOptions{
		Ref1: b.Upstream.Reference.Hash().String(),
		Ref2: head,
	})

	if err != nil {
		// possibly found nothing or no upstream set
	} else {
		for _, c := range pushables {
			commit := commit(c, LocalCommit)
			notPushedCommits = append(notPushedCommits, commit)
		}
	}
	return notPushedCommits, nil
}

func commit(c *object.Commit, t CommitType) *Commit {
	commit := &Commit{
		Hash: c.Hash.String(),
		Author: &Contributor{
			Name:  c.Author.Name,
			Email: c.Author.Email,
			When:  c.Author.When,
		},
		Commiter: &Contributor{
			Name:  c.Committer.Name,
			Email: c.Committer.Email,
			When:  c.Committer.When,
		},
		Message:    c.Message,
		CommitType: t,
		C:          c,
	}
	return commit
}

// DiffStat Show diff stat
func (c *Commit) DiffStat(done chan bool) string {

	var str string
	defer recoverDiff(&str)
	if c.C == nil {
		return ""
	}
	d, err := c.C.Stats()
	if err != nil {
		return ""
	}
	str = d.String()
	done <- true
	return str
}

func (c *Commit) String() string {
	d := "Hash:" + " " + c.Hash
	d = d + "\n" + "Author:" + " " + c.Author.Name + " <" + c.Author.Email + ">"
	d = d + "\n" + "Date:" + " " + c.Author.When.String() + "\n"
	re := regexp.MustCompile(`\r?\n`)
	s := re.Split(c.Message, -1)
	for _, l := range s {
		d = d + "\n" + " " + l
	}
	return d
}

func recoverDiff(str *string) {
	if r := recover(); r != nil {
		*str = "diffstat overloaded"
	}
}
