package git

import (
	"math/rand"
	"time"
)

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString generates a random string of n length
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[r.Intn(len(characterRunes))]
	}
	return string(b)
}
