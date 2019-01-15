package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/core/git"
)

func TestStatusWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	_, err = testFile("file")
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		input *git.Repository
	}{
		{r},
	}
	for _, test := range tests {
		if output, err := statusWithGit(test.input); err != nil || len(output) <= 0 {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestStatusWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	_, err = testFile("file")
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		input *git.Repository
	}{
		{r},
	}
	for _, test := range tests {
		if output, err := statusWithGoGit(test.input); err != nil || len(output) <= 0 {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
