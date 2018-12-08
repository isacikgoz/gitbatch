package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

// stash view
func (gui *Gui) openStashView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(stashViewFeature.Name, 6, int(0.75*float32(maxY)), maxX-6, maxY-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stashViewFeature.Title
		v.Wrap = true
	}
	entity := gui.getSelectedRepository()
	if err := refreshStashView(g, entity); err != nil {
		return err
	}
	return nil
}

// refresh the main view and re-render the repository representations
func refreshStashView(g *gocui.Gui, entity *git.RepoEntity) error {
	stashView, err := g.View(stashViewFeature.Name)
	if err != nil {
		return err
	}
	stashView.Clear()
	_, cy := stashView.Cursor()
	_, oy := stashView.Origin()
	stashedItems := entity.Stasheds
	for i, stashedItem := range stashedItems {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stashView, "%s%d %s: %s (%s)\n", prefix, stashedItem.StashID, cyan.Sprint(stashedItem.BranchName), stashedItem.Description, stashedItem.Hash)
	}
	return nil
}