package load

import (
	"fmt"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestSyncLoad(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
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
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	testChannel := make(chan bool)
	testAsyncMockFunc := func(r *git.Repository) {
		go func() {
			if <-testChannel {
				fmt.Println(r.Name)
			}
		}()
	}

	var tests = []struct {
		inp1 []string
		inp2 AsyncAdd
		inp3 chan bool
	}{
		{[]string{th.BasicRepoPath(), th.DirtyRepoPath()}, testAsyncMockFunc, testChannel},
	}
	for _, test := range tests {
		err := AsyncLoad(test.inp1, test.inp2, test.inp3)
		require.NoError(t, err)
	}
}
