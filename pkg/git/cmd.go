package git

import (
	"github.com/isacikgoz/gitbatch/pkg/helpers"
)

// GenericGitCommand runs any git command without expecting output
func GenericGitCommand(repoPath string, args []string) error {
	_, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

// GenericGitCommandWithOutput runs any git command with returning output
func GenericGitCommandWithOutput(repoPath string, args []string) (string, error) {
	out, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?", err
	}
	return helpers.TrimTrailingNewline(out), nil
}

// GenericGitCommandWithErrorOutput runs any git command with returning output
func GenericGitCommandWithErrorOutput(repoPath string, args []string) (string, error) {
	out, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return helpers.TrimTrailingNewline(out), err
	}
	return helpers.TrimTrailingNewline(out), nil
}

// GitShow is conventional git show command without any argument
func GitShow(repoPath, hash string) string {
	args := []string{"show", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// GitShowEmail gets author's e-mail with git show command
func GitShowEmail(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ae", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// GitShowBody gets body of the commit with git show
func GitShowBody(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%B", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return err.Error()
	}
	return diff
}

// GitShowDate gets commit's date with git show as string
func GitShowDate(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ai", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// StatusWithGit returns the plaintext short status of the repo
func (entity *RepoEntity) StatusWithGit() string {
	args := []string{"status"}
	status, err := helpers.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return "?"
	}
	return status
}
