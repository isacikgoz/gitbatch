package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/core/git"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) focusToRepository(g *gocui.Gui, v *gocui.View) error {
	mainViews = focusViews
	r := gui.getSelectedRepository()
	gui.order = focus

	if _, err := g.SetCurrentView(commitViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, commitViewFeature.Name)

	r.State.Branch.InitializeCommits(r)

	if err := gui.renderCommits(r); err != nil {
		return err
	}
	if err := gui.initFocusStat(r); err != nil {
		return err
	}
	if err := gui.initStashedView(r); err != nil {
		return err
	}
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.renderMain()
	})
	return nil
}

func (gui *Gui) focusBackToMain(g *gocui.Gui, v *gocui.View) error {
	mainViews = overviewViews
	gui.order = overview

	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) commitCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		ly := len(v.BufferLines()) - 1

		// if we are at the end we just return
		if cy+oy == ly-1 {
			return nil
		}
		v.EditDelete(true)
		pos := cy + oy + 1
		adjustAnchor(pos, ly, v)
		gui.commitDetail(pos)
	}
	return nil
}

// moves the cursor upwards for the main view
func (gui *Gui) commitCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		_, cy := v.Cursor()
		ly := len(v.BufferLines()) - 1
		v.EditDelete(true)
		pos := cy + oy - 1
		adjustAnchor(pos, ly, v)
		if pos >= 0 {
			gui.commitDetail(cy + oy - 1)
		}
	}
	return nil
}

// updates the commitsview for given entity
func (gui *Gui) renderCommits(r *git.Repository) error {
	v, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	cs := r.State.Branch.Commits
	bc := r.State.Branch.State.Commit
	si := 0
	fmt.Fprintln(v, " "+yellow.Sprint("*******")+" "+yellow.Sprint("Current State"))

	for i, c := range cs {
		if c.Hash == bc.Hash {
			si = i
			fmt.Fprintln(v, ws+commitLabel(c, false))
			continue
		}

		fmt.Fprintln(v, tab+commitLabel(c, false))
	}
	adjustAnchor(si, len(cs), v)
	return nil
}
