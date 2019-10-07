package job

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	ggit "gopkg.in/src-d/go-git.v4"
)

var (
	testRepoDir, _ = ioutil.TempDir("", "dirty-repo")
)

func TestStart(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var (
		mockJob1 = &Job{
			JobType:    PullJob,
			Repository: r,
		}
		mockJob2 = &Job{
			JobType:    FetchJob,
			Repository: r,
		}
		mockJob3 = &Job{
			JobType:    MergeJob,
			Repository: r,
		}
	)
	var tests = []struct {
		input *Job
	}{
		{mockJob1},
		{mockJob2},
		{mockJob3},
	}
	for _, test := range tests {
		if err := test.input.start(); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func testRepo() (*git.Repository, error) {
	testRepoURL := "https://gitlab.com/isacikgoz/dirty-repo.git"
	_, err := ggit.PlainClone(testRepoDir, false, &ggit.CloneOptions{
		URL:               testRepoURL,
		RecurseSubmodules: ggit.DefaultSubmoduleRecursionDepth,
	})
	if err != nil && err != ggit.NoErrAlreadyUpToDate {
		return nil, err
	}
	return git.InitializeRepo(testRepoDir)
}

func cleanRepo() error {
	return os.RemoveAll(testRepoDir)
}
