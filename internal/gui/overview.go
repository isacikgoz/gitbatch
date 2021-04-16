package gui

import (
	"github.com/jroimartin/gocui"
)

// set the layout and create views with their default size, name etc. values
// TODO: window sizes can be handled better
func (gui *Gui) overviewLayout(g *gocui.Gui) error {
	g.SelFgColor = gocui.ColorDefault
	maxX, maxY := g.Size()
	// dx := int(0.55 * float32(maxX))
	dx := -2
	if v, err := g.SetView(mainViewFrameFeature.Name, 0, 0, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFrameFeature.Title
	}
	if v, err := g.SetView(mainViewFeature.Name, 1, 1, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFeature.Title
		v.Overwrite = true
		if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
			return err
		}
		v.Frame = false
	}
	if v, err := g.SetView(remoteViewFeature.Name, dx, 0, -1, int(0.15*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(remoteBranchViewFeature.Name, dx, int(0.15*float32(maxY)), -1, int(0.55*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteBranchViewFeature.Title
		v.Wrap = false
		v.Overwrite = false
	}
	if v, err := g.SetView(branchViewFeature.Name, int(0.25*float32(maxX)), int(0.25*float32(maxY)), int(0.75*float32(maxX)), int(0.75*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = branchViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
		_, _ = g.SetViewOnBottom(v.Name())
	}
	if v, err := g.SetView(batchBranchViewFeature.Name, int(0.25*float32(maxX)), int(0.25*float32(maxY)), int(0.75*float32(maxX)), int(0.75*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = batchBranchViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
		_, _ = g.SetViewOnBottom(v.Name())
	}
	if v, err := g.SetView(suggestBranchViewFeature.Name, int(0.30*float32(maxX)), int(0.45*float32(maxY)), int(0.70*float32(maxX)), int(0.55*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = suggestBranchViewFeature.Title
		v.Editable = true
		v.Wrap = false
		v.Autoscroll = false
		_, _ = g.SetViewOnBottom(v.Name())
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
		_ = gui.updateKeyBindingsView(g, mainViewFeature.Name)
	}
	return nil
}

// close confirmation view
func (gui *Gui) openBranchesView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop(branchViewFeature.Name); err != nil {
		return err
	}
	return gui.focusToView(branchViewFeature.Name)
}

// close confirmation view
func (gui *Gui) closeBranchesView(g *gocui.Gui, v *gocui.View) error {
	if gui.order == focus {
		return nil
	}
	if _, err := g.SetViewOnBottom(branchViewFeature.Name); err != nil {
		return err
	}
	return gui.focusToView(mainViewFeature.Name)
}
