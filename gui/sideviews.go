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

// refreshes the side views of the application for given repository.Repository struct
func (gui *Gui) renderSideViews(r *git.Repository) error {
	if r == nil {
		return nil
	}

	if err := gui.renderRemotes(r); err != nil {
		return err
	}
	if err := gui.renderBranch(r); err != nil {
		return err
	}
	if err := gui.renderRemoteBranches(r); err != nil {
		return err
	}
	if err := gui.renderCommits(r); err != nil {
		return err
	}
	return nil
}

// updates the remotesview for given entity
func (gui *Gui) renderRemotes(r *git.Repository) error {
	var err error
	out, err := gui.g.View(remoteViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	totalRemotes := len(r.Remotes)
	if totalRemotes > 0 {
		for i, rm := range r.Remotes {
			_, shortURL := trimRemoteURL(rm.URL[0])
			if rm.Name == r.State.Remote.Name {
				currentindex = i
				fmt.Fprintln(out, selectionIndicator+rm.Name+": "+shortURL)
				continue
			}
			fmt.Fprintln(out, tab+rm.Name+": "+shortURL)
		}
		if err = gui.smartAnchorRelativeToLine(out, currentindex, totalRemotes); err != nil {
			return err
		}
	}
	return nil
}

// updates the remotebranchview for given entity
func (gui *Gui) renderRemoteBranches(r *git.Repository) error {
	var err error
	out, err := gui.g.View(remoteBranchViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	trb := len(r.State.Remote.Branches)
	if trb > 0 {
		for i, rm := range r.State.Remote.Branches {
			if rm.Name == r.State.Remote.Branch.Name {
				currentindex = i
				fmt.Fprintln(out, selectionIndicator+rm.Name)
				continue
			}
			fmt.Fprintln(out, tab+rm.Name)
		}
		if err = gui.smartAnchorRelativeToLine(out, currentindex, trb); err != nil {
			return err
		}
	}
	return nil
}

// updates the branchview for given entity
func (gui *Gui) renderBranch(r *git.Repository) error {
	var err error
	out, err := gui.g.View(branchViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	totalbranches := len(r.Branches)
	for i, b := range r.Branches {
		if b.Name == r.State.Branch.Name {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator+b.Name)
			continue
		}
		fmt.Fprintln(out, tab+b.Name)
	}

	return gui.smartAnchorRelativeToLine(out, currentindex, totalbranches)
}

// updates the commitsview for given entity
func (gui *Gui) renderCommits(r *git.Repository) error {
	var err error
	out, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	totalcommits := len(r.Commits)
	for i, c := range r.Commits {
		if c.Hash == r.State.Commit.Hash {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator+commitLabel(c))
			continue
		}
		fmt.Fprintln(out, tab+commitLabel(c))
	}
	return gui.smartAnchorRelativeToLine(out, currentindex, totalcommits)
}

// cursor down variant for sideviews
func (gui *Gui) sideViewsNextItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	r := gui.getSelectedRepository()
	switch viewName := v.Name(); viewName {
	case remoteBranchViewFeature.Name:
		return r.State.Remote.NextRemoteBranch(r)
	case remoteViewFeature.Name:
		return r.NextRemote()
	case branchViewFeature.Name:
		if err = r.Checkout(r.NextBranch()); err != nil {
			err = gui.openErrorView(g, err.Error(),
				"You should manually resolve this issue",
				branchViewFeature.Name)
			return err
		}
	case commitViewFeature.Name:
		r.NextCommit()
		return gui.renderCommits(r)
	}
	return err
}

// cursor up variant for sideviews
func (gui *Gui) sideViewsPreviousItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	r := gui.getSelectedRepository()
	switch viewName := v.Name(); viewName {
	case remoteBranchViewFeature.Name:
		return r.State.Remote.PreviousRemoteBranch(r)
	case remoteViewFeature.Name:
		return r.PreviousRemote()
	case branchViewFeature.Name:
		if err = r.Checkout(r.PreviousBranch()); err != nil {
			err = gui.openErrorView(g, err.Error(),
				"You should manually resolve this issue",
				branchViewFeature.Name)
			return err
		}
	case commitViewFeature.Name:
		r.PreviousCommit()
		return gui.renderCommits(r)
	}
	return err
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
