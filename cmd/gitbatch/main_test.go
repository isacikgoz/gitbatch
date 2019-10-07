package main

import (
	"os"
	"testing"
)

var (
	sp = string(os.PathSeparator)
)

func TestRun(t *testing.T) {
	var tests = []struct {
		inp1 []string
		inp2 string
		inp3 int
		inp4 bool
		inp5 string
	}{
		// TODO: generate scenarios
		{
			[]string{},
			"debug",
			0,
			true,
			"fetch",
		},
	}
	for _, test := range tests {
		if err := run(test.inp1, test.inp2, test.inp3, test.inp4, test.inp5); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}
