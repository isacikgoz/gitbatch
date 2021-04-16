package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/isacikgoz/gitbatch/internal/testlib"
	"github.com/stretchr/testify/require"
)

type TestHelper struct {
	Repository *Repository
	RepoPath   string
}

func InitTestRepositoryFromLocal(t *testing.T) *TestHelper {
	testPathDir, err := ioutil.TempDir("", "gitbatch")
	require.NoError(t, err)

	p, err := testlib.ExtractTestRepository(testPathDir)
	require.NoError(t, err)

	r, err := InitializeRepo(p)
	require.NoError(t, err)

	return &TestHelper{
		Repository: r,
		RepoPath:   p,
	}
}

func InitTestRepository(t *testing.T) *TestHelper {
	testRepoDir, err := ioutil.TempDir("", "test-data")
	require.NoError(t, err)

	testRepoURL := "https://gitlab.com/isacikgoz/test-data.git"
	_, err = git.PlainClone(testRepoDir, false, &git.CloneOptions{
		URL:               testRepoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	time.Sleep(time.Second)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		require.FailNow(t, err.Error())
		return nil
	}

	r, err := InitializeRepo(testRepoDir)
	require.NoError(t, err)

	return &TestHelper{
		Repository: r,
		RepoPath:   testRepoDir,
	}
}

func (h *TestHelper) CleanUp(t *testing.T) {
	err := os.RemoveAll(filepath.Dir(h.RepoPath))
	require.NoError(t, err)
}

func (h *TestHelper) DirtyRepoPath() string {
	return filepath.Join(h.RepoPath, "dirty-repo")
}

func (h *TestHelper) BasicRepoPath() string {
	return filepath.Join(h.RepoPath, "basic-repo")
}

func (h *TestHelper) NonRepoPath() string {
	return filepath.Join(h.RepoPath, "non-repo")
}
