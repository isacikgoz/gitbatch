package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	"log"
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
    	// commitstring := string([]rune(c.Hash.String())[:Hashlimit]) + " " + c.Message
    	re := regexp.MustCompile(`\r?\n`)
    	// commitstring = re.ReplaceAllString(commitstring, " ")
    	commit := newCommit(re.ReplaceAllString(c.Hash.String(), " "), c.Author.Email, re.ReplaceAllString(c.Message, " "), c.Author.When)
        commits = append(commits, commit)

        return nil
	})
	if err != nil {
		return commits, err
	}
    return commits, nil
}

func (entity *RepoEntity) CommitDetail() (commitDetail string, err error) {

	commit := entity.Commit
	commitDetail = "Hash: " + commit.Hash + "\n" + "Author: " + commit.Author + "\n" + commit.Message

    return commitDetail, nil
}

// resolve blob at given path from obj. obj can be a commit, tag, tree, or blob.
func resolve(obj object.Object, path string) (*object.Blob, error) {
	switch o := obj.(type) {
	case *object.Commit:
		t, err := o.Tree()
		if err != nil {
			return nil, err
		}
		return resolve(t, path)
	case *object.Tag:
		target, err := o.Object()
		if err != nil {
			return nil, err
		}
		return resolve(target, path)
	case *object.Tree:
		file, err := o.File(path)
		if err != nil {
			return nil, err
		}
		return &file.Blob, nil
	case *object.Blob:
		return o, nil
	default:
		return nil, object.ErrUnsupportedObject
	}
}

func (entity *RepoEntity) Diff(hash string) (diff string, err error) {
	commits, err := entity.Repository.Log(&git.LogOptions{
    	From: plumbing.NewHash(hash),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		log.Fatal(hash)
		return "", err
	}
	defer commits.Close()

	var prevCommit *object.Commit
	var prevTree   *object.Tree

	for {
		commit, err := commits.Next()
		if err != nil {
			break
		}
		currentTree, err := commit.Tree()
		if err != nil {
			return diff, err
		}

		if prevCommit == nil {
			prevCommit = commit
			prevTree = currentTree
			continue
		}

		changes, err := currentTree.Diff(prevTree)
		if err != nil {
			return "", err
		}

		for _, c := range changes {
			// if c.To.Name == node {
				patch, err := c.Patch()
				if err != nil {
					break
				}
				diff = diff + patch.String() + "\n"
			// 	break
			// }
		}

		prevCommit = commit
		prevTree = currentTree
	}
	return diff, nil
}