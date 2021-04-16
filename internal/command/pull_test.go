package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

var (
	testPullopts1 = &PullOptions{
		RemoteName: "origin",
	}

	testPullopts2 = &PullOptions{
		RemoteName: "origin",
		Force:      true,
	}

	testPullopts3 = &PullOptions{
		RemoteName: "origin",
		Progress:   true,
	}
)

func TestPullWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *PullOptions
	}{
		{th.Repository, testPullopts1},
		{th.Repository, testPullopts2},
	}
	for _, test := range tests {
		err := pullWithGit(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}

func TestPullWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *PullOptions
	}{
		{th.Repository, testPullopts1},
		{th.Repository, testPullopts3},
	}
	for _, test := range tests {
		err := pullWithGoGit(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}
