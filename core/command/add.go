package command

import (
	"errors"

	"github.com/isacikgoz/gitbatch/core/git"
	log "github.com/sirupsen/logrus"
)

var (
	addCmdMode string

	addCommand       = "add"
	addCmdModeLegacy = "git"
	addCmdModeNative = "go-git"
)

// AddOptions defines the rules for "git add" command
type AddOptions struct {
	// Update
	Update bool
	// Force
	Force bool
	// DryRun
	DryRun bool
}

// Add is a wrapper function for "git add" command
func Add(r *git.Repository, f *git.File, o *AddOptions) error {
	addCmdMode = addCmdModeNative
	if o.Update || o.Force || o.DryRun {
		addCmdMode = addCmdModeLegacy
	}
	switch addCmdMode {
	case addCmdModeLegacy:
		err := addWithGit(r, f, o)
		return err
	case addCmdModeNative:
		err := addWithGoGit(r, f)
		return err
	}
	return errors.New("Unhandled add operation")
}

// AddAll function is the wrapper of "git add ." command
func AddAll(r *git.Repository, o *AddOptions) error {
	args := make([]string, 0)
	args = append(args, addCommand)
	if o.DryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, ".")
	out, err := Run(r.AbsPath, "git", args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

func addWithGit(r *git.Repository, f *git.File, o *AddOptions) error {
	args := make([]string, 0)
	args = append(args, addCommand)
	args = append(args, f.Name)
	if o.Update {
		args = append(args, "--update")
	}
	if o.Force {
		args = append(args, "--force")
	}
	if o.DryRun {
		args = append(args, "--dry-run")
	}
	out, err := Run(r.AbsPath, "git", args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

func addWithGoGit(r *git.Repository, f *git.File) error {
	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(f.Name)
	return err
}
