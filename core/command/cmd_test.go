package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/isacikgoz/gitbatch/core/git"
	ggit "gopkg.in/src-d/go-git.v4"
)

var (
	testRepoDir, _ = ioutil.TempDir("", "dirty-repo")
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

func testRepo() (*git.Repository, error) {
	testRepoURL := "https://gitlab.com/isacikgoz/dirty-repo.git"
	_, err := ggit.PlainClone(testRepoDir, false, &ggit.CloneOptions{
		URL:               testRepoURL,
		RecurseSubmodules: ggit.DefaultSubmoduleRecursionDepth,
	})
	if err != nil && err != ggit.NoErrAlreadyUpToDate {
		return nil, err
	}
	return git.InitializeRepo(testRepoDir)
}

func testFile(name string) (*git.File, error) {
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

func cleanRepo() error {
	return os.RemoveAll(testRepoDir)
}
