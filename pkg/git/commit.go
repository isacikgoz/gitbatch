package git

import (
	"regexp"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	Hashlimit = 7
)

type Commit struct {
	Hash    string
	Author  string
	Message string
	Time    string
	CommitType    CommitType
}

type CommitType string

const (
	LocalCommit CommitType = "local"
	RemoteCommit CommitType = "remote"
)

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
	for _, rmc := range rmcs {
		entity.Commits = append(entity.Commits, rmc)
	}
	// ... just iterates over the commits
	err = cIter.ForEach(func(c *object.Commit) error {
		re := regexp.MustCompile(`\r?\n`)
		commit := &Commit{
			Hash: re.ReplaceAllString(c.Hash.String(), " "),
			Author: c.Author.Email,
			Message: re.ReplaceAllString(c.Message, " "),
			Time: c.Author.When.String(),
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

	commits, err := entity.Repository.Log(&git.LogOptions{
		From:  plumbing.NewHash(entity.Commit.Hash), //plumbing.NewHash(entity.Commits[currentCommitIndex].Hash),
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

	for _, c := range changes {
		patch, err := c.Patch()
		if err != nil {
			break
		}
		diff = diff + patch.String() + "\n"
	}
	return diff, nil
}
