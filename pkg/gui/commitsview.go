package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

// updates the commitsview for given entity
func (gui *Gui) updateCommits(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	out, err := g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()

	currentindex := 0
	totalcommits := len(entity.Commits)
	for i, c := range entity.Commits {
		var body string
		if c.CommitType == git.LocalCommit {
			body = cyan.Sprint(c.Hash[:hashLength]) + " " + c.Message
		} else {
			body = yellow.Sprint(c.Hash[:hashLength]) + " " + c.Message
		}
		if c.Hash == entity.Commit.Hash {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator+body)
			continue
		}
		fmt.Fprintln(out, tab+body)
	}
	if err = gui.smartAnchorRelativeToLine(out, currentindex, totalcommits); err != nil {
		return err
	}
	return err
}

// iteration handler for the commitsview
func (gui *Gui) nextCommit(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = entity.NextCommit(); err != nil {
		return err
	}
	if err = gui.updateCommits(g, entity); err != nil {
		return err
	}
	return err
}
