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
	data      = parent + sp + "test-data"
	basic     = data + sp + "basic-repo"
	dirty     = data + sp + "dirty-repo"
	non       = data + sp + "non-repo"
	subbasic  = non + sp + "basic-repo"
)

func TestGenerateDirectories(t *testing.T) {
	var tests = []struct {
		inp1     []string
		inp2     int
		expected []string
	}{
		{[]string{relparent}, 0, []string{data}},
		{[]string{data}, 0, []string{basic, dirty}},
		{[]string{data}, 2, []string{basic, dirty, subbasic}},
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
			[]string{data},
			[]string{""},
			[]string{data + sp + ".git", data + sp + ".gitmodules", non},
			[]string{"", basic, dirty},
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
			data,
			[]string{data + sp + ".git", data + sp + ".gitmodules", non},
			[]string{basic, dirty},
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
