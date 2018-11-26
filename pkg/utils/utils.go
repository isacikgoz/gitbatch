package utils

import (
	"strings"
)

func TrimTrailingNewline(str string) string {
	if strings.HasSuffix(str, "\n") {
		return str[:len(str)-1]
	}
	return str
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}