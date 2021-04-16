package command

import (
	"os"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
)

func TestRun(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Test Failed.")
	}
	var tests = []struct {
		inp1 string
		inp2 string
		inp3 []string
	}{
		{wd, "git", []string{"status"}},
	}
	for _, test := range tests {
		if output, err := Run(test.inp1, test.inp2, test.inp3); err != nil || len(output) <= 0 {
			t.Errorf("Test Failed. {%s, %s, %s} inputted, output: %s", test.inp1, test.inp2, test.inp3, output)
		}
	}
}

func TestReturn(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Test Failed.")
	}
	var tests = []struct {
		inp1     string
		inp2     string
		inp3     []string
		expected int
	}{
		{wd, "foo", []string{}, -1},
	}
	for _, test := range tests {
		if output, _ := Return(test.inp1, test.inp2, test.inp3); output != test.expected {
			t.Errorf("Test Failed. {%s, %s, %s} inputted, output: %d, expected : %d", test.inp1, test.inp2, test.inp3, output, test.expected)
		}
	}
}

func TestTrimTrailingNewline(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"foo", "foo"},
		{"foo\n", "foo"},
		{"foo\r", "foo"},
	}
	for _, test := range tests {
		if output := trimTrailingNewline(test.input); output != test.expected {
			t.Errorf("Test Failed. %s inputted, output: %s, expected: %s", test.input, output, test.expected)
		}
	}
}

func testFile(testRepoDir, name string) (*git.File, error) {
	_, err := os.Create(testRepoDir + string(os.PathSeparator) + name)
	if err != nil {
		return nil, err
	}
	f := &git.File{
		Name:    name,
		AbsPath: testRepoDir + string(os.PathSeparator) + name,
		X:       git.StatusUntracked,
		Y:       git.StatusUntracked,
	}
	return f, nil
}
