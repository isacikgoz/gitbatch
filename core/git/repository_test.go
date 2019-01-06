package git

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	ggit "gopkg.in/src-d/go-git.v4"
)

var (
	testRepoDir, _ = ioutil.TempDir("", "dirty-repo")
)

func TestInitializeRepo(t *testing.T) {
	defer cleanRepo()
	_, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
}

func testRepo() (*Repository, error) {
	testRepoURL := "https://gitlab.com/isacikgoz/dirty-repo.git"
	_, err := ggit.PlainClone(testRepoDir, false, &ggit.CloneOptions{
		URL:               testRepoURL,
		RecurseSubmodules: ggit.DefaultSubmoduleRecursionDepth,
	})
	time.Sleep(time.Second)
	if err != nil && err != ggit.NoErrAlreadyUpToDate {
		return nil, err
	}
	return InitializeRepo(testRepoDir)
}

func cleanRepo() error {
	return os.RemoveAll(testRepoDir)
}
