package git

import (
	"fmt"
	"testing"
)

func TestNextBranch(t *testing.T) {

}

func TestPreviousBranch(t *testing.T) {

}

func TestCheckout(t *testing.T) {

}

func TestRevlistNew(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	// HEAD..@{u}
	upstream := r.State.Remote.Branch.Reference.Hash().String()
	headref, err := r.Repo.Head()
	head := headref.Hash().String()
	fmt.Printf("HEAD (%s) @: %s\n", headref.Name(), head)
	fmt.Printf("REMOTE (%s) @ %s\n", r.State.Remote.Branch.Name, upstream)
	fmt.Printf("\n")
	pullables, err := RevList(r, RevListOptions{
		Ref1: head,
		Ref2: upstream,
	})
	if err != nil {
		t.Errorf("Test Failed.")
	}
	for _, pullable := range pullables {
		fmt.Println(pullable.Hash.String())
	}
	fmt.Printf("\n")
	pushables, err := RevList(r, RevListOptions{
		Ref1: upstream,
		Ref2: head,
	})
	if err != nil {
		t.Errorf("Test Failed.")
	}
	for _, pushable := range pushables {
		fmt.Println(pushable.Hash.String())
	}
}

func testRevRepo() (*Repository, error) {
	return InitializeRepo("/home/isacikgoz/git-testing/gitbatch")
}
