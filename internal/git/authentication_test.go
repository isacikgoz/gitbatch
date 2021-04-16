package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthProtocol(t *testing.T) {
	var tests = []struct {
		input    *Remote
		expected string
	}{
		{&Remote{
			URL: []string{"https://gitlab.com/isacikgoz/dirty-repo.git", ""},
		}, "https"},
		{&Remote{
			URL: []string{"http://gitlab.com/isacikgoz/dirty-repo.git", ""},
		}, "http"},
		{&Remote{
			URL: []string{"git@gitlab.com:isacikgoz/dirty-repo.git", ""},
		}, "ssh"},
	}
	for _, test := range tests {
		protocol, err := AuthProtocol(test.input)
		require.NoError(t, err)
		require.Equal(t, test.expected, protocol)
	}
}
