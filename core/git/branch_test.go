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

func TestRevlist(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	pullables, err := RevList(r, RevListOptions{
		Ref1: "HEAD",
		Ref2: "@{u}",
	})
	if err != nil {
		t.Errorf("Test Failed.")
	}
	for _, pullable := range pullables {
		fmt.Println(pullable)
	}
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
	pullables, err := RevListNew(r, RevListOptions{
		Ref1: head,
		Ref2: upstream,
	})
	if err != nil {
		t.Errorf("Test Failed.")
	}
	for _, pullable := range pullables {
		fmt.Println(pullable)
	}
}

func testRevRepo() (*Repository, error) {
	return InitializeRepo("/home/isacikgoz/test-git/git-1/testing")
}
