package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	"log"
)

var (
	Hashlimit = 40
)

func lastCommit(r *git.Repository) (hash string, err error) {
	ref, err := r.Head()
    if err != nil {
        return hash, err
    }

    cIter, _ := r.Log(&git.LogOptions{
    	From: ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	defer cIter.Close()

	c, err := cIter.Next()
	if err != nil {
        return hash, err
    }
    hash = string([]rune(c.Hash.String())[:Hashlimit])
	return hash, nil
}

func (entity *RepoEntity) NextCommit() error {

	currentCommit := entity.Commit
	commits, err :=  entity.Commits()
    if err != nil {
        return err
    }

    currentCommitIndex := 0
	for i, cs := range commits {
		if cs[:Hashlimit] == currentCommit {
			currentCommitIndex = i
		}
	}
	if currentCommitIndex == len(commits)-1 {
		entity.Commit = commits[0][:Hashlimit]
		return nil
	}
	entity.Commit = commits[currentCommitIndex+1][:Hashlimit]
	return nil
}

func (entity *RepoEntity) Commits() (commits []string, err error) {
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
    	commitstring := string([]rune(c.Hash.String())[:Hashlimit]) + " " + c.Message
    	re := regexp.MustCompile(`\r?\n`)
    	commitstring = re.ReplaceAllString(commitstring, " ")
        commits = append(commits, commitstring)

        return nil
	})
	if err != nil {
		return commits, err
	}
    return commits, nil
}

func (entity *RepoEntity) CommitDetail() (commitDetail string, err error) {
	r := entity.Repository
	
	ref, err := r.Head()
    if err != nil {
        return commitDetail, err
    }

    cIter, _ := r.Log(&git.LogOptions{
    	From: ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
    var commit *object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		if string([]rune(c.Hash.String())[:Hashlimit]) == entity.Commit {
			commit = c
		}

        return nil
	})
	// commit, err := cIter.Next()
	// if err != nil {
 //        return commitDetail, err
 //    }
	commitDetail = "Hash: " + commit.Hash.String() + "\n" + "Author: " + commit.Author.Email + "\n" + commit.Message

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
				diff = diff + c.String() + "\n"
			// 	break
			// }
		}

		prevCommit = commit
		prevTree = currentTree
	}
	return diff, nil
}