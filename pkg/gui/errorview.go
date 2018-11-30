package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// open an error view to inform user with a message and a useful note 
func (gui *Gui) openErrorView(g *gocui.Gui, message string, note string) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(errorViewFeature.Name, maxX/2-30, maxY/2-3, maxX/2+30, maxY/2+3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = errorViewFeature.Title
		v.Wrap = true
		ps := red.Sprint("Note:") + " " + note
		fmt.Fprintln(v, message)
		fmt.Fprintln(v, ps)
	}
	gui.updateKeyBindingsView(g, errorViewFeature.Name)
	if _, err := g.SetCurrentView(errorViewFeature.Name); err != nil {
		return err
	}
	return nil
}

// close the opened error view
func (gui *Gui) closeErrorView(g *gocui.Gui, v *gocui.View) error {

	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}
