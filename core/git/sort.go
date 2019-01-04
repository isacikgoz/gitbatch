package git

import (
	"unicode"
)

// Alphabetical slice is the re-ordered *RepoEntity slice that sorted according
// to alphabetical order (A-Z)
type Alphabetical []*RepoEntity

// Len is the interface implementation for Alphabetical sorting function
func (s Alphabetical) Len() int { return len(s) }

// Swap is the interface implementation for Alphabetical sorting function
func (s Alphabetical) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less is the interface implementation for Alphabetical sorting function
func (s Alphabetical) Less(i, j int) bool {
	iRunes := []rune(s[i].Name)
	jRunes := []rune(s[j].Name)

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		// the lowercase runes are the same, so compare the original
		if ir != jr {
			return ir < jr
		}
	}
	return false
}

// LastModified slice is the re-ordered *RepoEntity slice that sorted according
// to last modified date of the repository directory
type LastModified []*RepoEntity

// Len is the interface implementation for LastModified sorting function
func (s LastModified) Len() int { return len(s) }

// Swap is the interface implementation for LastModified sorting function
func (s LastModified) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less is the interface implementation for LastModified sorting function
func (s LastModified) Less(i, j int) bool {
	return s[i].ModTime.Unix() > s[j].ModTime.Unix()
}
