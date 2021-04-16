package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

var (
	testFetchopts1 = &FetchOptions{
		RemoteName: "origin",
	}

	testFetchopts2 = &FetchOptions{
		RemoteName: "origin",
		Prune:      true,
	}

	testFetchopts3 = &FetchOptions{
		RemoteName: "origin",
		DryRun:     true,
	}

	testFetchopts4 = &FetchOptions{
		RemoteName: "origin",
		Progress:   true,
	}
)

func TestFetchWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *FetchOptions
	}{
		{th.Repository, testFetchopts1},
		{th.Repository, testFetchopts2},
		{th.Repository, testFetchopts3},
	}
	for _, test := range tests {
		err := fetchWithGit(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}

func TestFetchWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	refspec := "+" + "refs/heads/" + th.Repository.State.Branch.Name + ":" + "/refs/remotes/" + th.Repository.State.Branch.Name
	var tests = []struct {
		inp1 *git.Repository
		inp2 *FetchOptions
		inp3 string
	}{
		{th.Repository, testFetchopts1, refspec},
		{th.Repository, testFetchopts4, refspec},
	}
	for _, test := range tests {
		err := fetchWithGoGit(test.inp1, test.inp2, test.inp3)
		require.NoError(t, err)
	}
}
