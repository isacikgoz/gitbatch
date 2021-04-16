package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestConfigWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	testConfigopt := &ConfigOptions{
		Section: "remote.origin",
		Option:  "url",
		Site:    ConfigSiteLocal,
	}

	var tests = []struct {
		inp1     *git.Repository
		inp2     *ConfigOptions
		expected string
	}{
		{th.Repository, testConfigopt, "https://gitlab.com/isacikgoz/test-data.git"},
	}
	for _, test := range tests {
		output, err := configWithGit(test.inp1, test.inp2)
		require.NoError(t, err)
		require.Equal(t, test.expected, output)
	}
}

func TestConfigWithGoGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	testConfigopt := &ConfigOptions{
		Section: "core",
		Option:  "bare",
		Site:    ConfigSiteLocal,
	}

	var tests = []struct {
		inp1     *git.Repository
		inp2     *ConfigOptions
		expected string
	}{
		{th.Repository, testConfigopt, "false"},
	}
	for _, test := range tests {
		output, err := configWithGoGit(test.inp1, test.inp2)
		require.NoError(t, err)
		require.Equal(t, output, test.expected)
	}
}

func TestAddConfigWithGit(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	testConfigopt := &ConfigOptions{
		Section: "user",
		Option:  "name",
		Site:    ConfigSiteLocal,
	}

	var tests = []struct {
		inp1 *git.Repository
		inp2 *ConfigOptions
		inp3 string
	}{
		{th.Repository, testConfigopt, "foo"},
	}
	for _, test := range tests {
		err := addConfigWithGit(test.inp1, test.inp2, test.inp3)
		require.NoError(t, err)
	}
}
