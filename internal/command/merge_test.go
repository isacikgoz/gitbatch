package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	opts := &MergeOptions{
		BranchName: th.Repository.State.Branch.Upstream.Name,
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *MergeOptions
	}{
		{th.Repository, opts},
	}
	for _, test := range tests {
		err := Merge(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}
