package git

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// File represents the status of a file in an index or work tree
type File struct {
	Name    string
	AbsPath string
	X       FileStatus
	Y       FileStatus
}

// FileStatus is the short representation of state of a file
type FileStatus byte

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

// Diff is a wrapper of "git diff" command for a file to compare with HEAD rev
func (f *File) Diff() (output string, err error) {
	args := make([]string, 0)
	args = append(args, "diff")
	args = append(args, "HEAD")
	args = append(args, f.Name)
	output, err = GenericGitCommandWithErrorOutput(strings.TrimSuffix(f.AbsPath, f.Name), args)
	if err != nil {
		log.Warn(err)
	}
	return output, err
}
