package gui

import (
	"github.com/jroimartin/gocui"
)

var (
	focusViews = []viewFeature{commitViewFeature, dynamicViewFeature, remoteViewFeature, remoteBranchViewFeature, branchViewFeature, stashViewFeature}
)

// set the layout and create views with their default size, name etc. values
// TODO: window sizes can be handled better
func (gui *Gui) focusLayout(g *gocui.Gui) error {

	g.SelFgColor = gocui.ColorGreen
	maxX, maxY := g.Size()
	dx := int(0.35 * float32(maxX))
	rx := int(0.80 * float32(maxX))
	if v, err := g.SetView(mainViewFeature.Name, -2*dx, 0, 0, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFeature.Title
		v.Overwrite = true
	}
	if v, err := g.SetView(remoteViewFeature.Name, rx, 0, maxX-1, int(0.15*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(remoteBranchViewFeature.Name, rx, int(0.15*float32(maxY)), maxX-1, int(0.50*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteBranchViewFeature.Title
		v.Wrap = false
		v.Overwrite = false
	}
	if v, err := g.SetView(branchViewFeature.Name, rx, int(0.50*float32(maxY)), maxX-1, int(0.85*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = branchViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(stashViewFeature.Name, rx, int(0.85*float32(maxY)), maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stashViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(commitViewFeature.Name, 0, 0, dx-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = commitViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(dynamicViewFeature.Name, dx, 0, rx-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = dynamicViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(keybindingsViewFeature.Name, -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorBlack
		v.Frame = false
		gui.updateKeyBindingsView(g, commitFrameViewFeature.Name)
	}
	return nil
}

func (gui *Gui) focusToRepository(g *gocui.Gui, v *gocui.View) error {
	// mainViews = focusViews
	r := gui.getSelectedRepository()
	gui.order = focus

	if _, err := g.SetCurrentView(commitViewFeature.Name); err != nil {
		return err
	}

	r.State.Branch.InitializeCommits(r)

	if err := gui.renderCommits(r); err != nil {
		return err
	}
	if err := gui.initStashedView(r); err != nil {
		return err
	}
	if err := gui.initFocusStat(r); err != nil {
		return err
	}

	gui.updateKeyBindingsView(g, commitViewFeature.Name)
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.renderMain()
	})
	return nil
}

func (gui *Gui) focusBackToMain(g *gocui.Gui, v *gocui.View) error {
	// mainViews = overviewViews
	gui.order = overview

	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

// focus to next view
func (gui *Gui) nextFocusView(g *gocui.Gui, v *gocui.View) error {
	return gui.nextViewOfGroup(g, v, focusViews)
}

// focus to previous view
func (gui *Gui) previousFocusView(g *gocui.Gui, v *gocui.View) error {
	return gui.previousViewOfGroup(g, v, focusViews)
}
