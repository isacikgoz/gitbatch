package command

import (
	"path/filepath"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestDiffFile(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	f := &git.File{
		AbsPath: filepath.Join(th.RepoPath, ".gitignore"),
		Name:    ".gitignore",
	}

	_, err := testFile(th.RepoPath, f.Name)
	require.NoError(t, err)

	var tests = []struct {
		input    *git.File
		expected string
	}{
		{f, ""},
	}
	for _, test := range tests {
		output, err := DiffFile(test.input)
		require.NoError(t, err)
		require.Equal(t, test.expected, output)
	}
}

func TestDiffWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	headRef, err := th.Repository.Repo.Head()
	require.NoError(t, err)
	var tests = []struct {
		inp1     *git.Repository
		inp2     string
		expected string
	}{
		{th.Repository, headRef.Hash().String(), ""},
	}
	for _, test := range tests {
		output, err := diffWithGoGit(test.inp1, test.inp2)
		require.NoError(t, err)
		require.False(t, len(output) == len(test.expected))
	}
}
