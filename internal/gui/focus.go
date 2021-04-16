package gui

import (
	"github.com/jroimartin/gocui"
)

var (
	focusViews = []viewFeature{commitViewFeature, dynamicViewFeature, remoteViewFeature, branchViewFeature, stashViewFeature}
)

// set the layout and create views with their default size, name etc. values
// TODO: window sizes can be handled better
func (gui *Gui) focusLayout(g *gocui.Gui) error {

	g.SelFgColor = gocui.ColorGreen
	maxX, maxY := g.Size()
	dx := int(0.35 * float32(maxX))
	rx := int(0.75 * float32(maxX))
	if v, err := g.SetView(mainViewFeature.Name, -2*dx, 0, 0, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFeature.Title
		v.Overwrite = true
	}
	if v, err := g.SetView(remoteViewFeature.Name, rx, 0, maxX-1, int(0.25*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(remoteBranchViewFeature.Name, int(0.25*float32(maxX)), int(0.25*float32(maxY)), int(0.75*float32(maxX)), int(0.75*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteBranchViewFeature.Title
		v.Wrap = false
		v.Overwrite = false
		_, _ = g.SetViewOnBottom(v.Name())
	}
	if v, err := g.SetView(branchViewFeature.Name, rx, int(0.25*float32(maxY)), maxX-1, int(0.75*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = branchViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(stashViewFeature.Name, rx, int(0.75*float32(maxY)), maxX-1, maxY-2); err != nil {
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
		_ = gui.updateKeyBindingsView(g, commitFrameViewFeature.Name)
	}
	return nil
}

// evolve the layout to focus layout and focus to commitview also initialize
// some stuff
func (gui *Gui) focusToRepository(g *gocui.Gui, v *gocui.View) error {
	// mainViews = focusViews
	r := gui.getSelectedRepository()
	if r == nil {
		return nil
	}
	gui.order = focus

	if _, err := g.SetCurrentView(commitViewFeature.Name); err != nil {
		return err
	}
	if err := gui.sendOverviewViewsToBottom(g, v); err != nil {
		return err
	}

	_ = r.State.Branch.InitializeCommits(r)

	if err := gui.renderCommits(r); err != nil {
		return err
	}
	if err := gui.initStashedView(r); err != nil {
		return err
	}
	if err := gui.initFocusStat(r); err != nil {
		return err
	}

	_ = gui.updateKeyBindingsView(g, commitViewFeature.Name)
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.renderMain()
	})
	return nil
}

// return back to overview layout
func (gui *Gui) focusBackToMain(g *gocui.Gui, v *gocui.View) error {
	gui.order = overview

	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	if err := gui.sendFocusViewsToBottom(g, v); err != nil {
		return err
	}
	_ = gui.updateKeyBindingsView(g, mainViewFeature.Name)
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

// send view to bottom so that view won't block others
func (gui *Gui) sendFocusViewsToBottom(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnBottom(branchViewFeature.Name); err != nil {
		return err
	}
	return nil
}

// send view to bottom so that view won't block others
func (gui *Gui) sendOverviewViewsToBottom(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop(branchViewFeature.Name); err != nil {
		return err
	}
	if _, err := g.SetViewOnTop(commitViewFeature.Name); err != nil {
		return err
	}
	return nil
}
