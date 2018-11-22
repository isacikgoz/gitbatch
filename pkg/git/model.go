package git

import (
	"github.com/fatih/color"
	"github.com/isacikgoz/gitbatch/pkg/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
)

func (entity *RepoEntity) GetRemotes() (remotes []string, err error) {
	r := entity.Repository
    red := color.New(color.FgRed)
    if list, err := r.Remotes(); err != nil {
        return remotes, err
    } else {
        for _, r := range list {
        	remoteString := r.Config().Name + string(red.Sprint(" → ")) + r.Config().URLs[0]
            remotes = append(remotes, remoteString)
        }
    }
    return remotes, nil
}

func getRemotes(r *git.Repository) (remotes []string, err error) {

    if list, err := r.Remotes(); err != nil {
        return remotes, err
    } else {
        for _, r := range list {
        	remoteString := r.Config().Name
            remotes = append(remotes, remoteString)
        }
    }
    return remotes, nil
}

func (entity *RepoEntity) GetCommits() (commits []string, err error) {
	r := entity.Repository
	//TODO: Handle Errors
	ref, _ := r.Head()

    cIter, _ := r.Log(&git.LogOptions{
    	From: ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})

// ... just iterates over the commits
    err = cIter.ForEach(func(c *object.Commit) error {
    	commitstring := utils.ColoredString(string([]rune(c.Hash.String())[:7]), color.FgGreen) + " " + c.Message
    	re := regexp.MustCompile(`\r?\n`)
    	commitstring = re.ReplaceAllString(commitstring, " ")
        commits = append(commits, commitstring)

        return nil
	})
	if err != nil {
		return commits, err
	}
    return commits, nil
}

func (entity *RepoEntity) GetStatus() (status string) {
	status = "↑ " + entity.Pushables + " ↓ " + entity.Pullables + " → " + entity.Branch
	re := regexp.MustCompile(`\r?\n`)
    status = re.ReplaceAllString(status, " ")
    return status
}

func (entity *RepoEntity) GetDisplayString() string{

    blue := color.New(color.FgBlue)
    green := color.New(color.FgGreen)
    red := color.New(color.FgRed)
    cyan := color.New(color.FgCyan)
    orange := color.New(color.FgYellow)
    white := color.New(color.FgWhite)

    prefix := string(blue.Sprint("↑")) + " " + entity.Pushables + " " +
     string(blue.Sprint("↓")) + " " + entity.Pullables + string(red.Sprint(" → ")) + string(cyan.Sprint(entity.Branch)) + " "

	if entity.Marked {
		return prefix + string(green.Sprint(entity.Name))
	} else if !entity.Clean {
		return prefix + string(orange.Sprint(entity.Name))
	} else {
		return prefix + string(white.Sprint(entity.Name))
	}
}

func  (entity *RepoEntity) GetBranches() (branches []string, err error) {
    localBranches, err := entity.LocalBranches()
    red := color.New(color.FgRed)
    if err != nil {
        return nil, err
    }
    for _, b := range localBranches {
        prefix := "   "
        if b == entity.GetActiveBranch() {
            prefix = " → "
        }
        branches = append(branches, (string(red.Sprint(prefix)) + b))
    }
    return branches, nil
}