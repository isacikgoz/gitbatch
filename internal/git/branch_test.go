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
	defer func() {
		if err := cleanRepo(); err != nil {
			t.Fatalf("Test Failed. error: %s", err.Error())
		}
	}()

	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	// HEAD..@{u}
	headref, err := r.Repo.Head()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}

	head := headref.Hash().String()
	fmt.Printf("HEAD (%s) @: %s\n", headref.Name(), head)
	fmt.Printf("REMOTE (%s) @ %s\n", r.State.Remote.Name, r.State.Branch.Upstream.Name)
	fmt.Printf("\n")
	pullables, err := RevList(r, RevListOptions{
		Ref1: head,
		Ref2: r.State.Branch.Upstream.Reference.Hash().String(),
	})
	if err != nil {
		t.Errorf("Test Failed.")
	}
	for _, pullable := range pullables {
		fmt.Println(pullable.Hash.String())
	}
	fmt.Printf("\n")
	pushables, err := RevList(r, RevListOptions{
		Ref1: r.State.Branch.Upstream.Reference.Hash().String(),
		Ref2: head,
	})
	if err != nil {
		t.Errorf("Test Failed.")
	}
	for _, pushable := range pushables {
		fmt.Println(pushable.Hash.String())
	}
}
