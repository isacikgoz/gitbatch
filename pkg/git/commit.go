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
func (entity *RepoEntity) NextCommit() error {
	currentCommitIndex := entity.findCurrentCommitIndex()
	if currentCommitIndex == len(entity.Commits)-1 {
		entity.Commit = entity.Commits[0]
		return nil
	}
	entity.Commit = entity.Commits[currentCommitIndex+1]
	return nil
}

// PreviousCommit iterates to opposite direction
func (entity *RepoEntity) PreviousCommit() error {
	currentCommitIndex := entity.findCurrentCommitIndex()
	if currentCommitIndex == 0 {
		entity.Commit = entity.Commits[len(entity.Commits)-1]
		return nil
	}
	entity.Commit = entity.Commits[currentCommitIndex-1]
	return nil
}

// returns the active commit index
func (entity *RepoEntity) findCurrentCommitIndex() int {
	currentCommitIndex := 0
	for i, cs := range entity.Commits {
		if cs.Hash == entity.Commit.Hash {
			currentCommitIndex = i
		}
	}
	return currentCommitIndex
}

// loads the local commits by simply using git log way. ALso, gets the upstream
// diff commits
func (entity *RepoEntity) loadCommits() error {
	r := entity.Repository
	entity.Commits = make([]*Commit, 0)
	ref, err := r.Head()
	if err != nil {
		log.Trace("Cannot get HEAD " + err.Error())
		return err
	}

	cIter, err := r.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		log.Trace("git log failed " + err.Error())
		return err
	}
	defer cIter.Close()
	rmcs, err := entity.pullDiffsToUpstream()
	if err != nil {
		log.Trace("git rev-list failed " + err.Error())
		return err
	}
	for _, rmc := range rmcs {
		entity.Commits = append(entity.Commits, rmc)
	}
	lcs, err := entity.pushDiffsToUpstream()
	if err != nil {
		log.Trace("git rev-list failed " + err.Error())
		return err
	}
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
		entity.Commits = append(entity.Commits, commit)

		return nil
	})
	if err != nil {
		return err
	}
	// entity.Commits = commits
	return nil
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
