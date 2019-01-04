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
func Add(r *git.Repository, file *File, option AddOptions) error {
	addCmdMode = addCmdModeNative
	if option.Update || option.Force || option.DryRun {
		addCmdMode = addCmdModeLegacy
	}
	switch addCmdMode {
	case addCmdModeLegacy:
		err := addWithGit(r, file, option)
		return err
	case addCmdModeNative:
		err := addWithGoGit(r, file)
		return err
	}
	return errors.New("Unhandled add operation")
}

// AddAll function is the wrapper of "git add ." command
func AddAll(r *git.Repository, option AddOptions) error {
	args := make([]string, 0)
	args = append(args, addCommand)
	if option.DryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, ".")
	out, err := GenericGitCommandWithOutput(r.AbsPath, args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

func addWithGit(r *git.Repository, file *File, option AddOptions) error {
	args := make([]string, 0)
	args = append(args, addCommand)
	args = append(args, file.Name)
	if option.Update {
		args = append(args, "--update")
	}
	if option.Force {
		args = append(args, "--force")
	}
	if option.DryRun {
		args = append(args, "--dry-run")
	}
	out, err := GenericGitCommandWithOutput(r.AbsPath, args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

func addWithGoGit(r *git.Repository, file *File) error {
	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(file.Name)
	return nil
}
