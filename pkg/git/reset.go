package git

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

var resetCommand = "reset"

type ResetOptions struct {
	Hard bool
	Merge bool
	Keep bool
}

func (file *File) Reset(option ResetOptions) error {
	args := make([]string, 0)
	args = append(args, resetCommand)
	args = append(args, "--")
	args = append(args, file.Name)
	if option.Hard {
		args = append(args, "--hard")
	}
	if option.Merge {
		args = append(args, "--merge")
	}
	if option.Keep {
		args = append(args, "--keep")
	}
	out, err := GenericGitCommandWithOutput(strings.TrimSuffix(file.AbsPath, file.Name), args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}

func (entity *RepoEntity) ResetAll(option ResetOptions) error {
	args := make([]string, 0)
	args = append(args, resetCommand)
	if option.Hard {
		args = append(args, "--hard")
	}
	out, err := GenericGitCommandWithOutput(entity.AbsPath, args)
	if err != nil {
		log.Warn("Error while add command")
		return errors.New(out + "\n" + err.Error())
	}
	return nil
}