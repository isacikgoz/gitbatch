package app

import (
	"os"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestOverrideConfig(t *testing.T) {
	config1 := &Config{
		Directories: []string{},
		LogLevel:    "info",
		Depth:       1,
		QuickMode:   false,
		Mode:        "fetch",
	}
	config2 := &Config{
		Directories: []string{string(os.PathSeparator) + "tmp"},
		LogLevel:    "error",
		Depth:       1,
		QuickMode:   true,
		Mode:        "pull",
	}

	var tests = []struct {
		inp1     *Config
		inp2     *Config
		expected *Config
	}{
		{config1, config2, config1},
	}
	for _, test := range tests {
		output := overrideConfig(test.inp1, test.inp2)
		require.Equal(t, test.expected, output)
		require.Equal(t, test.inp2.Mode, output.Mode)
	}
}

func TestExecQuickMode(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	var tests = []struct {
		inp1 []string
	}{
		{[]string{th.BasicRepoPath()}},
	}
	a := App{
		Config: &Config{
			Mode: "fetch",
		},
	}
	for _, test := range tests {
		err := a.execQuickMode(test.inp1)
		require.NoError(t, err)
	}
}
