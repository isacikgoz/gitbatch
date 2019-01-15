package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/core/git"
)

var (
	testAddopt1 = &AddOptions{}
)

func TestAddAll(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	_, err = testFile("file")
	if err != nil {
		t.Errorf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *AddOptions
	}{
		{r, testAddopt1},
	}
	for _, test := range tests {
		if err := AddAll(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestAddWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	f, err := testFile("file")
	if err != nil {
		t.Errorf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *git.File
		inp3 *AddOptions
	}{
		{r, f, testAddopt1},
	}
	for _, test := range tests {
		if err := addWithGit(test.inp1, test.inp2, test.inp3); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestAddWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	f, err := testFile("file")
	if err != nil {
		t.Errorf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *git.File
	}{
		{r, f},
	}
	for _, test := range tests {
		if err := addWithGoGit(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
