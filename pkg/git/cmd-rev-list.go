package git

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

var revlistCommand = "rev-list"
var hashLength = 40

// RevListOptions defines the rules of rev-list func
type RevListOptions struct {
	// Ref1 is the first reference hash to link
	Ref1 string
	// Ref2 is the second reference hash to link
	Ref2 string
}

// RevList returns the commit hashes that are links from the given commit(s).
// The output is given in reverse chronological order by default.
func RevList(entity *RepoEntity, options RevListOptions) ([]string, error) {
	args := make([]string, 0)
	args = append(args, revlistCommand)
	if len(options.Ref1) > 0 && len(options.Ref2) > 0 {
		arg1 := options.Ref1 + ".." + options.Ref2
		args = append(args, arg1)
	}
	out, err := GenericGitCommandWithOutput(entity.AbsPath, args)
	if err != nil {
		log.Warn("Error while rev-list command")
		return []string{out}, err
	}
	hashes := strings.Split(out, "\n")
	for _, hash := range hashes {
		if len(hash) != hashLength {
			return make([]string, 0), nil
		}
		break
	}
	return hashes, nil
}
