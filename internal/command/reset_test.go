package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

var (
	testResetopt1 = &ResetOptions{}
)

func TestResetWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	f, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	err = AddAll(th.Repository, testAddopt1)
	require.NoError(t, err)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *git.File
		inp3 *ResetOptions
	}{
		{th.Repository, f, testResetopt1},
	}
	for _, test := range tests {
		err := resetWithGit(test.inp1, test.inp2, test.inp3)
		require.NoError(t, err)
	}
}

func TestResetAllWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	_, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	err = AddAll(th.Repository, testAddopt1)
	require.NoError(t, err)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *ResetOptions
	}{
		{th.Repository, testResetopt1},
	}
	for _, test := range tests {
		err := resetAllWithGit(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}

func TestResetAllWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	_, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	err = AddAll(th.Repository, testAddopt1)
	require.NoError(t, err)

	ref, err := th.Repository.Repo.Head()
	require.NoError(t, err)

	opt := &ResetOptions{
		Hash:      ref.Hash().String(),
		ResetType: ResetMixed,
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *ResetOptions
	}{
		{th.Repository, opt},
	}
	for _, test := range tests {
		err := resetAllWithGoGit(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}
