package git

import "testing"

func TestRandomString(t *testing.T) {
	stringLength := 8
	randString := RandomString(stringLength)
	if len(randString) != stringLength {
		t.Errorf("The length of the string should be equal.")
	}
}
