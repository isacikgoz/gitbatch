package git

import (
	"strings"

	"github.com/isacikgoz/gitbatch/pkg/helpers"
)

// UpstreamDifferenceCount checks how many pushables/pullables there are for the
// current branch
// TODO: get pull pushes to remote branch vs local branch
func UpstreamDifferenceCount(repoPath string) (string, string) {
	args := []string{"rev-list", "@{u}..HEAD", "--count"}
	pushableCount, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?", "?"
	}
	args = []string{"rev-list", "HEAD..@{u}", "--count"}
	pullableCount, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

// Instead of returning the count, this method returns the hash list
func UpstreamPushDiffs(repoPath string) string {
	args := []string{"rev-list", "@{u}..HEAD"}
	pushableCount, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return pushableCount
}

// Instead of returning the count, this method returns the hash list
func UpstreamPullDiffs(repoPath string) string {
	args := []string{"rev-list", "HEAD..@{u}"}
	pullableCount, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return pullableCount
}

// Conventional git show command without any argument
func GitShow(repoPath, hash string) string {
	args := []string{"show", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// get author's e-mail with git show command
func GitShowEmail(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ae", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// get body of the commit with git show
func GitShowBody(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%B", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return err.Error()
	}
	return diff
}

// get commit's date with git show as string
func GitShowDate(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ai", hash}
	diff, err := helpers.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// wrapper of the git fetch <remote> command
func (entity *RepoEntity) FetchWithGit(remote string) error {
	args := []string{"fetch", remote}
	_, err := helpers.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

// wrapper of the git pull <remote>/<branch> command
func (entity *RepoEntity) PullWithGit(remote, branch string) error {
	args := []string{"pull", remote, branch}
	_, err := helpers.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

// wrapper of the git merge <branch> command
func (entity *RepoEntity) MergeWithGit(mergeFrom string) error {
	args := []string{"merge", mergeFrom}
	_, err := helpers.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

// wrapper of the git checkout <branch> command
func (entity *RepoEntity) CheckoutWithGit(branch string) error {
	args := []string{"checkout", branch}
	_, err := helpers.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

// GitStatus returns the plaintext short status of the repo
func (entity *RepoEntity) StatusWithGit() string {
	args := []string{"status"}
	status, err := helpers.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return "?"
	}
	return status
}
