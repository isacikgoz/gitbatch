package load

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/isacikgoz/gitbatch/internal/testlib"
	"github.com/stretchr/testify/require"
)

func TestSyncLoad(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		input []string
	}{
		{[]string{th.BasicRepoPath(), th.DirtyRepoPath()}},
	}
	for _, test := range tests {
		output, err := SyncLoad(test.input)
		require.NoError(t, err)
		require.NotEmpty(t, output)
	}
}

func TestAsyncLoad(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var loaded []string
	add := func(r *git.Repository) {
		loaded = append(loaded, r.Name)
	}

	dirs := []string{th.BasicRepoPath(), th.DirtyRepoPath(), th.NonRepoPath()}
	err := AsyncLoad(dirs, add)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"basic-repo", "dirty-repo"}, loaded)
}
