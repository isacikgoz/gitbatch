package main

import (
	"os"
	"strings"
	"testing"
)

var (
	wd, _ = os.Getwd()
	sp    = string(os.PathSeparator)
	d     = strings.TrimSuffix(wd, sp+"app")

	relparent = ".." + sp + "test"
	parent    = d + sp + "test"
	data      = parent + sp + "test-data"
)

func TestRun(t *testing.T) {
	var tests = []struct {
		inp1 []string
		inp2 string
		inp3 int
		inp4 bool
		inp5 string
	}{
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
