package git

import (
	"unicode"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// Alphabetical slice is the re-ordered *Repository slice that sorted according
// to alphabetical order (A-Z)
type Alphabetical []*Repository

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

// LastModified slice is the re-ordered *Repository slice that sorted according
// to last modified date of the repository directory
type LastModified []*Repository

// Len is the interface implementation for LastModified sorting function
func (s LastModified) Len() int { return len(s) }

// Swap is the interface implementation for LastModified sorting function
func (s LastModified) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less is the interface implementation for LastModified sorting function
func (s LastModified) Less(i, j int) bool {
	return s[i].ModTime.Unix() > s[j].ModTime.Unix()
}

// Less returns a comparison between to repositories by name
func Less(ri, rj *Repository) bool {
	iRunes := []rune(ri.Name)
	jRunes := []rune(rj.Name)

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

// CommitTime slice is the re-ordered *object.Commit slice that sorted according
// commit date
type CommitTime []*object.Commit

// Len is the interface implementation for LastModified sorting function
func (s CommitTime) Len() int { return len(s) }

// Swap is the interface implementation for LastModified sorting function
func (s CommitTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less is the interface implementation for LastModified sorting function
func (s CommitTime) Less(i, j int) bool {
	return s[i].Author.When.Unix() > s[j].Author.When.Unix()
}
