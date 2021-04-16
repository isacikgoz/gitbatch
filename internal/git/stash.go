package git

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// StashedItem holds the required fields for a stashed change
type StashedItem struct {
	StashID     int
	BranchName  string
	Hash        string
	Description string
	EntityPath  string
}

func (r *Repository) loadStashedItems() error {
	r.Stasheds = make([]*StashedItem, 0)
	output := stashGet(r, "list")
	stashIDRegex := regexp.MustCompile(`stash@{[\d]+}:`)
	stashIDRegexInt := regexp.MustCompile(`[\d]+`)
	stashBranchRegex := regexp.MustCompile(`^(.*?): `)
	stashMsgRegex := regexp.MustCompile(`WIP on \(?([^)]*)\)?`)
	stashHashRegex := regexp.MustCompile(`[\w|\d]{7}\s`)

	stashlist := strings.Split(output, "\n")
	for _, stashitem := range stashlist {
		// find id
		id := stashIDRegexInt.FindString(stashIDRegex.FindString(stashitem))
		i, err := strconv.Atoi(id)
		if err != nil {
			// probably something isn't right let's continue over this iteration
			continue
		}
		// trim id section
		trimmed := stashIDRegex.Split(stashitem, 2)[1]

		// let's ignore autostash
		if trimmed[1:] == "autostash" {
			continue
		}
		// find branch
		stashBranchRegexMatch := stashBranchRegex.FindString(trimmed)
		branchName := stashBranchRegexMatch[:len(stashBranchRegexMatch)-2]

		branchMatches := stashMsgRegex.FindStringSubmatch(branchName)
		if len(branchMatches) >= 2 {
			branchName = stashBranchRegexMatch[:len(stashBranchRegexMatch)-2]
		}

		// trim branch section
		trimmed = stashBranchRegex.Split(trimmed, 2)[1]
		hash := ""

		var desc string
		if stashHashRegex.MatchString(trimmed) {
			hash = stashHashRegex.FindString(trimmed)[:7]
			desc = stashHashRegex.Split(trimmed, 2)[1]
		} else {
			desc = trimmed
		}
		// trim hash

		r.Stasheds = append(r.Stasheds, &StashedItem{
			StashID:     i,
			BranchName:  branchName,
			Hash:        hash,
			Description: desc,
			EntityPath:  r.AbsPath,
		})
	}
	return nil
}

func stashGet(r *Repository, option string) string {
	args := make([]string, 0)
	args = append(args, "stash")
	args = append(args, option)
	cmd := exec.Command("git", args...)
	cmd.Dir = r.AbsPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "?"
	}
	return string(output)
}

// Pop is the wrapper of "git stash pop" command that used for a file
func (stashedItem *StashedItem) Pop() (string, error) {
	args := make([]string, 0)
	args = append(args, "stash")
	args = append(args, "pop")
	args = append(args, "stash@{"+strconv.Itoa(stashedItem.StashID)+"}")
	cmd := exec.Command("git", args...)
	cmd.Dir = stashedItem.EntityPath
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Show is the wrapper of "git stash show -p " command
func (stashedItem *StashedItem) Show() (string, error) {
	args := make([]string, 0)
	args = append(args, "stash")
	args = append(args, "show")
	args = append(args, "-p")
	args = append(args, "stash@{"+strconv.Itoa(stashedItem.StashID)+"}")
	cmd := exec.Command("git", args...)
	cmd.Dir = stashedItem.EntityPath
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// Stash is the wrapper of conventional "git stash" command
func (r *Repository) Stash() (string, error) {
	args := make([]string, 0)
	args = append(args, "stash")

	cmd := exec.Command("git", args...)
	cmd.Dir = r.AbsPath
	output, err := cmd.CombinedOutput()
	_ = r.Refresh()
	return string(output), err
}
