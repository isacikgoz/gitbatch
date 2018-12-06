package gui

import (
	log "github.com/sirupsen/logrus"
	"github.com/jroimartin/gocui"
)

// focus to next view
func (gui *Gui) nextView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == mainViews[len(mainViews)-1].Name {
		focusedViewName = mainViews[0].Name
	} else {
		for i := range mainViews {
			if v.Name() == mainViews[i].Name {
				focusedViewName = mainViews[i+1].Name
				break
			}
			if i == len(mainViews)-1 {
				return nil
			}
		}
	}
	if _, err := g.SetCurrentView(focusedViewName); err != nil {
			log.Warn("Loading view cannot be focused.")
			return nil
	}
	gui.updateKeyBindingsView(g, focusedViewName)
	return nil
}

// focus to previous view
func (gui *Gui) previousView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == mainViews[0].Name {
		focusedViewName = mainViews[len(mainViews)-1].Name
	} else {
		for i := range mainViews {
			if v.Name() == mainViews[i].Name {
				focusedViewName = mainViews[i-1].Name
				break
			}
			if i == len(mainViews)-1 {
				return nil
			}
		}
	}
	if _, err := g.SetCurrentView(focusedViewName); err != nil {
			log.Warn("Loading view cannot be focused.")
			return nil
	}
	gui.updateKeyBindingsView(g, focusedViewName)
	return nil
}
