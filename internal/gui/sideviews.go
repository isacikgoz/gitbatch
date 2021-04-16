package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/internal/command"
	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/jroimartin/gocui"
)

var (
	confirmationViewFeature = viewFeature{Name: "confirmation", Title: " Confirmation "}
	sideViews               = []viewFeature{remoteViewFeature, remoteBranchViewFeature, branchViewFeature, batchBranchViewFeature}
)

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) sideCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		ly := len(v.BufferLines()) - 1

		// if we are at the end we just return
		if cy+oy == ly-1 {
			return nil
		}
		v.EditDelete(true)
		_ = adjustAnchor(cy+oy+1, ly, v)
	}
	return nil
}

// moves the cursor upwards for the main view
func (gui *Gui) sideCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		_, cy := v.Cursor()
		ly := len(v.BufferLines()) - 1
		v.EditDelete(true)
		_ = adjustAnchor(cy+oy-1, ly, v)
	}
	return nil
}

func (gui *Gui) renderSideViews(r *git.Repository) error {
	if r == nil {
		return nil
	}
	if err := gui.resetSideCursors(); err != nil {
		return err
	}
	if err := gui.renderBranches(r); err != nil {
		return err
	}
	if err := gui.renderRemoteBranches(r); err != nil {
		return err
	}
	if err := gui.renderRemotes(r); err != nil {
		return err
	}
	return nil
}

// updates the branchesview for given entity
func (gui *Gui) renderBranches(r *git.Repository) error {
	v, err := gui.g.View(branchViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	bs := r.Branches
	if r.State.Branch == nil {
		return nil
	}
	bc := r.State.Branch
	si := 0
	for i, b := range bs {
		if b.Name == bc.Name {
			si = i
			fmt.Fprintln(v, ws+green.Sprint(b.Name))
			continue
		}
		fmt.Fprintln(v, tab+b.Name)
	}
	_ = adjustAnchor(si, len(bs), v)
	return nil
}

// updates the remotebranchesview for given entity
func (gui *Gui) renderRemotes(r *git.Repository) error {
	v, err := gui.g.View(remoteViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	rs := r.Remotes
	rc := r.State.Remote
	si := 0
	for i, rb := range rs {
		shortURL := "URL not found."
		if len(rb.URL) > 0 {
			_, shortURL = trimRemoteURL(rb.URL[0])
		}
		if rb.Name == rc.Name {
			si = i
			fmt.Fprintln(v, ws+green.Sprint(rb.Name+": "+shortURL))
			continue
		}
		fmt.Fprintln(v, tab+rb.Name+": "+shortURL)
	}
	_ = adjustAnchor(si, len(rs), v)
	return nil
}

// updates the remotesview for given entity
func (gui *Gui) renderRemoteBranches(r *git.Repository) error {
	v, err := gui.g.View(remoteBranchViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	rs := r.State.Remote.Branches
	// rc := r.State.Remote.Branch
	si := 0
	for _, rb := range rs {
		// if rb.Name == rc.Name {
		// 	si = i
		// 	fmt.Fprintln(v, ws+green.Sprint(rb.Name))
		// 	continue
		// }
		fmt.Fprintln(v, tab+rb.Name)
	}
	_ = adjustAnchor(si, len(rs), v)
	return nil
}

func (gui *Gui) selectSideItem(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	r := gui.getSelectedRepository()
	ix := oy + cy
	var err error
	if v.Name() == commitViewFeature.Name {
		r.State.Branch.State.Commit = r.State.Branch.Commits[ix]
		err = gui.renderCommits(r)
	} else if v.Name() == branchViewFeature.Name {
		_ = r.Checkout(r.Branches[ix])
		err = gui.renderBranches(r)
	} else if v.Name() == remoteBranchViewFeature.Name {
		// r.State.Remote.Branch = r.State.Remote.Branches[ix]
		err = gui.renderRemoteBranches(r)
	} else if v.Name() == remoteViewFeature.Name {
		r.State.Remote = r.Remotes[ix]
		_ = r.Refresh()
		err = gui.renderRemotes(r)
	} else if v.Name() == batchBranchViewFeature.Name {
		gui.State.targetBranch = gui.State.totalBranches[ix].BranchName
		err = gui.renderBatchBranches(false)
	}
	return err
}

func adjustAnchor(i, r int, v *gocui.View) error {
	_, y := v.Size()
	if i >= int(0.5*float32(y)) && r-i+int(0.5*float32(y)) >= y {
		if err := v.SetOrigin(0, i-int(0.5*float32(y))); err != nil {
			return err
		}
	} else if r-i < y && r > y {
		if err := v.SetOrigin(0, r-y); err != nil {
			return err
		}
	} else if r-i <= int(0.5*float32(y)) && r > y-1 && i > y {
		if err := v.SetOrigin(0, i-int(0.5*float32(y))); err != nil {
			return err
		}
	} else {
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
	}
	_, oy := v.Origin()
	c := i - oy
	_ = v.SetCursor(0, c)
	v.EditWrite('→')
	return nil
}

func (gui *Gui) resetSideCursors() error {
	for _, vf := range sideViews {
		v, err := gui.g.View(vf.Name)
		if err != nil {
			return err
		}
		_ = v.SetCursor(0, 0)
	}
	return nil
}

// basically does fetch --prune
func (gui *Gui) syncRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	return command.Fetch(r, &command.FetchOptions{
		RemoteName: r.State.Remote.Name,
		Prune:      true,
	})
}

// opens a confirmation view for setting default merge branch
func (gui *Gui) setUpstreamToBranch(g *gocui.Gui, _ *gocui.View) error {
	maxX, maxY := g.Size()

	r := gui.getSelectedRepository()
	v, err := g.SetView(confirmationViewFeature.Name, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "branch."+r.State.Branch.Name+"."+"remote"+"="+r.State.Remote.Name)
		fmt.Fprintln(v, "branch."+r.State.Branch.Name+"."+"merge"+"="+r.State.Branch.Reference.Name().String())
	}
	return gui.focusToView(confirmationViewFeature.Name)
}

// add config for upstream merge
func (gui *Gui) confirmSetUpstreamToBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	r := gui.getSelectedRepository()
	if err = command.AddConfig(r, &command.ConfigOptions{
		Section: "branch." + r.State.Branch.Name,
		Option:  "remote",
		Site:    command.ConfigSiteLocal,
	}, r.State.Remote.Name); err != nil {
		return err
	}
	if err = command.AddConfig(r, &command.ConfigOptions{
		Section: "branch." + r.State.Branch.Name,
		Option:  "merge",
		Site:    command.ConfigSiteLocal,
	}, r.State.Branch.Reference.Name().String()); err != nil {
		return err
	}
	_ = r.Refresh()
	return gui.closeConfirmationView(g, v)
}

// close confirmation view
func (gui *Gui) closeConfirmationView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	return gui.closeViewCleanup(branchViewFeature.Name)
}

// close confirmation view
func (gui *Gui) openRemoteBranchesView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop(remoteBranchViewFeature.Name); err != nil {
		return err
	}
	return gui.focusToView(remoteBranchViewFeature.Name)
}

// close confirmation view
func (gui *Gui) closeRemoteBranchesView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnBottom(remoteBranchViewFeature.Name); err != nil {
		return err
	}
	return gui.focusToView(remoteViewFeature.Name)
}
