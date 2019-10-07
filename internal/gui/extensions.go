package gui

import (
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// focus to next view
func (gui *Gui) nextViewOfGroup(g *gocui.Gui, v *gocui.View, group []viewFeature) error {
	var focusedViewName string
	if v == nil || v.Name() == group[len(group)-1].Name {
		focusedViewName = group[0].Name
	} else {
		for i := range group {
			if v.Name() == group[i].Name {
				focusedViewName = group[i+1].Name
				break
			}
			if i == len(group)-1 {
				return nil
			}
		}
	}
	if _, err := g.SetCurrentView(focusedViewName); err != nil {
		log.WithFields(log.Fields{
			"view": focusedViewName,
		}).Warn("View cannot be focused.")
		return nil
	}

	return gui.updateKeyBindingsView(g, focusedViewName)
}

// focus to previous view
func (gui *Gui) previousViewOfGroup(g *gocui.Gui, v *gocui.View, group []viewFeature) error {
	var focusedViewName string
	if v == nil || v.Name() == group[0].Name {
		focusedViewName = group[len(group)-1].Name
	} else {
		for i := range group {
			if v.Name() == group[i].Name {
				focusedViewName = group[i-1].Name
				break
			}
			if i == len(group)-1 {
				return nil
			}
		}
	}
	if _, err := g.SetCurrentView(focusedViewName); err != nil {
		log.WithFields(log.Fields{
			"view": focusedViewName,
		}).Warn("View cannot be focused.")
		return nil
	}

	return gui.updateKeyBindingsView(g, focusedViewName)
}

// siwtch the app's mode to fetch
func (gui *Gui) switchToFetchMode(g *gocui.Gui, v *gocui.View) error {
	gui.State.Mode = fetchMode
	return gui.updateKeyBindingsView(g, mainViewFeature.Name)
}

// siwtch the app's mode to pull
func (gui *Gui) switchToPullMode(g *gocui.Gui, v *gocui.View) error {
	gui.State.Mode = pullMode
	return gui.updateKeyBindingsView(g, mainViewFeature.Name)
}

// siwtch the app's mode to merge
func (gui *Gui) switchToMergeMode(g *gocui.Gui, v *gocui.View) error {
	gui.State.Mode = mergeMode
	return gui.updateKeyBindingsView(g, mainViewFeature.Name)
}

// siwtch the app's mode to checkout
func (gui *Gui) switchToCheckoutMode(g *gocui.Gui, v *gocui.View) error {
	gui.State.Mode = checkoutMode
	return gui.updateKeyBindingsView(g, mainViewFeature.Name)
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
	newCy := min(ly, maxY)
	if err := v.SetCursor(cx, newCy); err != nil {
		return err
	}
	err := v.SetOrigin(ox, ly-newCy)
	return err
}

// min finds the minimum value of two int
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// cursor down acts like half-page down for faster scrolling
func (gui *Gui) fastCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, vy := v.Size()
		if len(v.BufferLines())+len(v.ViewBufferLines()) <= vy+oy || len(v.ViewBufferLines()) < vy {
			return nil
		}
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

// closeViewCleanup both updates the keybidings view and focuses to returning view
func (gui *Gui) closeViewCleanup(returningViewName string) (err error) {
	if _, err = gui.g.SetCurrentView(returningViewName); err != nil {
		return err
	}
	err = gui.updateKeyBindingsView(gui.g, returningViewName)
	return err
}

// focus to view same as closeViewCleanup but its just a wrapper for easy reading
func (gui *Gui) focusToView(viewName string) (err error) {
	return gui.closeViewCleanup(viewName)
}
