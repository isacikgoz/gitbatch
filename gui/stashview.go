package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/core/git"
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
	r := gui.getSelectedRepository()
	err = refreshStashView(g, r)
	return err
}

//
func (gui *Gui) stashChanges(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	output, err := r.Stash()
	if err != nil {
		if err = gui.openErrorView(g, output,
			"You should manually resolve this issue",
			stashViewFeature.Name); err != nil {
			return err
		}
	}
	err = refreshAllStatusView(g, r, true)
	return err
}

//
func (gui *Gui) popStash(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	_, oy := v.Origin()
	_, cy := v.Cursor()
	if len(r.Stasheds) <= 0 {
		return nil
	}
	stashedItem := r.Stasheds[oy+cy]
	output, err := stashedItem.Pop()
	if err != nil {
		if err = gui.openErrorView(g, output,
			"You should manually resolve this issue",
			stashViewFeature.Name); err != nil {
			return err
		}
	}
	// since the pop is a func of stashed item, we need to refresh entity here
	if err := r.Refresh(); err != nil {
		return err
	}

	return refreshAllStatusView(g, r, true)
}

// refresh the main view and re-render the repository representations
func refreshStashView(g *gocui.Gui, r *git.Repository) error {
	stashView, err := g.View(stashViewFeature.Name)
	if err != nil {
		return err
	}
	stashView.Clear()
	_, cy := stashView.Cursor()
	_, oy := stashView.Origin()
	stashedItems := r.Stasheds
	for i, stashedItem := range stashedItems {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stashView, "%s%d %s: %s (%s)\n", prefix, stashedItem.StashID, cyan.Sprint(stashedItem.BranchName), stashedItem.Description, cyan.Sprint(stashedItem.Hash))
	}
	return nil
}
