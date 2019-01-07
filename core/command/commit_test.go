package command

import (
	"testing"

	giterr "github.com/isacikgoz/gitbatch/core/errors"
	"github.com/isacikgoz/gitbatch/core/git"
)

var (
	testCommitopt1 = &CommitOptions{
		CommitMsg: "test",
		User:      "foo",
		Email:     "foo@bar.com",
	}
)

func TestCommitWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	f, err := testFile("file")
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	if err := addWithGit(r, f, testAddopt1); err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *CommitOptions
	}{
		{r, testCommitopt1},
	}
	for _, test := range tests {
		if err := commitWithGit(test.inp1, test.inp2); err != nil && err == giterr.ErrUserEmailNotSet {
			t.Errorf("Test Failed.")
		}
	}
}

func TestCommitWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	f, err := testFile("file")
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	if err := addWithGit(r, f, testAddopt1); err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *CommitOptions
	}{
		{r, testCommitopt1},
	}
	for _, test := range tests {
		if err := commitWithGoGit(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed.")
		}
	}
}
