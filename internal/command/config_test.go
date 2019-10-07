package command

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
)

var (
	testConfigopt1 = &ConfigOptions{
		Section: "remote.origin",
		Option:  "url",
		Site:    ConfigSiteLocal,
	}
	testConfigopt2 = &ConfigOptions{
		Section: "core",
		Option:  "bare",
		Site:    ConfigSiteLocal,
	}
	testConfigopt3 = &ConfigOptions{
		Section: "user",
		Option:  "name",
		Site:    ConfigSiteLocal,
	}
)

func TestConfigWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1     *git.Repository
		inp2     *ConfigOptions
		expected string
	}{
		{r, testConfigopt1, "https://gitlab.com/isacikgoz/dirty-repo.git"},
	}
	for _, test := range tests {
		if output, err := configWithGit(test.inp1, test.inp2); err != nil || output != test.expected {
			t.Errorf("Test Failed: %s expected, %s was the output.", test.expected, output)
		}
	}
}

func TestConfigWithGoGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1     *git.Repository
		inp2     *ConfigOptions
		expected string
	}{
		{r, testConfigopt2, "false"},
	}
	for _, test := range tests {
		if output, err := configWithGoGit(test.inp1, test.inp2); err != nil || output != test.expected {
			t.Errorf("Test Failed: %s expected, %s was the output.", test.expected, output)
		}
	}
}

func TestAddConfigWithGit(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 *git.Repository
		inp2 *ConfigOptions
		inp3 string
	}{
		{r, testConfigopt3, "foo"},
	}
	for _, test := range tests {
		if err := addConfigWithGit(test.inp1, test.inp2, test.inp3); err != nil {
			t.Errorf("Test Failed: error: %s", err.Error())
		}
	}
}
