package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var stashReturnView string

func (gui *Gui) openStashView(g *gocui.Gui, returnViewName string) error {
	maxX, maxY := g.Size()
	stashReturnView = returnViewName
	v, err := g.SetView("stash", maxX/2-30, maxY/2-3, maxX/2+30, maxY/2+3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Stashed Items "
		v.Wrap = true
		entity := gui.getSelectedRepository()
		stashedItems, err := entity.LoadStashedItems()
		if err != nil {
			return err
		}
		for _, stashedItem := range stashedItems {
			fmt.Fprintln(v, stashedItem)
		}

	}
	gui.updateKeyBindingsView(g, "stash")
	if _, err := g.SetCurrentView("stash"); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) closeStashView(g *gocui.Gui, v *gocui.View) error {

	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	if _, err := g.SetCurrentView(stashReturnView); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, stashReturnView)
	return nil
}

func (gui *Gui) showStashView(g *gocui.Gui, v *gocui.View) (err error) {
	gui.openStashView(g, mainViewFeature.Name)
	return nil
}
