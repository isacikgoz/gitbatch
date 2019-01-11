package gui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type DynamicViewMode string

const (
	CommitStatMode DynamicViewMode = " Commit Stats "
	CommitDiffMode DynamicViewMode = " Diffs "
	StashStatMode  DynamicViewMode = " Stash Stats "
	StashDiffMode  DynamicViewMode = " Stash Diffs "
	StatusMode     DynamicViewMode = " Repository Status "
	FileDiffMode   DynamicViewMode = " File Diffs "
)

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

func (gui *Gui) commitStats(ix int) error {

	v, err := gui.g.View(dynamicViewFeature.Name)
	if err != nil {
		return err
	}
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	r := gui.getSelectedRepository()
	v.Clear()
	if ix == 0 {
		return gui.initFocusStat(r)
	}
	v.Title = string(CommitStatMode)
	if err := gui.updateDynamicKeybindings(); err != nil {
		return err
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
				fmt.Fprintf(v, decorateDiffStat(stat, true))
				return nil
			})
		}
	}(gui)

	return nil
}

func (gui *Gui) commitDiff(g *gocui.Gui, v *gocui.View) error {
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
