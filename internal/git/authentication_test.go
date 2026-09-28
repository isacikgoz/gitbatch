package git_test

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestAuthProtocol(t *testing.T) {
	var tests = []struct {
		input    *git.Remote
		expected string
	}{
		{&git.Remote{
			URL: []string{"https://gitlab.com/isacikgoz/dirty-repo.git", ""},
		}, "https"},
		{&git.Remote{
			URL: []string{"http://gitlab.com/isacikgoz/dirty-repo.git", ""},
		}, "http"},
		{&git.Remote{
			URL: []string{"git@gitlab.com:isacikgoz/dirty-repo.git", ""},
		}, "ssh"},
	}
	for _, test := range tests {
		protocol, err := git.AuthProtocol(test.input)
		require.NoError(t, err)
		require.Equal(t, test.expected, protocol)
	}
}
