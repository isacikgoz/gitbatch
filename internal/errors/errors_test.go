package errors

import (
	"testing"
)

func TestParseGitError(t *testing.T) {
	var tests = []struct {
		input    string
		expected error
	}{
		{"", ErrUnclassified},
	}
	for _, test := range tests {
		if output := ParseGitError(test.input, nil); output != test.expected {
			t.Errorf("Test Failed. %s expected, output: %s", test.expected.Error(), output.Error())
		}
	}
}
