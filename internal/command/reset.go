package command

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/internal/git"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// ResetOptions defines the rules of git reset command
type ResetOptions struct {
	// Hash is the reference to be resetted
	Hash string
	// Type is the mode of a reset operation
	ResetType ResetType
	// Mode is the command mode
	CommandMode Mode
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
func Reset(r *git.Repository, file *git.File, o *ResetOptions) error {
	mode := o.CommandMode

	switch mode {
	case ModeLegacy:
		err := resetWithGit(r, file, o)
		return err
	case ModeNative:

	}
	return fmt.Errorf("unhandled reset operation")
}

func resetWithGit(r *git.Repository, file *git.File, option *ResetOptions) error {
	args := make([]string, 0)
	args = append(args, "reset")

	args = append(args, "--")
	args = append(args, file.Name)
	if len(option.ResetType) > 0 {
		args = append(args, "--"+string(option.ResetType))
	}
	_, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return fmt.Errorf("could not reset file %s: %v", file.AbsPath, err)
	}
	return nil
}

// ResetAll resets the changes in a repository, should be used wise
func ResetAll(r *git.Repository, o *ResetOptions) error {

	switch o.CommandMode {
	case ModeLegacy:
		err := resetAllWithGit(r, o)
		return err
	case ModeNative:
		err := resetAllWithGoGit(r, o)
		return err
	}
	return fmt.Errorf("unhandled reset operation")
}

func resetAllWithGit(r *git.Repository, option *ResetOptions) error {
	args := make([]string, 0)
	args = append(args, "reset")
	if len(option.ResetType) > 0 {
		args = append(args, "--"+string(option.ResetType))
	}
	_, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return fmt.Errorf("could not reset all: %v", err)
	}
	return nil
}

func resetAllWithGoGit(r *git.Repository, option *ResetOptions) error {
	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}
	var mode gogit.ResetMode
	switch option.ResetType {
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
