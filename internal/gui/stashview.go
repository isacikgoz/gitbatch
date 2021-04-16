package gui

import (
	"fmt"
	"strings"

	"github.com/isacikgoz/gitbatch/internal/command"
	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/jroimartin/gocui"
)

// initialize stashed items
func (gui *Gui) initStashedView(r *git.Repository) error {
	v, err := gui.g.View(stashViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	st := r.Stasheds
	for _, s := range st {
		fmt.Fprintf(v, " %d %s: %s\n", s.StashID, cyan.Sprint(s.BranchName), s.Description)
	}
	if len(st) > 0 {
		_ = adjustAnchor(0, len(st), v)
	}
	return nil
}

// moves the cursor downwards for the stash view and if it goes to bottom it
// prevents from going further
func (gui *Gui) stashCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		ly := len(v.BufferLines()) - 1

		// if we are at the end we just return
		if cy+oy == ly-1 || ly < 0 {
			return nil
		}
		v.EditDelete(true)
		pos := cy + oy + 1
		_ = adjustAnchor(pos, ly, v)
		_ = gui.renderStashDiff(pos)
	}
	return nil
}

// moves the cursor upwards for the stash view
func (gui *Gui) stashCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		_, cy := v.Cursor()
		ly := len(v.BufferLines()) - 1
		// if we are at the end we just return
		if ly < 0 {
			return nil
		}
		v.EditDelete(true)
		pos := cy + oy - 1
		_ = adjustAnchor(pos, ly, v)
		if pos >= 0 {
			_ = gui.renderStashDiff(cy + oy - 1)
		}
	}
	return nil
}

// open stashed item diff which given with index
func (gui *Gui) renderStashDiff(ix int) error {
	r := gui.getSelectedRepository()
	st := r.Stasheds
	if len(st) <= 0 {
		return nil
	}
	v, err := gui.g.View(dynamicViewFeature.Name)
	if err != nil {
		return err
	}
	v.Title = string(StashDiffMode)
	if err := gui.updateDynamicKeybindings(); err != nil {
		return err
	}
	v.Clear()
	d, err := command.StashDiff(r, st[ix].StashID)
	if err != nil {
		return err
	}
	s := colorizeDiff(d)
	fmt.Fprintf(v, "%s", strings.Join(s, "\n"))
	return nil
}

// open stash item diff
func (gui *Gui) stashDiff(g *gocui.Gui, v *gocui.View) error {

	_, oy := v.Origin()
	_, cy := v.Cursor()

	return gui.renderStashDiff(oy + cy)
}

// pop out the stash
func (gui *Gui) stashPop(g *gocui.Gui, v *gocui.View) error {
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
	_ = r.Refresh()
	if err := gui.focusToRepository(g, v); err != nil {
		return err
	}
	if err := gui.initStashedView(r); err != nil {
		return err
	}
	return nil
}
