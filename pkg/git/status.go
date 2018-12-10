package git

import (
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var statusCommand = "status"

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
	StatusModified  FileStatus = 'M'
	// StatusAdded says file is added to index
	StatusAdded     FileStatus = 'A'
	// StatusDeleted says file is deleted
	StatusDeleted   FileStatus = 'D'
	// StatusRenamed says file is renamed
	StatusRenamed   FileStatus = 'R'
	// StatusCopied says file is copied
	StatusCopied    FileStatus = 'C'
	// StatusUpdated says file is updated
	StatusUpdated   FileStatus = 'U'
	// StatusUntracked says file is untraced
	StatusUntracked FileStatus = '?'
	// StatusIgnored says file is ignored
	StatusIgnored   FileStatus = '!'
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

// LoadFiles function simply commands a git status and collects output in a
// structured way
func (entity *RepoEntity) LoadFiles() ([]*File, error) {
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
