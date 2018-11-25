package git

import (
	"github.com/fatih/color"
	"regexp"
)

func (entity *RepoEntity) GetRemotes() (remotes []string, err error) {

	r := entity.Repository
    // green := color.New(color.FgGreen)
    if list, err := remoteBranches(&r); err != nil {
        return remotes, err
    } else {
        for _, r := range list {
            remotes = append(remotes, r)
        }
    }

    return remotes, nil
}

func (entity *RepoEntity) GetStatus() (status string) {
	status = "↑ " + entity.Pushables + " ↓ " + entity.Pullables + " → " + entity.Branch
	re := regexp.MustCompile(`\r?\n`)
    status = re.ReplaceAllString(status, " ")
    return status
}

func (entity *RepoEntity) DisplayString() string{

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