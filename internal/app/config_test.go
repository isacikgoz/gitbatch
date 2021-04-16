package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfiguration(t *testing.T) {
	cfg, err := loadConfiguration()
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestReadConfiguration(t *testing.T) {
	err := readConfiguration()
	require.NoError(t, err)
}

func TestInitializeConfigurationManager(t *testing.T) {
	err := initializeConfigurationManager()
	require.NoError(t, err)
}

func TestOsConfigDirectory(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"linux", ".config"},
		{"darwin", "Application Support"},
	}
	for _, test := range tests {
		output := osConfigDirectory(test.input)
		require.Contains(t, output, test.expected)
	}
}
