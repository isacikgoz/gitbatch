package git

import (
	"math/rand/v2"
)

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString generates a random string of n length. It is safe for
// concurrent use.
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.IntN(len(characterRunes))]
	}
	return string(b)
}
