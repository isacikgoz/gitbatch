package gui

import (
	"fmt"
	"strings"

	"github.com/isacikgoz/gitbatch/core/command"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) commitDetail(ix int) error {

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
	v.Title = detailViewFeature.Title
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
	v.Title = detailViewFeature.Title
	r := gui.getSelectedRepository()
	if ix == 0 {
		p, err := command.PlainDiff(r)
		if err != nil {
			return err
		}
		var s string
		for _, d := range colorizeDiff(p) {
			s = s + "\n" + d
		}
		if len(s) <= 1 {
			return nil
		}

		v.Clear()
		fmt.Fprintf(v, s)
		return nil
	}

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
	if v != nil && v.Title == detailViewFeature.Title {
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
	if v != nil && v.Title == detailViewFeature.Title {
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
