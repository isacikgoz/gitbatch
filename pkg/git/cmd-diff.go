package git

import (
	"errors"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	diffCmdMode string

	diffCommand       = "diff"
	diffCmdModeLegacy = "git"
	diffCmdModeNative = "go-git"
)

// Diff is a wrapper function for "git diff" command
// Diff function returns the diff to previous commit detail of the given has
// of a specific commit
func Diff(entity *RepoEntity, hash string) (diff string, err error) {
	diffCmdMode = diffCmdModeNative

	switch diffCmdMode {
	case diffCmdModeLegacy:
		return diffWithGit(entity, hash)
	case diffCmdModeNative:
		return diffWithGoGit(entity, hash)
	}
	return diff, errors.New("Unhandled diff operation")
}

func diffWithGit(entity *RepoEntity, hash string) (diff string, err error) {
	return diff, nil
}

func diffWithGoGit(entity *RepoEntity, hash string) (diff string, err error) {
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
