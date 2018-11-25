package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	"time"
)

var (
	Hashlimit = 6
)

type Commit struct {
	Hash string
	Author string
	Message string
	Time time.Time
}

func newCommit(hash, author, message string, time time.Time) (commit *Commit) {
	commit = &Commit{hash, author, message, time}
	return commit
}

func lastCommit(r *git.Repository) (commit *Commit, err error) {
	ref, err := r.Head()
    if err != nil {
        return nil, err
    }

    cIter, _ := r.Log(&git.LogOptions{
    	From: ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	defer cIter.Close()

	c, err := cIter.Next()
	if err != nil {
        return nil, err
    }
    re := regexp.MustCompile(`\r?\n`)
    commit = newCommit(re.ReplaceAllString(c.Hash.String(), " "), c.Author.Email, re.ReplaceAllString(c.Message, " "), c.Author.When)
	return commit, nil
}

func (entity *RepoEntity) NextCommit() error {

	currentCommit := entity.Commit
	commits, err :=  entity.Commits()
    if err != nil {
        return err
    }

    currentCommitIndex := 0
	for i, cs := range commits {
		if cs.Hash == currentCommit.Hash {
			currentCommitIndex = i
		}
	}
	if currentCommitIndex == len(commits)-1 {
		entity.Commit = commits[0]
		return nil
	}
	entity.Commit = commits[currentCommitIndex+1]
	return nil
}

func (entity *RepoEntity) Commits() (commits []*Commit, err error) {
	r := entity.Repository
	
	ref, err := r.Head()
    if err != nil {
        return commits, err
    }

    cIter, _ := r.Log(&git.LogOptions{
    	From: ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	defer cIter.Close()

    // ... just iterates over the commits
    err = cIter.ForEach(func(c *object.Commit) error {
    	re := regexp.MustCompile(`\r?\n`)
    	commit := newCommit(re.ReplaceAllString(c.Hash.String(), " "), c.Author.Email, re.ReplaceAllString(c.Message, " "), c.Author.When)
        commits = append(commits, commit)

        return nil
	})
	if err != nil {
		return commits, err
	}
    return commits, nil
}

func (entity *RepoEntity) Diff(hash string) (diff string, err error) {

	cms, err :=  entity.Commits()
    if err != nil {
        return "", err
    }

	currentCommitIndex := 0
	for i, cs := range cms {
		if cs.Hash == hash {
			currentCommitIndex = i
		}
	}
	if len(cms) -currentCommitIndex <= 1 {
		return "there is no diff", nil
	}

	commits, err := entity.Repository.Log(&git.LogOptions{
    	From: plumbing.NewHash(cms[currentCommitIndex].Hash),
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