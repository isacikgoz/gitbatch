package git

import (
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var statusCommand = "status"

type File struct {
	Name string
	AbsPath string
	X FileStatus
	Y FileStatus
}

type FileStatus rune

var (
	StatusNotupdated FileStatus = ' '
	StatusModified FileStatus = 'M'
	StatusAdded FileStatus = 'A'
	StatusDeleted FileStatus = 'D'
	StatusRenamed FileStatus = 'R'
	StatusCopied FileStatus = 'C'
	StatusUntracked FileStatus = '?'
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
		relativePathRegex := regexp.MustCompile(`[(\w|/|.)]+`)
		path := relativePathRegex.FindString(file[2:])

		files = append(files, &File{
			Name: path,
			AbsPath: entity.AbsPath + string(os.PathSeparator) + path,
			X: FileStatus(x),
			Y: FileStatus(y),
			})
	}
	return files, nil
}
