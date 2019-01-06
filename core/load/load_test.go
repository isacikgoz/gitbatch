package load

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/isacikgoz/gitbatch/core/git"
)

var (
	testChannel = make(chan bool)
)

func TestSyncLoad(t *testing.T) {
	var (
		wd, _ = os.Getwd()
		sp    = string(os.PathSeparator)
		d     = strings.TrimSuffix(wd, sp+"core"+sp+"load")

		parent = d + sp + "test"
		data   = parent + sp + "test-data"
		basic  = data + sp + "basic-repo"
		dirty  = data + sp + "dirty-repo"
	)
	var tests = []struct {
		input []string
	}{
		{[]string{basic, dirty}},
	}
	for _, test := range tests {
		if output, err := SyncLoad(test.input); err != nil || len(output) <= 0 {
			t.Errorf("Test Failed. %s inputted, found %d repos.", test.input, len(output))
		}
	}
}

func TestAsyncLoad(t *testing.T) {
	var (
		wd, _ = os.Getwd()
		sp    = string(os.PathSeparator)
		d     = strings.TrimSuffix(wd, sp+"core"+sp+"load")

		parent = d + sp + "test"
		data   = parent + sp + "test-data"
		basic  = data + sp + "basic-repo"
		dirty  = data + sp + "dirty-repo"
	)
	var tests = []struct {
		inp1 []string
		inp2 AsyncAdd
		inp3 chan bool
	}{
		{[]string{basic, dirty}, testAsyncMockFunc, testChannel},
	}
	for _, test := range tests {
		if err := AsyncLoad(test.inp1, test.inp2, test.inp3); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func testAsyncMockFunc(r *git.Repository) {
	go func() {
		if <-testChannel {
			fmt.Println(r.Name)
		}
	}()
}
