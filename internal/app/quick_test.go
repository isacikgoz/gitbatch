package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thorstenhirsch/gitbatch/internal/git"
)

func TestQuick(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1 []string
		inp2 string
	}{
		{
			[]string{th.DirtyRepoPath()},
			"fetch",
		}, {
			[]string{th.DirtyRepoPath()},
			"pull",
		},
	}
	for _, test := range tests {
		err := quick(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}
