package gui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

// DynamicViewMode is a indicator of dynamic view mode and it is also used as
// the views title
type DynamicViewMode string

const (
	// CommitStatMode when dynamic mode morphed into commit stat mode
	CommitStatMode DynamicViewMode = " Commit Stats "
	// CommitDiffMode when dynamic mode morphed into commit diff mode
	CommitDiffMode DynamicViewMode = " Diffs "
	// StashStatMode when dynamic mode morphed into stash status mode
	StashStatMode DynamicViewMode = " Stash Stats "
	// StashDiffMode when dynamic mode morphed into stash diff mode
	StashDiffMode DynamicViewMode = " Stash Diffs "
	// StatusMode when dynamic mode morphed into repository status mode
	StatusMode DynamicViewMode = " Repository Status "
	// FileDiffMode when dynamic mode morphed into file diff mode
	FileDiffMode DynamicViewMode = " File Diffs "
)

// shows the stats of current commit
func (gui *Gui) commitStat(g *gocui.Gui, v *gocui.View) error {
	vc, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	_, oy := vc.Origin()
	_, cy := vc.Cursor()

	return gui.commitStats(oy + cy)

}

// render the commit stats of the give index
func (gui *Gui) commitStats(ix int) error {

	v, err := gui.g.View(dynamicViewFeature.Name)
	if err != nil {
		return err
	}
	_ = v.SetOrigin(0, 0)
	_ = v.SetCursor(0, 0)
	r := gui.getSelectedRepository()
	v.Clear()
	if ix == 0 {
		return gui.initFocusStat(r)
	}
	v.Title = string(CommitStatMode)
	if err := gui.updateDynamicKeybindings(); err != nil {
		return err
	}
	if len(r.State.Branch.Commits) <= 0 {
		return nil
	}
	c := r.State.Branch.Commits[ix-1]
	done := make(chan bool)
	var stat string
	go func() {
		stat = c.DiffStat(done)
	}()
	ld := "loading stats..."
	fmt.Fprintf(v, "%s\n", decorateCommit(c.String()))
	fmt.Fprintf(v, "%s\n", red.Sprint(ld))

	go func(gui *Gui) {
		if <-done {
			if !strings.Contains(strings.Join(v.BufferLines(), "\n"), c.Hash) {
				return
			}
			v.Clear()
			fmt.Fprintf(v, "%s\n", decorateCommit(c.String()))
			gui.g.Update(func(g *gocui.Gui) error {
				v, err := gui.g.View(dynamicViewFeature.Name)
				if err != nil {
					return err
				}
				fmt.Fprintf(v, "%s", decorateDiffStat(stat, true))
				return nil
			})
		}
	}(gui)

	return nil
}

// show the diff of current commit
func (gui *Gui) commitDiff(g *gocui.Gui, _ *gocui.View) error {
	v, err := gui.g.View(dynamicViewFeature.Name)
	if err != nil {
		return err
	}
	vcm, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	_, oy := vcm.Origin()
	_, cy := vcm.Cursor()
	ix := oy + cy
	if ix == 0 {
		return nil
	}

	v.Title = string(CommitDiffMode)
	if err := gui.updateDynamicKeybindings(); err != nil {
		return err
	}
	r := gui.getSelectedRepository()

	v.Clear()
	c := r.State.Branch.Commits[ix-1]
	if ix+1 > len(r.State.Branch.Commits) {
		ix = ix - 1
	}
	p, err := r.State.Branch.Commits[ix].C.Patch(c.C)
	if err != nil {
		return err
	}
	var s string
	for _, d := range colorizeDiff(p.String()) {
		s = s + "\n" + d
	}
	fmt.Fprint(v, s)
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
