package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var (
	statusHeaderViewFeature = viewFeature{Name: "status-header", Title: " Status Header "}
	statusViewFeature       = viewFeature{Name: "status", Title: " Status "}
	stageViewFeature        = viewFeature{Name: "staged", Title: " Staged "}
	unstageViewFeature      = viewFeature{Name: "unstaged", Title: " Unstaged "}
	stashViewFeature        = viewFeature{Name: "stash", Title: " Stash "}
)

// open the status layout
func (gui *Gui) openStatusView(g *gocui.Gui, v *gocui.View) error {
	gui.openStatusHeaderView(g)
	gui.openStageView(g)
	gui.openUnStagedView(g)
	gui.openStashView(g)
	return nil
}

// iteration handler for the status layout
func (gui *Gui) nextStatusView(g *gocui.Gui, v *gocui.View) error {
	var err error
	return err
}

// header og the status layout
func (gui *Gui) openStatusHeaderView(g *gocui.Gui) error {
	maxX, _ := g.Size()
	entity := gui.getSelectedRepository()
	v, err := g.SetView(statusHeaderViewFeature.Name, 6, 2,  maxX-6, 4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, entity.AbsPath)
		// v.Frame = false
		v.Wrap = true
	}
	gui.updateKeyBindingsView(g, statusHeaderViewFeature.Name)
	if _, err := g.SetCurrentView(statusHeaderViewFeature.Name); err != nil {
		return err
	}
	return nil
}

// staged view
func (gui *Gui) openStageView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(stageViewFeature.Name, 6, 5, maxX/2-1, int(0.75*float32(maxY))-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stageViewFeature.Title
		v.Wrap = true
	}
	return nil
}

// not staged view
func (gui *Gui) openUnStagedView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(unstageViewFeature.Name, maxX/2+1, 5, maxX-6, int(0.75*float32(maxY))-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = unstageViewFeature.Title
		v.Wrap = true
	}
	return nil
}


// stash view
func (gui *Gui) openStashView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(stashViewFeature.Name, 6, int(0.75*float32(maxY)), maxX-6, maxY-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stashViewFeature.Title
		v.Wrap = true
	}
	return nil
}

// close the opened stat views
func (gui *Gui) closeStatusView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(stashViewFeature.Name); err != nil {
		return err
	}
	if err := g.DeleteView(unstageViewFeature.Name); err != nil {
		return err
	}
	if err := g.DeleteView(stageViewFeature.Name); err != nil {
		return err
	}
	if err := g.DeleteView(statusHeaderViewFeature.Name); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

 		// 	{
		// 	View:        statusHeaderViewFeature.Name,
		// 	Key:         'c',
		// 	Modifier:    gocui.ModNone,
		// 	Handler:     gui.closeStatusView,
		// 	Display:     "c",
		// 	Description: "close/cancel",
		// 	Vital:       true,
		// },

		//  {
		// 	View:        mainViewFeature.Name,
		// 	Key:         't',
		// 	Modifier:    gocui.ModNone,
		// 	Handler:     gui.openStatusView,
		// 	Display:     "t",
		// 	Description: "Open Status",
		// 	Vital:       true,
		// },