package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
)

func TestMerge(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	opts := &MergeOptions{
		BranchName: r.State.Branch.Upstream.Name,
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *MergeOptions
	}{
		{r, opts},
	}
	for _, test := range tests {
		if err := Merge(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
