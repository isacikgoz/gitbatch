package app

import (
	"os"
	"strings"
	"testing"
)

var (
	wd, _ = os.Getwd()
	sp    = string(os.PathSeparator)
	d     = strings.TrimSuffix(wd, sp+"app")

	relparent = ".." + sp + "test"
	parent    = d + sp + "test"
	// data      = parent + sp + "test-data"
	basic    = testRepoDir + sp + "basic-repo"
	dirty    = testRepoDir + sp + "dirty-repo"
	non      = testRepoDir + sp + "non-repo"
	subbasic = non + sp + "basic-repo"
)

func TestGenerateDirectories(t *testing.T) {
	defer cleanRepo()
	_, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1     []string
		inp2     int
		expected []string
	}{
		{[]string{testRepoDir}, 1, []string{basic, dirty}},
		{[]string{testRepoDir}, 2, []string{basic, dirty}}, // maybe move one repo to a sub folder
	}
	for _, test := range tests {
		if output := generateDirectories(test.inp1, test.inp2); !testEq(output, test.expected) {
			t.Errorf("Test Failed: {%s, %d} inputted, received: %s, expected: %s", test.inp1, test.inp2, output, test.expected)
		}
	}
}

func TestWalkRecursive(t *testing.T) {
	defer cleanRepo()
	_, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 []string
		inp2 []string
		exp1 []string
		exp2 []string
	}{
		{
			[]string{testRepoDir},
			[]string{""},
			[]string{testRepoDir + sp + ".git", testRepoDir + sp + ".gitmodules", non},
			[]string{"", basic, dirty},
		},
	}
	for _, test := range tests {
		if out1, out2 := walkRecursive(test.inp1, test.inp2); !testEq(out1, test.exp1) || !testEq(out2, test.exp2) {
			t.Errorf("Test Failed: {%s, %s} inputted, received: {%s, %s}, expected: {%s, %s}", test.inp1, test.inp2, out1, out2, test.exp1, test.exp2)
		}
	}
}

func TestSeparateDirectories(t *testing.T) {
	defer cleanRepo()
	_, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
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
			testRepoDir,
			[]string{testRepoDir + sp + ".git", testRepoDir + sp + ".gitmodules", non},
			[]string{basic, dirty},
		},
	}
	for _, test := range tests {
		if out1, out2, err := separateDirectories(test.input); !testEq(out1, test.exp1) || !testEq(out2, test.exp2) || err != nil {
			if err != nil {
				t.Errorf("Test failed with error: %s ", err.Error())
				return
			}
			t.Errorf("Test Failed: %s inputted, received: {%s, %s}, expected: {%s, %s}", test.input, out1, out2, test.exp1, test.exp2)
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
