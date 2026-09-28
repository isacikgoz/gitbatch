package git_test

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	stringLength := 8
	randString := git.RandomString(stringLength)
	require.Equal(t, len(randString), stringLength)
}
