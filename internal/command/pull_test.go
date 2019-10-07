package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
)

var (
	testPullopts1 = &PullOptions{
		RemoteName: "origin",
	}

	testPullopts2 = &PullOptions{
		RemoteName: "origin",
		Force:      true,
	}

	testPullopts3 = &PullOptions{
		RemoteName: "origin",
		Progress:   true,
	}
)

func TestPullWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *PullOptions
	}{
		{r, testPullopts1},
		{r, testPullopts2},
	}
	for _, test := range tests {
		if err := pullWithGit(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestPullWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *PullOptions
	}{
		{r, testPullopts1},
		{r, testPullopts3},
	}
	for _, test := range tests {
		if err := pullWithGoGit(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
