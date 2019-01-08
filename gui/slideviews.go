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

	gui.renderCommits(r)

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

func (gui *Gui) selectCommit(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	r := gui.getSelectedRepository()
	ix := oy + cy

	r.State.Branch.State.Commit = r.State.Branch.Commits[ix]
	return gui.renderCommits(r)
}

func (gui *Gui) commitDetail(ix int) error {
	v, err := gui.g.View(detailViewFeature.Name)
	if err != nil {
		return err
	}

	r := gui.getSelectedRepository()
	v.Clear()
	c := r.State.Branch.Commits[ix]

	fmt.Fprintf(v, "%s\n", c.String())
	fmt.Fprintf(v, decorateDiffStat(c.DiffStat()))
	return nil
}

func (gui *Gui) commitDiff(g *gocui.Gui, v *gocui.View) error {
	v, err := gui.g.View(detailViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	vcm, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	_, oy := vcm.Origin()
	_, cy := vcm.Cursor()
	ix := oy + cy
	r := gui.getSelectedRepository()
	c := r.State.Branch.Commits[ix]
	if ix+1 > len(r.State.Branch.Commits) {
		ix = ix - 1
	}
	p, err := r.State.Branch.Commits[ix+1].C.Patch(c.C)
	var s string
	for _, d := range colorizeDiff(p.String()) {
		s = s + "\n" + d
	}
	fmt.Fprintf(v, s)
	return nil
}

// moves cursor down for a page size
func (gui *Gui) dpageDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, _ := v.Cursor()
		_, vy := v.Size()
		lr := len(v.BufferLines())
		if lr < vy {
			return nil
		}
		if oy+vy >= lr-vy {
			if err := v.SetOrigin(ox, lr-vy); err != nil {
				return err
			}
		} else if err := v.SetOrigin(ox, oy+vy); err != nil {
			return err
		}
		if err := v.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	return nil
}

// moves cursor up for a page size
func (gui *Gui) dpageUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		_, vy := v.Size()
		if oy == 0 || oy+cy < vy {
			if err := v.SetOrigin(ox, 0); err != nil {
				return err
			}
		} else if oy <= vy {
			if err := v.SetOrigin(ox, oy+cy-vy); err != nil {
				return err
			}
		} else if err := v.SetOrigin(ox, oy-vy); err != nil {
			return err
		}
		if err := v.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	return nil
}
