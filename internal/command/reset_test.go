package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
)

var (
	testResetopt1 = &ResetOptions{}
)

func TestResetWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	f, err := testFile("file")
	AddAll(r, testAddopt1)
	if err != nil {
		t.Errorf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *git.File
		inp3 *ResetOptions
	}{
		{r, f, testResetopt1},
	}
	for _, test := range tests {
		if err := resetWithGit(test.inp1, test.inp2, test.inp3); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestResetAllWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	_, err = testFile("file")
	AddAll(r, testAddopt1)
	if err != nil {
		t.Errorf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *ResetOptions
	}{
		{r, testResetopt1},
	}
	for _, test := range tests {
		if err := resetAllWithGit(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestResetAllWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	_, err = testFile("file")
	AddAll(r, testAddopt1)
	if err != nil {
		t.Errorf("Test Failed. error: %s", err.Error())
	}
	ref, err := r.Repo.Head()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	opt := &ResetOptions{
		Hash:      ref.Hash().String(),
		ResetType: ResetMixed,
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *ResetOptions
	}{
		{r, opt},
	}
	for _, test := range tests {
		if err := resetAllWithGoGit(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
