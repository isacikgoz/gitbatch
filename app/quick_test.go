package app

import (
	"testing"
)

func TestQuick(t *testing.T) {
	var tests = []struct {
		inp1 []string
		inp2 string
	}{
		{
			[]string{dirty},
			"fetch",
		}, {
			[]string{dirty},
			"pull",
		},
	}
	for _, test := range tests {
		quick(test.inp1, test.inp2)
	}
}
