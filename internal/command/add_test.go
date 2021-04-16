package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

var (
	testAddopt1 = &AddOptions{}
)

func TestAddAll(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	_, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *AddOptions
	}{
		{th.Repository, testAddopt1},
	}
	for _, test := range tests {
		err := AddAll(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}

func TestAddWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	f, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *git.File
		inp3 *AddOptions
	}{
		{th.Repository, f, testAddopt1},
	}
	for _, test := range tests {
		err := addWithGit(test.inp1, test.inp2, test.inp3)
		require.NoError(t, err)
	}
}

func TestAddWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	f, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *git.File
	}{
		{th.Repository, f},
	}
	for _, test := range tests {
		err := addWithGoGit(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}
