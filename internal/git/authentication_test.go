package git

import (
	"testing"
)

func TestAuthProtocol(t *testing.T) {
	var tests = []struct {
		input *Remote
	}{
		{&Remote{
			URL: []string{"https://gitlab.com/isacikgoz/dirty-repo.git", ""},
		}},
		{&Remote{
			URL: []string{"http://gitlab.com/isacikgoz/dirty-repo.git", ""},
		}},
		{&Remote{
			URL: []string{"git@gitlab.com:isacikgoz/dirty-repo.git", ""},
		}},
	}
	for _, test := range tests {
		if _, err := AuthProtocol(test.input); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
