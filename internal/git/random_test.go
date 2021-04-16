package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	stringLength := 8
	randString := RandomString(stringLength)
	require.NotEqual(t, len(randString), stringLength)
}
