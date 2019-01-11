package gui

import (
	"github.com/jroimartin/gocui"
)

var (
	overviewViews = []viewFeature{mainViewFeature, remoteViewFeature, remoteBranchViewFeature, branchViewFeature}
)

// set the layout and create views with their default size, name etc. values
// TODO: window sizes can be handled better
func (gui *Gui) overviewLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	dx := int(0.55 * float32(maxX))
	if v, err := g.SetView(mainViewFeature.Name, 0, 0, dx-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFeature.Title
		v.Overwrite = true
		if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
			return err
		}
	}
	if v, err := g.SetView(remoteViewFeature.Name, dx, 0, maxX-1, int(0.15*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(remoteBranchViewFeature.Name, dx, int(0.15*float32(maxY)), maxX-1, int(0.55*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteBranchViewFeature.Title
		v.Wrap = false
		v.Overwrite = false
	}
	if v, err := g.SetView(branchViewFeature.Name, dx, int(0.55*float32(maxY)), maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = branchViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(stashViewFeature.Name, -1*int(0.20*float32(maxX)), 0, -1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stashViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(commitViewFeature.Name, -1*int(0.20*float32(maxX)), 0, -1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = commitViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(dynamicViewFeature.Name, -1*int(0.20*float32(maxX)), 0, -1, maxY); err != nil {
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
		gui.updateKeyBindingsView(g, mainViewFeature.Name)
	}
	return nil
}
