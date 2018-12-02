package git

import (
	"regexp"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

// type of the commit; it can be local or remote (upstream diff)
type CommitType string

const (
	LocalCommit  CommitType = "local"
	RemoteCommit CommitType = "remote"
)

// iterate over next commit of a branch
// TODO: the commits entites can tied to branch instead ot the repository
func (entity *RepoEntity) NextCommit() error {
	currentCommitIndex := 0
	for i, cs := range entity.Commits {
		if cs.Hash == entity.Commit.Hash {
			currentCommitIndex = i
		}
	}
	if currentCommitIndex == len(entity.Commits)-1 {
		entity.Commit = entity.Commits[0]
		return nil
	}
	entity.Commit = entity.Commits[currentCommitIndex+1]
	return nil
}

// loads the local commits by simply using git log way. ALso, gets the upstream
// diff commits
func (entity *RepoEntity) loadCommits() error {
	r := entity.Repository
	entity.Commits = make([]*Commit, 0)
	ref, err := r.Head()
	if err != nil {
		return err
	}

	cIter, err := r.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return err
	}
	defer cIter.Close()
	rmcs, err := entity.pullDiffsToUpstream()
	if err != nil {
		return err
	}
	for _, rmc := range rmcs {
		entity.Commits = append(entity.Commits, rmc)
	}
	// ... just iterates over the commits
	err = cIter.ForEach(func(c *object.Commit) error {
		re := regexp.MustCompile(`\r?\n`)
		commit := &Commit{
			Hash:       re.ReplaceAllString(c.Hash.String(), " "),
			Author:     c.Author.Email,
			Message:    re.ReplaceAllString(c.Message, " "),
			Time:       c.Author.When.String(),
			CommitType: LocalCommit,
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

// returns the diff to previous commit detail of the given hash of a specific
// commit
func (entity *RepoEntity) Diff(hash string) (diff string, err error) {

	currentCommitIndex := 0
	for i, cs := range entity.Commits {
		if cs.Hash == hash {
			currentCommitIndex = i
		}
	}
	if len(entity.Commits)-currentCommitIndex <= 1 {
		return "there is no diff", nil
	}

	// maybe we dont need to log the repo again?
	commits, err := entity.Repository.Log(&git.LogOptions{
		From:  plumbing.NewHash(entity.Commit.Hash),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return "", err
	}

	currentCommit, err := commits.Next()
	if err != nil {
		return "", err
	}
	currentTree, err := currentCommit.Tree()
	if err != nil {
		return diff, err
	}

	prevCommit, err := commits.Next()
	if err != nil {
		return "", err
	}
	prevTree, err := prevCommit.Tree()
	if err != nil {
		return diff, err
	}

	changes, err := prevTree.Diff(currentTree)
	if err != nil {
		return "", err
	}

	// here we collect the actual diff
	for _, c := range changes {
		patch, err := c.Patch()
		if err != nil {
			break
		}
		diff = diff + patch.String() + "\n"
	}
	return diff, nil
}
