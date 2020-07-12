package app

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/isacikgoz/gitbatch/internal/git"
	ggit "gopkg.in/src-d/go-git.v4"
)

var (
	config1 = &Config{
		Directories: []string{},
		LogLevel:    "info",
		Depth:       1,
		QuickMode:   false,
		Mode:        "fetch",
	}
	config2 = &Config{
		Directories: []string{string(os.PathSeparator) + "tmp"},
		LogLevel:    "error",
		Depth:       1,
		QuickMode:   true,
		Mode:        "pull",
	}

	testRepoDir, _ = ioutil.TempDir("", "test-data")
)

func TestSetup(t *testing.T) {
	mockApp := &App{Config: config1}
	var tests = []struct {
		input    *Config
		expected *App
	}{
		{config2, nil},
		{config1, mockApp},
	}
	for _, test := range tests {

		app, err := New(test.input)
		if err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
		q := test.input.QuickMode
		if q && app != nil {
			t.Errorf("Test Failed.")
		} else if !q && app == nil {
			t.Errorf("Test Failed.")
		}

	}
}

func TestOverrideConfig(t *testing.T) {
	var tests = []struct {
		inp1     *Config
		inp2     *Config
		expected *Config
	}{
		{config1, config2, config1},
	}
	for _, test := range tests {
		if output := overrideConfig(test.inp1, test.inp2); output != test.expected || test.inp2.Mode != output.Mode {
			t.Errorf("Test Failed: {%s, %s} inputted, output: %s, expected: %s", test.inp1.Directories, test.inp2.Directories, output.Directories, test.expected.Directories)
		}
	}
}

func TestExecQuickMode(t *testing.T) {
	defer cleanRepo()
	_, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	var tests = []struct {
		inp1 []string
	}{
		{[]string{basic}},
	}
	a := App{
		Config: &Config{
			Mode: "fetch",
		},
	}
	for _, test := range tests {
		if err := a.execQuickMode(test.inp1); err != nil {
			t.Errorf("Test Failed: %s", err.Error())
		}
	}
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

func testFile(name string) (*git.File, error) {
	_, err := os.Create(testRepoDir + string(os.PathSeparator) + name)
	if err != nil {
		return nil, err
	}
	f := &git.File{
		Name:    name,
		AbsPath: testRepoDir + string(os.PathSeparator) + name,
		X:       git.StatusUntracked,
		Y:       git.StatusUntracked,
	}
	return f, nil
}

func cleanRepo() error {
	return os.RemoveAll(testRepoDir)
}
