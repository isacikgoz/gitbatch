package app

import (
	"os"
	"strings"
	"testing"
)

var (
	wd, _ = os.Getwd()
	d     = strings.TrimSuffix(wd, string(os.PathSeparator)+"app")
)

func TestGenerateDirectories(t *testing.T) {
	var tests = []struct {
		inp1     []string
		inp2     int
		expected []string
	}{
		{[]string{"../test"}, 0, []string{d + "/test/test-data"}},
		{[]string{"../test/test-data"}, 0, []string{d + "/test/test-data/basic-repo", d + "/test/test-data/dirty-repo"}},
		{[]string{"../test/test-data"}, 2, []string{d + "/test/test-data/basic-repo", d + "/test/test-data/dirty-repo", d + "/test/test-data/non-repo/basic-repo"}},
	}
	for _, test := range tests {
		if output := generateDirectories(test.inp1, test.inp2); !testEq(output, test.expected) {
			t.Errorf("Test Failed: {%s, %d} inputted, recieved: %s, expected: %s", test.inp1, test.inp2, output, test.expected)
		}
	}
}

func TestWalkRecursive(t *testing.T) {
	var tests = []struct {
		inp1 []string
		inp2 []string
		exp1 []string
		exp2 []string
	}{
		{
			[]string{"../test/test-data"},
			[]string{""},
			[]string{d + "/test/test-data/.git", d + "/test/test-data/.gitmodules", d + "/test/test-data/non-repo"},
			[]string{"", d + "/test/test-data/basic-repo", d + "/test/test-data/dirty-repo"},
		},
	}
	for _, test := range tests {
		if out1, out2 := walkRecursive(test.inp1, test.inp2); !testEq(out1, test.exp1) || !testEq(out2, test.exp2) {
			t.Errorf("Test Failed: {%s, %s} inputted, recieved: {%s, %s}, expected: {%s, %s}", test.inp1, test.inp2, out1, out2, test.exp1, test.exp2)
		}
	}
}

func TestSeperateDirectories(t *testing.T) {
	var tests = []struct {
		input string
		exp1  []string
		exp2  []string
	}{
		{
			"",
			nil,
			nil,
		},
		{
			"../test/test-data",
			[]string{d + "/test/test-data/.git", d + "/test/test-data/.gitmodules", d + "/test/test-data/non-repo"},
			[]string{d + "/test/test-data/basic-repo", d + "/test/test-data/dirty-repo"},
		},
	}
	for _, test := range tests {
		if out1, out2, err := seperateDirectories(test.input); !testEq(out1, test.exp1) || !testEq(out2, test.exp2) || err != nil {
			if err != nil {
				t.Errorf("Test failed with error: %s ", err.Error())
				return
			}
			t.Errorf("Test Failed: %s inputted, recieved: {%s, %s}, expected: {%s, %s}", test.input, out1, out2, test.exp1, test.exp2)
		}
	}
}

func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
