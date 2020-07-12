package load

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	ggit "github.com/go-git/go-git/v5"
	"github.com/isacikgoz/gitbatch/internal/git"
)

var (
	testChannel = make(chan bool)

	sp       = string(os.PathSeparator)
	basic    = testRepoDir + sp + "basic-repo"
	dirty    = testRepoDir + sp + "dirty-repo"
	non      = testRepoDir + sp + "non-repo"
	subbasic = non + sp + "basic-repo"

	testRepoDir, _ = ioutil.TempDir("", "test-data")
)

func TestSyncLoad(t *testing.T) {
	defer cleanRepo()
	_, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}

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
	defer cleanRepo()
	_, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}

	var tests = []struct {
		inp1 []string
		inp2 AsyncAdd
		inp3 chan bool
	}{
		{[]string{basic, dirty}, testAsyncMockFunc, testChannel},
	}
	for _, test := range tests {
		err := AsyncLoad(test.inp1, test.inp2, test.inp3)
		if err != nil {
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

func testRepo() (*git.Repository, error) {
	testRepoURL := "https://gitlab.com/isacikgoz/test-data.git"
	_, err := ggit.PlainClone(testRepoDir, false, &ggit.CloneOptions{
		URL:               testRepoURL,
		RecurseSubmodules: ggit.DefaultSubmoduleRecursionDepth,
	})
	time.Sleep(time.Second)
	if err != nil && err != ggit.NoErrAlreadyUpToDate {
		return nil, err
	}
	return git.InitializeRepo(testRepoDir)
}

func cleanRepo() error {
	return os.RemoveAll(testRepoDir)
}
