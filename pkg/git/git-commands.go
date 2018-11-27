package git

import (
	"strings"

	"github.com/isacikgoz/gitbatch/pkg/command"
)

// UpstreamDifferenceCount checks how many pushables/pullables there are for the
// current branch
// TODO: get pull pushes to remote branch vs local branch
func UpstreamDifferenceCount(repoPath string) (string, string) {
	args := []string{"rev-list", "@{u}..HEAD", "--count"}
	pushableCount, err := command.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?", "?"
	}
	args = []string{"rev-list", "HEAD..@{u}", "--count"}
	pullableCount, err := command.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

func UpstreamPushDiffs(repoPath string) string {
	args := []string{"rev-list", "@{u}..HEAD"}
	pushableCount, err := command.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return pushableCount
}

func UpstreamPullDiffs(repoPath string) string {
	args := []string{"rev-list", "HEAD..@{u}"}
	pullableCount, err := command.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return pullableCount
}

func GitShow(repoPath, hash string) string {
	args := []string{"show", hash}
	diff, err := command.RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

func (entity *RepoEntity) FetchWithGit(remote string) error {
	args := []string{"fetch", remote}
	_, err := command.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

func (entity *RepoEntity) PullWithGit(remote, branch string) error {
	args := []string{"pull", remote, branch}
	_, err := command.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

func (entity *RepoEntity) MergeWithGit(mergeFrom string) error {
	args := []string{"merge", mergeFrom}
	_, err := command.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

func (entity *RepoEntity) CheckoutWithGit(branch string) error {
	args := []string{"checkout", branch}
	_, err := command.RunCommandWithOutput(entity.AbsPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}
