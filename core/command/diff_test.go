package command

import (
	"os"
	"testing"

	"github.com/isacikgoz/gitbatch/core/git"
)

func TestDiffFile(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	f := &git.File{
		AbsPath: r.AbsPath + string(os.PathSeparator) + ".gitignore",
		Name:    ".gitignore",
	}
	var tests = []struct {
		input    *git.File
		expected string
	}{
		{f, ""},
	}
	for _, test := range tests {
		if output, err := DiffFile(test.input); err != nil || output != test.expected {
			t.Errorf("Test Failed: %s expected, %s was the output.", test.expected, output)
		}
	}
}

func TestDiffWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	headRef, err := r.Repo.Head()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1     *git.Repository
		inp2     string
		expected string
	}{
		{r, headRef.Hash().String(), ""},
	}
	for _, test := range tests {
		if output, err := diffWithGoGit(test.inp1, test.inp2); err != nil || len(output) == len(test.expected) {
			t.Errorf("Test Failed: %s expected, %s was the output.", test.expected, output)
		}
	}
}
