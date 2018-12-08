package git

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

var addCommand = "add"

type AddOptions struct {
	Update bool
	Force bool
	DryRun bool
}

func (file *File) Add(option AddOptions) error {
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
	out, err := GenericGitCommandWithOutput(strings.TrimSuffix(file.AbsPath, file.Name), args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

func (entity *RepoEntity) AddAll(option AddOptions) error {
	args := make([]string, 0)
	args = append(args, addCommand)
	if option.DryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, ".")
	out, err := GenericGitCommandWithOutput(entity.AbsPath, args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}