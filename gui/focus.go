package gui

import (
	"github.com/jroimartin/gocui"
)

var (
	focusViews = []viewFeature{commitViewFeature, detailViewFeature, remoteViewFeature, remoteBranchViewFeature, branchViewFeature, stashViewFeature}
)

// set the layout and create views with their default size, name etc. values
// TODO: window sizes can be handled better
func (gui *Gui) focusLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	dx := int(0.20 * float32(maxX))
	rx := int(0.60 * float32(maxX))
	if v, err := g.SetView(mainViewFeature.Name, -2*dx, 0, 0, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFeature.Title
		v.Overwrite = true
	}
	if v, err := g.SetView(remoteViewFeature.Name, 0, 0, dx-1, int(0.15*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(remoteBranchViewFeature.Name, 0, int(0.15*float32(maxY)), dx-1, int(0.50*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteBranchViewFeature.Title
		v.Wrap = false
		v.Overwrite = false
	}
	if v, err := g.SetView(branchViewFeature.Name, 0, int(0.50*float32(maxY)), dx-1, int(0.85*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = branchViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(stashViewFeature.Name, 0, int(0.85*float32(maxY)), dx-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stashViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(commitViewFeature.Name, dx, 0, rx-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = commitViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(detailViewFeature.Name, rx, 0, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = detailViewFeature.Title
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
