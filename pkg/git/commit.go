package git

import (
	"regexp"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
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
func (e *RepoEntity) NextCommit() {
	e.Commit = e.Commits[(e.currentCommitIndex()+1)%len(e.Commits)]
}

// PreviousCommit iterates to opposite direction
func (e *RepoEntity) PreviousCommit() {
	e.Commit = e.Commits[(len(e.Commits)+e.currentCommitIndex()-1)%len(e.Commits)]
}

// returns the active commit index
func (e *RepoEntity) currentCommitIndex() int {
	cix := 0
	for i, c := range e.Commits {
		if c.Hash == e.Commit.Hash {
			cix = i
		}
	}
	return cix
}

// loads the local commits by simply using git log way. ALso, gets the upstream
// diff commits
func (e *RepoEntity) loadCommits() error {
	r := e.Repository
	e.Commits = make([]*Commit, 0)
	ref, err := r.Head()
	if err != nil {
		log.Trace("Cannot get HEAD " + err.Error())
		return err
	}
	// git log first
	cIter, err := r.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		log.Trace("git log failed " + err.Error())
		return err
	}
	defer cIter.Close()
	// find commits that fetched from upstream but not merged commits
	rmcs, _ := e.pullDiffsToUpstream()
	e.Commits = append(e.Commits, rmcs...)

	// find commits that not pushed to upstream
	lcs, _ := e.pushDiffsToUpstream()

	// ... just iterates over the commits
	err = cIter.ForEach(func(c *object.Commit) error {
		re := regexp.MustCompile(`\r?\n`)
		cmType := EvenCommit
		for _, lc := range lcs {
			if lc == re.ReplaceAllString(c.Hash.String(), " ") {
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
		e.Commits = append(e.Commits, commit)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// this function creates the commit entities according to active branchs diffs
// to *its* configured upstream
func (e *RepoEntity) pullDiffsToUpstream() ([]*Commit, error) {
	remoteCommits := make([]*Commit, 0)
	pullables, err := RevList(e, RevListOptions{
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
				Author:     GitShowEmail(e.AbsPath, s),
				Message:    re.ReplaceAllString(GitShowBody(e.AbsPath, s), " "),
				Time:       GitShowDate(e.AbsPath, s),
				CommitType: RemoteCommit,
			}
			remoteCommits = append(remoteCommits, commit)
		}
	}
	return remoteCommits, nil
}

// this function returns the hashes of the commits that are not pushed to the
// upstream of the specific branch
func (e *RepoEntity) pushDiffsToUpstream() ([]string, error) {
	pushables, err := RevList(e, RevListOptions{
		Ref1: "@{u}",
		Ref2: "HEAD",
	})
	if err != nil {
		return make([]string, 0), nil
	}
	return pushables, nil
}
