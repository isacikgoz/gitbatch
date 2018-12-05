package gui

import (
	log "github.com/sirupsen/logrus"
	"github.com/jroimartin/gocui"
)

var cyclableViews = []string{mainViewFeature.Name,
	remoteViewFeature.Name,
	remoteBranchViewFeature.Name,
	branchViewFeature.Name,
	commitViewFeature.Name}

func (gui *Gui) nextView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == cyclableViews[len(cyclableViews)-1] {
		focusedViewName = cyclableViews[0]
	} else {
		for i := range cyclableViews {
			if v.Name() == cyclableViews[i] {
				focusedViewName = cyclableViews[i+1]
				break
			}
			if i == len(cyclableViews)-1 {
				return nil
			}
		}
	}
	if _, err := g.SetCurrentView(focusedViewName); err != nil {
			log.Warn("Loading view cannot be focused.")
			return nil
	}
	return nil
}

func (gui *Gui) previousView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == cyclableViews[0] {
		focusedViewName = cyclableViews[len(cyclableViews)-1]
	} else {
		for i := range cyclableViews {
			if v.Name() == cyclableViews[i] {
				focusedViewName = cyclableViews[i-1] // TODO: make this work properly
				break
			}
			if i == len(cyclableViews)-1 {
				return nil
			}
		}
	}
	if _, err := g.SetCurrentView(focusedViewName); err != nil {
			log.Warn("Loading view cannot be focused.")
			return nil
	}
	return nil
}
