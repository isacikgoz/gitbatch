package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestStatusWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	_, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	var tests = []struct {
		input *git.Repository
	}{
		{th.Repository},
	}
	for _, test := range tests {
		output, err := statusWithGit(test.input)
		require.NoError(t, err)
		require.NotEmpty(t, output)
	}
}

func TestStatusWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	_, err := testFile(th.RepoPath, "file")
	require.NoError(t, err)

	var tests = []struct {
		input *git.Repository
	}{
		{th.Repository},
	}
	for _, test := range tests {
		output, err := statusWithGoGit(test.input)
		require.NoError(t, err)
		require.NotEmpty(t, output)
	}
}
