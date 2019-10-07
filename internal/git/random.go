package git

import (
	"math/rand"
	"time"
)

// RandomString generates a random string of n length
func RandomString(n int) string {
	var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[r.Intn(len(characterRunes))]
	}
	return string(b)
}
