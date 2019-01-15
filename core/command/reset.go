package command

import (
	"errors"

	"github.com/isacikgoz/gitbatch/core/git"
	log "github.com/sirupsen/logrus"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	resetCmdMode string

	resetCommand       = "reset"
	resetCmdModeLegacy = "git"
	resetCmdModeNative = "go-git"
)

// ResetOptions defines the rules of git reset command
type ResetOptions struct {
	// Hash is the reference to be resetted
	Hash string
	// Type is the mode of a reset operation
	Rtype ResetType
}

// ResetType defines a string type for reset git command.
type ResetType string

const (
	// ResetHard Resets the index and working tree. Any changes to tracked
	// files in the working tree since <commit> are discarded.
	ResetHard ResetType = "hard"

	// ResetMixed Resets the index but not the working tree (i.e., the changed
	// files are preserved but not marked for commit) and reports what has not
	// been updated. This is the default action.
	ResetMixed ResetType = "mixed"

	// ResetMerge Resets the index and updates the files in the working tree
	// that are different between <commit> and HEAD, but keeps those which are
	// different between the index and working tree
	ResetMerge ResetType = "merge"

	// ResetSoft Does not touch the index file or the working tree at all
	// (but resets the head to <commit>
	ResetSoft ResetType = "soft"

	// ResetKeep Resets index entries and updates files in the working tree
	// that are different between <commit> and HEAD
	ResetKeep ResetType = "keep"
)

// Reset is the wrapper of "git reset" command
func Reset(r *git.Repository, file *git.File, option *ResetOptions) error {
	resetCmdMode = resetCmdModeLegacy

	switch resetCmdMode {
	case resetCmdModeLegacy:
		err := resetWithGit(r, file, option)
		return err
	case resetCmdModeNative:

	}
	return errors.New("Unhandled reset operation")
}

func resetWithGit(r *git.Repository, file *git.File, option *ResetOptions) error {
	args := make([]string, 0)
	args = append(args, resetCommand)

	args = append(args, "--")
	args = append(args, file.Name)
	if len(option.Rtype) > 0 {
		args = append(args, "--"+string(option.Rtype))
	}
	out, err := Run(r.AbsPath, "git", args)
	if err != nil {
		log.Warn("Error while reset command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

// ResetAll resets the changes in a repository, should be used wise
func ResetAll(r *git.Repository, option *ResetOptions) error {
	resetCmdMode = addCmdModeNative

	switch resetCmdMode {
	case resetCmdModeLegacy:
		err := resetAllWithGit(r, option)
		return err
	case resetCmdModeNative:
		err := resetAllWithGoGit(r, option)
		return err
	}
	return errors.New("Unhandled reset operation")
}

func resetAllWithGit(r *git.Repository, option *ResetOptions) error {
	args := make([]string, 0)
	args = append(args, resetCommand)
	if len(option.Rtype) > 0 {
		args = append(args, "--"+string(option.Rtype))
	}
	out, err := Run(r.AbsPath, "git", args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

func resetAllWithGoGit(r *git.Repository, option *ResetOptions) error {
	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}
	var mode gogit.ResetMode
	switch option.Rtype {
	case ResetHard:
		mode = gogit.HardReset
	case ResetMixed:
		mode = gogit.MixedReset
	case ResetMerge:
		mode = gogit.MergeReset
	case ResetSoft:
		mode = gogit.SoftReset
	}
	err = w.Reset(&gogit.ResetOptions{
		Commit: plumbing.NewHash(option.Hash),
		Mode:   mode,
	})
	return err
}
