package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/core/command"
	"github.com/isacikgoz/gitbatch/core/git"
	"github.com/jroimartin/gocui"
)

var (
	confirmationViewFeature = viewFeature{Name: "confirmation", Title: " Confirmation "}
	sideViews               = []viewFeature{remoteViewFeature, remoteBranchViewFeature, branchViewFeature, commitViewFeature}
)

// basically does fetch --prune
func (gui *Gui) syncRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	return command.Fetch(r, &command.FetchOptions{
		RemoteName: r.State.Remote.Name,
		Prune:      true,
	})
}

// opens a confirmation view for setting default merge branch
func (gui *Gui) setUpstreamToBranch(g *gocui.Gui, v *gocui.View) error {
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
	r.Refresh()
	return gui.closeConfirmationView(g, v)
}

// close confirmation view
func (gui *Gui) closeConfirmationView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	return gui.closeViewCleanup(branchViewFeature.Name)
}

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) sideCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		ox, oy := v.Origin()
		ly := len(v.BufferLines()) - 2

		// if we are at the end we just return
		if cy+oy == ly {
			return nil
		}

		v.EditDelete(true)
		if err := v.SetCursor(0, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		v.EditWrite('*')
	}
	return nil
	// return gui.renderSide(gui.getSelectedRepository(), v)
}

// moves the cursor upwards for the main view
func (gui *Gui) sideCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, cy := v.Cursor()
		v.EditDelete(true)
		if err := v.SetCursor(0, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		v.EditWrite('*')
	}
	return nil
	// return gui.renderSide(gui.getSelectedRepository(), v)
}

func (gui *Gui) renderSideViews(r *git.Repository) error {
	gui.resetCursors()
	gui.renderCommitsNew(r)
	gui.renderBranchesNew(r)
	gui.renderRemoteBranchesNew(r)
	gui.renderRemotesNew(r)

	return nil
}

func (gui *Gui) renderSide(r *git.Repository, v *gocui.View) error {
	if v.Name() == commitViewFeature.Name {
		gui.renderCommitsNew(r)
	} else if v.Name() == branchViewFeature.Name {
		gui.renderBranchesNew(r)
	} else if v.Name() == remoteBranchViewFeature.Name {
		gui.renderRemoteBranchesNew(r)
	} else if v.Name() == remoteViewFeature.Name {
		gui.renderRemotesNew(r)
	}
	return nil
}

// updates the commitsview for given entity
func (gui *Gui) renderCommitsNew(r *git.Repository) error {
	v, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	cs := r.State.Branch.Commits
	bc := r.State.Branch.State.Commit
	si := 0
	for i, c := range cs {
		if c.Hash == bc.Hash {
			si = i
			fmt.Fprintln(v, ws+green.Sprint(commitLabel(c)))
			continue
		}

		fmt.Fprintln(v, tab+commitLabel(c))
	}
	adjustAnchor(si, v)
	return nil
}

// updates the branchesview for given entity
func (gui *Gui) renderBranchesNew(r *git.Repository) error {
	v, err := gui.g.View(branchViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	bs := r.Branches
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
	adjustAnchor(si, v)
	return nil
}

// updates the remotebranchesview for given entity
func (gui *Gui) renderRemotesNew(r *git.Repository) error {
	v, err := gui.g.View(remoteViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	rs := r.Remotes
	rc := r.State.Remote
	si := 0
	for i, rb := range rs {
		_, shortURL := trimRemoteURL(rb.URL[0])
		if rb.Name == rc.Name {
			si = i
			fmt.Fprintln(v, ws+green.Sprint(rb.Name+": "+shortURL))
			continue
		}
		fmt.Fprintln(v, tab+rb.Name+": "+shortURL)
	}
	adjustAnchor(si, v)
	return nil
}

// updates the remotesview for given entity
func (gui *Gui) renderRemoteBranchesNew(r *git.Repository) error {
	v, err := gui.g.View(remoteBranchViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	rs := r.State.Remote.Branches
	rc := r.State.Remote.Branch
	si := 0
	for i, rb := range rs {
		if rb.Name == rc.Name {
			si = i
			fmt.Fprintln(v, ws+green.Sprint(rb.Name))
			continue
		}
		fmt.Fprintln(v, tab+rb.Name)
	}
	adjustAnchor(si, v)
	return nil
}

func (gui *Gui) selectSideItem(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	r := gui.getSelectedRepository()
	ix := oy + cy
	if v.Name() == commitViewFeature.Name {
		r.State.Branch.State.Commit = r.State.Branch.Commits[ix]
	} else if v.Name() == branchViewFeature.Name {
		r.Checkout(r.Branches[ix])
	} else if v.Name() == remoteBranchViewFeature.Name {
		r.State.Remote.Branch = r.State.Remote.Branches[ix]
	} else if v.Name() == remoteViewFeature.Name {
		r.State.Remote = r.Remotes[ix]
	}

	return gui.renderSide(r, v)
}

func adjustAnchor(i int, v *gocui.View) error {
	_, oy := v.Origin()
	v.SetOrigin(0, oy)
	v.SetCursor(0, i-oy)
	v.EditWrite('*')
	return nil
}

func (gui *Gui) resetCursors() error {
	for _, vf := range sideViews {
		v, err := gui.g.View(vf.Name)
		if err != nil {
			return err
		}
		v.SetCursor(0, 0)
	}
	return nil
}
