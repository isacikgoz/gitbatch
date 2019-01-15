package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/core/git"
)

var (
	testFetchopts1 = &FetchOptions{
		RemoteName: "origin",
	}

	testFetchopts2 = &FetchOptions{
		RemoteName: "origin",
		Prune:      true,
	}

	testFetchopts3 = &FetchOptions{
		RemoteName: "origin",
		DryRun:     true,
	}

	testFetchopts4 = &FetchOptions{
		RemoteName: "origin",
		Progress:   true,
	}
)

func TestFetchWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *FetchOptions
	}{
		{r, testFetchopts1},
		{r, testFetchopts2},
		{r, testFetchopts3},
	}
	for _, test := range tests {
		if err := fetchWithGit(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestFetchWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	refspec := "+" + "refs/heads/" + r.State.Branch.Name + ":" + "/refs/remotes/" + r.State.Branch.Name
	var tests = []struct {
		inp1 *git.Repository
		inp2 *FetchOptions
		inp3 string
	}{
		{r, testFetchopts1, refspec},
		{r, testFetchopts4, refspec},
	}
	for _, test := range tests {
		if err := fetchWithGoGit(test.inp1, test.inp2, test.inp3); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
