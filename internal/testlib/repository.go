package testlib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

// TestHelper holds a repository extracted from the embedded test fixture along
// with the path it was extracted to.
type TestHelper struct {
	Repository *git.Repository
	RepoPath   string
}

// InitTestRepositoryFromLocal extracts the embedded test fixture into a
// temporary directory and initializes a repository from it.
func InitTestRepositoryFromLocal(t *testing.T) *TestHelper {
	testPathDir, err := os.MkdirTemp("", "gitbatch")
	require.NoError(t, err)

	p, err := ExtractTestRepository(testPathDir)
	require.NoError(t, err)

	r, err := git.InitializeRepo(p)
	require.NoError(t, err)

	return &TestHelper{
		Repository: r,
		RepoPath:   p,
	}
}

// CleanUp removes the directory the test fixture was extracted to.
func (h *TestHelper) CleanUp(t *testing.T) {
	err := os.RemoveAll(filepath.Dir(h.RepoPath))
	require.NoError(t, err)
}

// DirtyRepoPath is the path of the fixture repository with a dirty worktree.
func (h *TestHelper) DirtyRepoPath() string {
	return filepath.Join(h.RepoPath, "dirty-repo")
}

// BasicRepoPath is the path of the fixture repository with a clean worktree.
func (h *TestHelper) BasicRepoPath() string {
	return filepath.Join(h.RepoPath, "basic-repo")
}

// NonRepoPath is the path of a fixture directory that is not a repository.
func (h *TestHelper) NonRepoPath() string {
	return filepath.Join(h.RepoPath, "non-repo")
}
