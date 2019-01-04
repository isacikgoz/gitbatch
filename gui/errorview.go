package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var errorReturnView string

// open an error view to inform user with a message and a useful note
func (gui *Gui) openErrorView(g *gocui.Gui, message, note, returnViewName string) error {
	maxX, maxY := g.Size()
	errorReturnView = returnViewName
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
	return gui.focusToView(errorViewFeature.Name)
}

// close the opened error view
func (gui *Gui) closeErrorView(g *gocui.Gui, v *gocui.View) error {

	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	return gui.closeViewCleanup(errorReturnView)
}
