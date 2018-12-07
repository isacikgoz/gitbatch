package git

import (
	"regexp"
	"strings"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var stashCommand = "stash"

type StashedItem struct {
	StashID int
	BranchName string
	Hash string
	Description string
}

// // StashOption is the option argument for git stash command
// type StashOption string

// var (
// 	StashPop StashOption = "pop"
// 	StashPush StashOption = "push"
// 	StashDrop StashOption = "drop"
// )

// // Stash used when you want to record the current state of the working
// // directory and the index, but want to go back to a clean working directory.
// func Stash(entity *RepoEntity, option StashOption) error {
// 	args := make([]string, 0)
// 	args = append(args, stashCommand)

// 	if err := GenericGitCommand(entity.AbsPath, args); err != nil {
// 		log.Warn("Error while stashing")
// 		return err
// 	}
// 	return nil
// }

func stashGet(entity *RepoEntity, option string) string {
	args := make([]string, 0)
	args = append(args, stashCommand)
	args = append(args, option)
	out, err := GenericGitCommandWithOutput(entity.AbsPath, args)
	if err != nil {
		log.Warn("Error while stash command")
		return "?"
	}
	return out
}

func (entity *RepoEntity) loadStashedItems() error {
	entity.Stasheds = make([]*StashedItem, 0)
	output := stashGet(entity, "list")
	stashIDRegex := regexp.MustCompile(`stash@{[\d]+}:`)
	stashIDRegexInt := regexp.MustCompile(`[\d]+`)
	stashBranchRegex := regexp.MustCompile(`[\w]+: `)
	stashHashRegex := regexp.MustCompile(`[\w]{7}`)
	
	stashlist := strings.Split(output, "\n")
	for _, stashitem := range stashlist {
		// find id
		id := stashIDRegexInt.FindString(stashIDRegex.FindString(stashitem))
		i, err := strconv.Atoi(id) 
		if err != nil {
			// probably something isn't right let's continue over this iteration
			log.Trace("cannot initiate stashed item")
			continue
		}
		// trim id section
		trimmed := stashIDRegex.Split(stashitem, 2)[1]
		
		// find branch
		stashBranchRegexMatch :=stashBranchRegex.FindString(trimmed)
		branchName := stashBranchRegexMatch[:len(stashBranchRegexMatch)-2]
		
		// trim branch section
		trimmed = stashBranchRegex.Split(trimmed, 2)[1]
		hash := stashHashRegex.FindString(trimmed)
		
		// trim hash
		desc := stashHashRegex.Split(trimmed, 2)[1][1:]

		entity.Stasheds = append(entity.Stasheds, &StashedItem{
			StashID: i,
			BranchName: branchName,
			Hash: hash,
			Description: desc,
			})
	}
	return nil
}

func (entity *RepoEntity) Stash() (error) {
	args := make([]string, 0)
	args = append(args, stashCommand)
	if err := GenericGitCommand(entity.AbsPath, args); err != nil {
		log.Warn("Error while stashing")
		return err
	}
	return nil
}