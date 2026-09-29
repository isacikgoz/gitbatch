package git_test

import (
	"sync"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	stringLength := 8
	randString := git.RandomString(stringLength)
	require.Equal(t, len(randString), stringLength)
}

func TestRandomStringConcurrentUnique(t *testing.T) {
	// repositories are loaded in parallel and RepoID must tell them apart
	const n = 1000
	ids := make([]string, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ids[i] = git.RandomString(8)
		}(i)
	}
	wg.Wait()

	seen := make(map[string]struct{}, n)
	for _, id := range ids {
		seen[id] = struct{}{}
	}
	require.Equal(t, n, len(seen), "RandomString returned duplicate values")
}
