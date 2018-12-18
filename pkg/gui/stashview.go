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
	}
	e := gui.getSelectedRepository()
	err = refreshStashView(g, e)
	return err
}

//
func (gui *Gui) stashChanges(g *gocui.Gui, v *gocui.View) error {
	e := gui.getSelectedRepository()
	output, err := e.Stash()
	if err != nil {
		if err = gui.openErrorView(g, output,
			"You should manually resolve this issue",
			stashViewFeature.Name); err != nil {
			return err
		}
	}
	err = refreshAllStatusView(g, e, true)
	return err
}

//
func (gui *Gui) popStash(g *gocui.Gui, v *gocui.View) error {
	e := gui.getSelectedRepository()
	_, oy := v.Origin()
	_, cy := v.Cursor()
	if len(e.Stasheds) <= 0 {
		return nil
	}
	stashedItem := e.Stasheds[oy+cy]
	output, err := stashedItem.Pop()
	if err != nil {
		if err = gui.openErrorView(g, output,
			"You should manually resolve this issue",
			stashViewFeature.Name); err != nil {
			return err
		}
	}
	// since the pop is a func of stashed item, we need to refresh entity here
	e.Refresh()

	err = refreshAllStatusView(g, e, true)
	return err
}

// refresh the main view and re-render the repository representations
func refreshStashView(g *gocui.Gui, e *git.RepoEntity) error {
	stashView, err := g.View(stashViewFeature.Name)
	if err != nil {
		return err
	}
	stashView.Clear()
	_, cy := stashView.Cursor()
	_, oy := stashView.Origin()
	stashedItems := e.Stasheds
	for i, stashedItem := range stashedItems {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stashView, "%s%d %s: %s (%s)\n", prefix, stashedItem.StashID, cyan.Sprint(stashedItem.BranchName), stashedItem.Description, cyan.Sprint(stashedItem.Hash))
	}
	return nil
}
