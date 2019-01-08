package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/core/git"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) focusToRepository(g *gocui.Gui, v *gocui.View) error {

	r := gui.getSelectedRepository()
	gui.mode = true

	if _, err := g.SetCurrentView(commitViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, commitViewFeature.Name)

	r.State.Branch.InitializeCommits(r)

	gui.renderCommits(r)

	gui.g.Update(func(g *gocui.Gui) error {
		return gui.renderMain()
	})
	return nil
}

func (gui *Gui) focusBackToMain(g *gocui.Gui, v *gocui.View) error {

	gui.mode = false

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
		adjustAnchor(cy+oy+1, ly, v)
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
		adjustAnchor(cy+oy-1, ly, v)
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
	for i, c := range cs {
		if c.Hash == bc.Hash {
			si = i
			fmt.Fprintln(v, ws+commitLabel(c, true))
			continue
		}

		fmt.Fprintln(v, tab+commitLabel(c, false))
	}
	adjustAnchor(si, len(cs), v)
	return nil
}
