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

    remotes, err = getRemotes(&r)
    if err !=nil {
    	return nil ,err
    }
    return remotes, nil
}

func getRemotes(r *git.Repository) (remotes []string, err error) {

    if list, err := r.Remotes(); err != nil {
        return remotes, err
    } else {
        for _, r := range list {
        	remoteString := r.Config().Name + " → " + r.Config().URLs[0]
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
	if entity.Marked {
		green := color.New(color.FgGreen)
		return string(green.Sprint(entity.Name))
	} else if !entity.Clean {
		orange := color.New(color.FgYellow)
		return string(orange.Sprint(entity.Name))
	} else {
		white := color.New(color.FgWhite)
		return string(white.Sprint(entity.Name))
	}
}