package git

import (
	"errors"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	statusCmdMode string

	statusCommand       = "status"
	statusCmdModeLegacy = "git"
	statusCmdModeNative = "go-git"
)

// File represents the status of a file in an index or work tree
type File struct {
	Name    string
	AbsPath string
	X       FileStatus
	Y       FileStatus
}

// FileStatus is the short representation of state of a file
type FileStatus rune

var (
	// StatusNotupdated says file not updated
	StatusNotupdated FileStatus = ' '
	// StatusModified says file is modifed
	StatusModified FileStatus = 'M'
	// StatusAdded says file is added to index
	StatusAdded FileStatus = 'A'
	// StatusDeleted says file is deleted
	StatusDeleted FileStatus = 'D'
	// StatusRenamed says file is renamed
	StatusRenamed FileStatus = 'R'
	// StatusCopied says file is copied
	StatusCopied FileStatus = 'C'
	// StatusUpdated says file is updated
	StatusUpdated FileStatus = 'U'
	// StatusUntracked says file is untraced
	StatusUntracked FileStatus = '?'
	// StatusIgnored says file is ignored
	StatusIgnored FileStatus = '!'
)

func shortStatus(entity *RepoEntity, option string) string {
	args := make([]string, 0)
	args = append(args, statusCommand)
	args = append(args, option)
	args = append(args, "--short")
	out, err := GenericGitCommandWithOutput(entity.AbsPath, args)
	if err != nil {
		log.Warn("Error while status command")
		return "?"
	}
	return out
}

func Status(entity *RepoEntity) ([]*File, error) {
	statusCmdMode = statusCmdModeNative

	switch statusCmdMode {
	case statusCmdModeLegacy:
		return statusWithGit(entity)
	case statusCmdModeNative:
		return statusWithGoGit(entity)
	}
	return nil, errors.New("Unhandled status operation")
}

// LoadFiles function simply commands a git status and collects output in a
// structured way
func statusWithGit(entity *RepoEntity) ([]*File, error) {
	files := make([]*File, 0)
	output := shortStatus(entity, "--untracked-files=all")
	if len(output) == 0 {
		return files, nil
	}
	fileslist := strings.Split(output, "\n")
	for _, file := range fileslist {
		x := rune(file[0])
		y := rune(file[1])
		relativePathRegex := regexp.MustCompile(`[(\w|/|.|\-)]+`)
		path := relativePathRegex.FindString(file[2:])

		files = append(files, &File{
			Name:    path,
			AbsPath: entity.AbsPath + string(os.PathSeparator) + path,
			X:       FileStatus(x),
			Y:       FileStatus(y),
		})
	}
	return files, nil
}

func statusWithGoGit(entity *RepoEntity) ([]*File, error) {
	files := make([]*File, 0)
	w, err := entity.Repository.Worktree()
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
			AbsPath: entity.AbsPath + string(os.PathSeparator) + k,
			X:       FileStatus(v.Staging),
			Y:       FileStatus(v.Worktree),
		})
	}
	return files, nil
}

// Diff is a wrapper of "git diff" command for a file to compare with HEAD rev
func (file *File) Diff() (output string, err error) {
	args := make([]string, 0)
	args = append(args, "diff")
	args = append(args, "HEAD")
	args = append(args, file.Name)
	output, err = GenericGitCommandWithErrorOutput(strings.TrimSuffix(file.AbsPath, file.Name), args)
	return output, err
}
