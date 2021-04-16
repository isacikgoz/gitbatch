package command

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/isacikgoz/gitbatch/internal/git"
)

func shortStatus(r *git.Repository, option string) string {
	args := make([]string, 0)
	args = append(args, "status")
	args = append(args, option)
	args = append(args, "--short")
	out, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return "?"
	}
	return out
}

// Status returns the dirty files
func Status(r *git.Repository) ([]*git.File, error) {
	// in case we want configure Status command externally
	mode := ModeLegacy

	switch mode {
	case ModeLegacy:
		return statusWithGit(r)
	case ModeNative:
		return statusWithGoGit(r)
	}
	return nil, fmt.Errorf("unhandled status operation")
}

// PlainStatus returns the plain status
func PlainStatus(r *git.Repository) (string, error) {
	args := make([]string, 0)
	args = append(args, "status")
	output, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\n?\r`)
	output = re.ReplaceAllString(output, "\n")
	return output, err
}

// LoadFiles function simply commands a git status and collects output in a
// structured way
func statusWithGit(r *git.Repository) ([]*git.File, error) {
	files := make([]*git.File, 0)
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

		files = append(files, &git.File{
			Name:    path,
			AbsPath: r.AbsPath + string(os.PathSeparator) + path,
			X:       git.FileStatus(x),
			Y:       git.FileStatus(y),
		})
	}
	sort.Sort(git.FilesAlphabetical(files))
	return files, nil
}

func statusWithGoGit(r *git.Repository) ([]*git.File, error) {
	files := make([]*git.File, 0)
	w, err := r.Repo.Worktree()
	if err != nil {
		return files, err
	}
	s, err := w.Status()
	if err != nil {
		return files, err
	}
	for k, v := range s {
		files = append(files, &git.File{
			Name:    k,
			AbsPath: r.AbsPath + string(os.PathSeparator) + k,
			X:       git.FileStatus(v.Staging),
			Y:       git.FileStatus(v.Worktree),
		})
	}
	sort.Sort(git.FilesAlphabetical(files))
	return files, nil
}
