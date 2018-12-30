package helpers

import "testing"

func TestRandomString(t *testing.T) {
	var rand_string = RandomString(8)
	if len(rand_string) != 8 {
		t.Errorf("The length of the string should be equal.")
	}
}
