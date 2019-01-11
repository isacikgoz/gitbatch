package gui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type DetailViewMode string

const (
	CommitStatMode DetailViewMode = " Commit Stats "
	CommitDiffMode DetailViewMode = " Diffs "
	StashStatMode  DetailViewMode = " Stash Stats "
	StashDiffMode  DetailViewMode = " Stash Diffs "
	StatusMode     DetailViewMode = " Repository Status "
	FileDiffMode   DetailViewMode = " File Diffs "
)

var (
	detailViewModes = []DetailViewMode{CommitStatMode, CommitDiffMode, StashDiffMode, StashStatMode, StatusMode, FileDiffMode}
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

	v, err := gui.g.View(detailViewFeature.Name)
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
	if err := gui.removeDetailViewKeybindings(); err != nil {
		return err
	}
	v.Title = string(CommitStatMode)
	if err := gui.updateDiffViewKeybindings(); err != nil {
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
			if !strings.Contains(strings.Join(v.BufferLines(), c.Hash), ld) {
				return
			}
			v.Clear()
			fmt.Fprintf(v, "%s\n", decorateCommit(c.String()))
			gui.g.Update(func(g *gocui.Gui) error {
				v, err := gui.g.View(detailViewFeature.Name)
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
	v, err := gui.g.View(detailViewFeature.Name)
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

	if err := gui.removeDetailViewKeybindings(); err != nil {
		return err
	}
	v.Title = string(CommitDiffMode)
	if err := gui.updateDiffViewKeybindings(); err != nil {
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

func (gui *Gui) updateDiffViewKeybindings() error {
	v, err := gui.g.View(detailViewFeature.Name)
	if err != nil {
		return err
	}

	if err := gui.generateKeybindingsForDetailView(v.Title); err != nil {
		return err
	}
	if err := gui.updateKeyBindingsView(gui.g, v.Title); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) removeDetailViewKeybindings() error {
	a := gui.KeyBindings
	gui.g.DeleteKeybindings(detailViewFeature.Name)
	// for i, b := range a {
	// 	if b.View == detailViewFeature.Name {
	// 		a = a[:i+copy(a[i:], a[i+1:])]
	// 	}
	// }
	gui.KeyBindings = a
	return nil
}
