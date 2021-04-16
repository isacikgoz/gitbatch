package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// open the application controls
// TODO: view size can handled better for such situations like too small
// application area
func (gui *Gui) openCheatSheetView(g *gocui.Gui, _ *gocui.View) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(cheatSheetViewFeature.Name, maxX/2-25, maxY/2-10, maxX/2+25, maxY/2+10)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = cheatSheetViewFeature.Title
		for _, k := range gui.KeyBindings {
			if k.View == mainViewFeature.Name || k.View == "" {
				binding := " " + cyan.Sprint(k.Display) + ": " + k.Description
				fmt.Fprintln(v, binding)
			}
		}
	}
	return gui.focusToView(cheatSheetViewFeature.Name)
}

// close the application controls and do the clean job
func (gui *Gui) closeCheatSheetView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	return gui.closeViewCleanup(mainViewFeature.Name)
}
