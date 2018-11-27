package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) updateBranch(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	out, err := g.View("branch")
	if err != nil {
		return err
	}
	out.Clear()

	currentindex := 0
	totalbranches := len(entity.Branches)
	for i, b := range entity.Branches {
		var suffix string
		if !b.Clean {
			suffix = " " + yellow.Sprint("✗")
		} else {
			suffix = " " + green.Sprint("✔")
		}
		if b.Name == entity.Branch.Name {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator()+b.Name+suffix)
			continue
		}
		fmt.Fprintln(out, tab()+b.Name)
	}
	if err = gui.smartAnchorRelativeToLine(out, currentindex, totalbranches); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) nextBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity, err := gui.getSelectedRepository(g, v)
	if err != nil {
		return err
	}
	if err = entity.Checkout(entity.NextBranch()); err != nil {
		if err = gui.openErrorView(g, "Stage your changes before checkout", "You should manually manage this issue"); err != nil {
			return err
		}
		return nil
	}
	if err = gui.updateBranch(g, entity); err != nil {
		return err
	}
	if err = gui.updateCommits(g, entity); err != nil {
		return err
	}
	if err = gui.refreshMain(g); err != nil {
		return err
	}
	return nil
}
