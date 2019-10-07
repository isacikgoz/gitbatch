package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
)

func TestCheckout(t *testing.T) {
	opts1 := &CheckoutOptions{
		TargetRef:      "master",
		CreateIfAbsent: false,
		CommandMode:    ModeLegacy,
	}
	opts2 := &CheckoutOptions{
		TargetRef:      "develop",
		CreateIfAbsent: true,
		CommandMode:    ModeLegacy,
	}
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *CheckoutOptions
	}{
		{r, opts1},
		{r, opts2},
	}
	for _, test := range tests {
		if err := Checkout(test.inp1, test.inp2); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
