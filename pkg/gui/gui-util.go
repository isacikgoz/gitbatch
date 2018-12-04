package gui

import (
	"sort"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/isacikgoz/gitbatch/pkg/helpers"
	"github.com/jroimartin/gocui"
)

// refreshes the side views of the application for given git.RepoEntity struct
func (gui *Gui) refreshViews(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	if err = gui.updateRemotes(g, entity); err != nil {
		return err
	}
	if err = gui.updateBranch(g, entity); err != nil {
		return err
	}
	if err = gui.updateRemoteBranches(g, entity); err != nil {
		return err
	}
	if err = gui.updateCommits(g, entity); err != nil {
		return err
	}
	return err
}

// siwtch the app mode
// TODO: switching can be made with conventional iteration
func (gui *Gui) switchMode(g *gocui.Gui, v *gocui.View) error {
	switch mode := gui.State.Mode.ModeID; mode {
	case FetchMode:
		gui.State.Mode = pullMode
	case PullMode:
		gui.State.Mode = mergeMode
	case MergeMode:
		gui.State.Mode = fetchMode
	default:
		gui.State.Mode = fetchMode
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

// bring the view on the top by its name
func (gui *Gui) setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

// if the cursor down past the last item, move it to the last line
func (gui *Gui) correctCursor(v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	width, height := v.Size()
	maxY := height - 1
	ly := width - 1
	if oy+cy <= ly {
		return nil
	}
	newCy := helpers.Min(ly, maxY)
	if err := v.SetCursor(cx, newCy); err != nil {
		return err
	}
	if err := v.SetOrigin(ox, ly-newCy); err != nil {
		return err
	}
	return nil
}

// this function handles the iteration of a side view and set its origin point
// so that the selected line can be in the middle of the view
func (gui *Gui) smartAnchorRelativeToLine(v *gocui.View, currentindex, totallines int) error {
	_, y := v.Size()
	if currentindex >= int(0.5*float32(y)) && totallines-currentindex+int(0.5*float32(y)) >= y {
		if err := v.SetOrigin(0, currentindex-int(0.5*float32(y))); err != nil {
			return err
		}
	} else if totallines-currentindex < y && totallines > y {
		if err := v.SetOrigin(0, totallines-y); err != nil {
			return err
		}
	} else if totallines-currentindex <= int(0.5*float32(y)) && totallines > y-1 && currentindex > y {
		if err := v.SetOrigin(0, currentindex-int(0.5*float32(y))); err != nil {
			return err
		}
	} else {
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
	}
	return nil
}

// this function writes the given text to rgiht hand side of the view
// cx and cy values are important to get the cursor to its old position
func writeRightHandSide(v *gocui.View, text string, cx, cy int) error {
	runes := []rune(text)
	tl := len(runes)
	lx, _ := v.Size()
	v.MoveCursor(lx-tl, cy-1, true)
	for i := tl - 1; i >= 0; i-- {
		v.EditDelete(true)
		v.EditWrite(runes[i])
	}
	v.SetCursor(cx, cy)
	return nil
}

// sortByName sorts the repositories by A to Z order
func (gui *Gui) sortByName(g *gocui.Gui, v *gocui.View) error {
	sort.Sort(git.Alphabetical(gui.State.Repositories))
	gui.refreshAfterSort(g)
	return nil
}

// sortByMod sorts the repositories according to last modifed date
// the top element will be the last modified
func (gui *Gui) sortByMod(g *gocui.Gui, v *gocui.View) error {
	sort.Sort(git.LastModified(gui.State.Repositories))
	gui.refreshAfterSort(g)
	return nil
}

// utility function that refreshes main and side views after that
func (gui *Gui) refreshAfterSort(g *gocui.Gui) error {
	gui.refreshMain(g)
	entity := gui.getSelectedRepository()
	gui.refreshViews(g, entity)
	return nil
}

// cursor down acts like half-page down for faster scrolling
func (gui *Gui) fastCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, vy := v.Size()

		// TODO: do something when it hits bottom
		if err := v.SetOrigin(ox, oy+vy/2); err != nil {
			return err
		}
	}
	return nil
}

// cursor up acts like half-page up for faster scrolling
func (gui *Gui) fastCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, vy := v.Size()

		if oy-vy/2 > 0 {
			if err := v.SetOrigin(ox, oy-vy/2); err != nil {
				return err
			}
		} else if oy-vy/2 <= 0 {
			if err := v.SetOrigin(0, 0); err != nil {
				return err
			}
		}
	}
	return nil
}
