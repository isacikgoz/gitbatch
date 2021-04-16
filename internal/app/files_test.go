package app

import (
	"path/filepath"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestGenerateDirectories(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1     []string
		inp2     int
		expected []string
	}{
		{[]string{th.RepoPath}, 1, []string{th.BasicRepoPath(), th.DirtyRepoPath()}},
		{[]string{th.RepoPath}, 2, []string{th.BasicRepoPath(), th.DirtyRepoPath()}}, // maybe move one repo to a sub folder
	}
	for _, test := range tests {
		output := generateDirectories(test.inp1, test.inp2)
		require.ElementsMatch(t, output, test.expected)
	}
}

func TestWalkRecursive(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1 []string
		inp2 []string
		exp1 []string
		exp2 []string
	}{
		{
			[]string{th.RepoPath},
			[]string{""},
			[]string{filepath.Join(th.RepoPath, ".git"), filepath.Join(th.RepoPath, ".gitmodules"), th.NonRepoPath()},
			[]string{"", th.BasicRepoPath(), th.DirtyRepoPath()},
		},
	}
	for _, test := range tests {
		out1, out2 := walkRecursive(test.inp1, test.inp2)
		require.ElementsMatch(t, out1, test.exp1)
		require.ElementsMatch(t, out2, test.exp2)
	}
}

func TestSeparateDirectories(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		input string
		exp1  []string
		exp2  []string
	}{
		{
			"",
			nil,
			nil,
		},
		{
			th.RepoPath,
			[]string{filepath.Join(th.RepoPath, ".git"), filepath.Join(th.RepoPath, ".gitmodules"), th.NonRepoPath()},
			[]string{th.BasicRepoPath(), th.DirtyRepoPath()},
		},
	}
	for _, test := range tests {
		out1, out2, err := separateDirectories(test.input)
		require.NoError(t, err)
		require.ElementsMatch(t, out1, test.exp1)
		require.ElementsMatch(t, out2, test.exp2)
	}
}
