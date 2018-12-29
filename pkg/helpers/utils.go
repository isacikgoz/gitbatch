package helpers

import (
	"math/rand"
	"strings"
	"time"
)

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// TrimTrailingNewline removes the trailing new line form a string. this method
// is used mostly on outputs of a command
func TrimTrailingNewline(str string) string {
	if strings.HasSuffix(str, "\n") {
		return str[:len(str)-1]
	}
	return str
}

// Min finds the minimum value of two int
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// RandomString generates a random string of n length
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[r.Intn(len(characterRunes))]
	}
	return string(b)
}
