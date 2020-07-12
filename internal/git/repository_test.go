package git

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
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
	_, err := git.PlainClone(testRepoDir, false, &git.CloneOptions{
		URL:               testRepoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	time.Sleep(time.Second)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, err
	}
	return InitializeRepo(testRepoDir)
}

func cleanRepo() error {
	return os.RemoveAll(testRepoDir)
}
