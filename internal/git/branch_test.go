package git_test

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/isacikgoz/gitbatch/internal/testlib"
	"github.com/stretchr/testify/require"
)

func TestNextBranch(t *testing.T) {

}

func TestPreviousBranch(t *testing.T) {

}

func TestRevlistNew(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	r := th.Repository
	// HEAD..@{u}
	headref, err := r.Repo.Head()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}

	head := headref.Hash().String()

	pullables, err := git.RevList(r, git.RevListOptions{
		Ref1: head,
		Ref2: r.State.Branch.Upstream.Reference.Hash().String(),
	})
	require.NoError(t, err)
	require.Empty(t, pullables)

	pushables, err := git.RevList(r, git.RevListOptions{
		Ref1: r.State.Branch.Upstream.Reference.Hash().String(),
		Ref2: head,
	})
	require.NoError(t, err)
	require.Empty(t, pushables)
}
