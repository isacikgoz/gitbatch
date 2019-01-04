package command

import (
	"errors"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/isacikgoz/gitbatch/core/git"
	log "github.com/sirupsen/logrus"
)

var (
	statusCmdMode string

	statusCommand       = "status"
	statusCmdModeLegacy = "git"
	statusCmdModeNative = "go-git"
)

func shortStatus(r *git.Repository, option string) string {
	args := make([]string, 0)
	args = append(args, statusCommand)
	args = append(args, option)
	args = append(args, "--short")
	out, err := GenericGitCommandWithOutput(r.AbsPath, args)
	if err != nil {
		log.Warn("Error while status command")
		return "?"
	}
	return out
}

// Status returns the dirty files
func Status(r *git.Repository) ([]*File, error) {
	statusCmdMode = statusCmdModeNative

	switch statusCmdMode {
	case statusCmdModeLegacy:
		return statusWithGit(r)
	case statusCmdModeNative:
		return statusWithGoGit(r)
	}
	return nil, errors.New("Unhandled status operation")
}

// LoadFiles function simply commands a git status and collects output in a
// structured way
func statusWithGit(r *git.Repository) ([]*File, error) {
	files := make([]*File, 0)
	output := shortStatus(r, "--untracked-files=all")
	if len(output) == 0 {
		return files, nil
	}
	fileslist := strings.Split(output, "\n")
	for _, file := range fileslist {
		x := byte(file[0])
		y := byte(file[1])
		relativePathRegex := regexp.MustCompile(`[(\w|/|.|\-)]+`)
		path := relativePathRegex.FindString(file[2:])

		files = append(files, &File{
			Name:    path,
			AbsPath: r.AbsPath + string(os.PathSeparator) + path,
			X:       FileStatus(x),
			Y:       FileStatus(y),
		})
	}
	sort.Sort(filesAlphabetical(files))
	return files, nil
}

func statusWithGoGit(r *git.Repository) ([]*File, error) {
	files := make([]*File, 0)
	w, err := r.Repo.Worktree()
	if err != nil {
		return files, err
	}
	s, err := w.Status()
	if err != nil {
		return files, err
	}
	for k, v := range s {
		files = append(files, &File{
			Name:    k,
			AbsPath: r.AbsPath + string(os.PathSeparator) + k,
			X:       FileStatus(v.Staging),
			Y:       FileStatus(v.Worktree),
		})
	}
	sort.Sort(filesAlphabetical(files))
	return files, nil
}
