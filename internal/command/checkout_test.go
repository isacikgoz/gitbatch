package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
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

	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1 *git.Repository
		inp2 *CheckoutOptions
	}{
		{th.Repository, opts1},
		{th.Repository, opts2},
	}
	for _, test := range tests {
		err := Checkout(test.inp1, test.inp2)
		require.NoError(t, err)
	}
}
