package command

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/isacikgoz/gitbatch/internal/git"
)

// Diff is a wrapper function for "git diff" command
// Diff function returns the diff to previous commit detail of the given has
// of a specific commit
func Diff(r *git.Repository, hash string) (diff string, err error) {
	// in case we want to configure diff options
	mode := ModeNative

	switch mode {
	case ModeLegacy:
		return diffWithGit(r, hash)
	case ModeNative:
		return diffWithGoGit(r, hash)
	}
	return diff, fmt.Errorf("unhandled diff operation")
}

// DiffFile is a wrapper of "git diff" command for a file to compare with HEAD rev
func DiffFile(f *git.File) (output string, err error) {
	args := make([]string, 0)
	args = append(args, "diff")
	args = append(args, "HEAD")
	args = append(args, f.Name)
	output, err = Run(strings.TrimSuffix(f.AbsPath, f.Name), "git", args)
	if err != nil {
		return "", err
	}
	return output, err
}

// DiffStat shows current working status "git diff --stat"
func DiffStat(r *git.Repository) (string, error) {
	args := make([]string, 0)
	args = append(args, "diff")
	args = append(args, "--stat")
	output, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\n?\r`)
	output = re.ReplaceAllString(output, "\n")
	return output, err
}

// DiffStatRefs shows diff stat of two refs  "git diff a1b2c3..e4f5g6 --stat"
func DiffStatRefs(r *git.Repository, ref1, ref2 string) (string, error) {
	args := make([]string, 0)
	args = append(args, "diff")
	args = append(args, ref1+".."+ref2)
	args = append(args, "--shortstat")
	output, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\n?\r`)
	output = re.ReplaceAllString(output, "\n")
	return output, err
}

// PlainDiff shows current working status "git diff"
func PlainDiff(r *git.Repository) (string, error) {
	args := make([]string, 0)
	args = append(args, "diff")
	output, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\n?\r`)
	output = re.ReplaceAllString(output, "\n")
	return output, err
}

// StashDiff shows diff of stash item "git show stash@{0}"
func StashDiff(r *git.Repository, id int) (string, error) {
	args := make([]string, 0)
	args = append(args, "show")
	args = append(args, "stash@{"+strconv.Itoa(id)+"}")
	output, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\n?\r`)
	output = re.ReplaceAllString(output, "\n")
	return output, err
}

func diffWithGit(r *git.Repository, hash string) (diff string, err error) {
	return diff, nil
}

func diffWithGoGit(r *git.Repository, hash string) (diff string, err error) {
	currentCommitIndex := 0
	for i, cs := range r.State.Branch.Commits {
		if cs.Hash == hash {
			currentCommitIndex = i
		}
	}
	if len(r.State.Branch.Commits)-currentCommitIndex <= 1 {
		return "there is no diff", nil
	}

	// maybe we don't need to log the repo again?
	commits, err := r.Repo.Log(&gogit.LogOptions{
		From:  plumbing.NewHash(r.State.Branch.State.Commit.Hash),
		Order: gogit.LogOrderCommitterTime,
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
