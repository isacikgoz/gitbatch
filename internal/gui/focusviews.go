package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/jroimartin/gocui"
)

// listens the event -> "branch.updated"
func (gui *Gui) branchUpdated(event *git.RepositoryEvent) error {
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.renderCommits(gui.getSelectedRepository())
	})
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
		_ = adjustAnchor(pos, ly, v)
		if err := gui.commitStats(pos); err != nil {
			return err
		}
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
		_ = adjustAnchor(pos, ly, v)
		if pos >= 0 {
			if err := gui.commitStats(cy + oy - 1); err != nil {
				return err
			}
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
	// bc := r.State.Branch.State.Commit
	si := 0
	fmt.Fprintln(v, " "+yellow.Sprint("*******")+" "+yellow.Sprint("Current State"))

	for _, c := range cs {
		// if c.Hash == bc.Hash {
		// 	si = i
		// 	fmt.Fprintln(v, ws+commitLabel(c, false))
		// 	continue
		// }

		fmt.Fprintln(v, tab+commitLabel(c, false))
	}
	_ = adjustAnchor(si, len(cs), v)
	return nil
}

// moves cursor down for a page size
func (gui *Gui) commitPageDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		_, vy := v.Size()
		_, cy := v.Cursor()
		lr := len(v.BufferLines())
		if lr < vy {
			return nil
		}
		v.EditDelete(true)

		_ = adjustAnchor(oy+cy+vy-1, lr, v)
		// if err := gui.commitStats(oy + cy + vy - 1); err != nil {
		// 	return err
		// }

	}
	return nil
}

// moves cursor to the top
func (gui *Gui) commitCursorTop(g *gocui.Gui, v *gocui.View) error {
	if v != nil {

		v.EditDelete(true)
		lr := len(v.BufferLines())

		_ = adjustAnchor(0, lr, v)
		if err := gui.commitStats(0); err != nil {
			return err
		}
	}
	return nil
}

// moves cursor up for a page size
func (gui *Gui) commitPageUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		_, cy := v.Cursor()
		_, vy := v.Size()
		lr := len(v.BufferLines())
		v.EditDelete(true)
		_ = adjustAnchor(oy+cy-vy+1, lr, v)
		// if err := gui.commitStats(oy + cy - vy + 1); err != nil {
		// 	return err
		// }
	}
	return nil
}
