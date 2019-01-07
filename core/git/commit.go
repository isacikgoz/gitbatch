package git

import (
	"regexp"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Commit is the lightweight version of go-git's Reference struct. it holds
// hash of the commit, author's e-mail address, Message (subject and body
// combined) commit date and commit type wheter it is local commit or a remote
type Commit struct {
	Hash       string
	Author     string
	Message    string
	Time       string
	CommitType CommitType
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

// NextCommit iterates over next commit of a branch
// TODO: the commits entites can tied to branch instead ot the repository
func (b *Branch) NextCommit() {
	b.State.Commit = b.Commits[(b.currentCommitIndex()+1)%len(b.Commits)]
}

// PreviousCommit iterates to opposite direction
func (b *Branch) PreviousCommit() {
	b.State.Commit = b.Commits[(len(b.Commits)+b.currentCommitIndex()-1)%len(b.Commits)]
}

// returns the active commit index
func (b *Branch) currentCommitIndex() int {
	cix := 0
	for i, c := range b.Commits {
		if c.Hash == b.State.Commit.Hash {
			cix = i
		}
	}
	return cix
}

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
		log.Trace("git log failed " + err.Error())
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
		re := regexp.MustCompile(`\r?\n`)
		cmType := EvenCommit
		for _, lc := range lcs {
			if lc.Hash == c.Hash.String() {
				cmType = LocalCommit
			}
		}
		commit := &Commit{
			Hash:       re.ReplaceAllString(c.Hash.String(), " "),
			Author:     c.Author.Email,
			Message:    re.ReplaceAllString(c.Message, " "),
			Time:       c.Author.When.String(),
			CommitType: cmType,
		}
		b.Commits = append(b.Commits, commit)
		return nil
	})
	if b.State.Commit == nil {
		b.State.Commit = b.Commits[0]
	}
	return err
}

// this function creates the commit entities according to active branchs diffs
// to *its* configured upstream
func (b *Branch) pullDiffsToUpstream(r *Repository) ([]*Commit, error) {
	remoteCommits := make([]*Commit, 0)
	upstream := r.State.Remote.Branch.Reference.Hash().String()

	head := b.Reference.Hash().String()

	pullables, err := RevList(r, RevListOptions{
		Ref1: head,
		Ref2: upstream,
	})

	re := regexp.MustCompile(`\r?\n`)
	if err != nil {
		// possibly found nothing or no upstream set
	} else {
		for _, c := range pullables {

			commit := &Commit{
				Hash:       c.Hash.String(),
				Author:     c.Author.Email,
				Message:    re.ReplaceAllString(c.Message, " "),
				Time:       c.Author.When.String(),
				CommitType: RemoteCommit,
			}
			remoteCommits = append(remoteCommits, commit)
		}
	}
	return remoteCommits, nil
}

// this function returns the hashes of the commits that are not pushed to the
// upstream of the specific branch
func (b *Branch) pushDiffsToUpstream(r *Repository) ([]*Commit, error) {
	notPushedCommits := make([]*Commit, 0)
	upstream := r.State.Remote.Branch.Reference.Hash().String()

	head := b.Reference.Hash().String()

	pushables, err := RevList(r, RevListOptions{
		Ref1: upstream,
		Ref2: head,
	})

	re := regexp.MustCompile(`\r?\n`)
	if err != nil {
		// possibly found nothing or no upstream set
	} else {
		for _, c := range pushables {

			commit := &Commit{
				Hash:       c.Hash.String(),
				Author:     c.Author.Email,
				Message:    re.ReplaceAllString(c.Message, " "),
				Time:       c.Author.When.String(),
				CommitType: LocalCommit,
			}
			notPushedCommits = append(notPushedCommits, commit)
		}
	}
	return notPushedCommits, nil
}
