package command

import (
	"testing"

	"github.com/stretchr/testify/require"
	giterr "github.com/thorstenhirsch/gitbatch/internal/errors"
	"github.com/thorstenhirsch/gitbatch/internal/git"
)

func TestCommitWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	f, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	err = addWithGit(th.Repository, f, testAddopt1)
	require.NoError(t, err)

	testCommitopt1 := &CommitOptions{
		CommitMsg: "test",
		User:      "foo",
		Email:     "foo@bar.com",
	}

	var tests = []struct {
		inp1 *git.Repository
		inp2 *CommitOptions
	}{
		{th.Repository, testCommitopt1},
	}
	for _, test := range tests {
		err = commitWithGit(test.inp1, test.inp2)
		require.False(t, err != nil && err == giterr.ErrUserEmailNotSet)
	}
}

func TestCommitWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	f, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	err = addWithGit(th.Repository, f, testAddopt1)
	require.NoError(t, err)

	testCommitopt1 := &CommitOptions{
		CommitMsg: "test",
		User:      "foo",
		Email:     "foo@bar.com",
	}

	var tests = []struct {
		inp1 *git.Repository
		inp2 *CommitOptions
	}{
		{th.Repository, testCommitopt1},
	}
	for _, test := range tests {
		err = commitWithGoGit(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}
