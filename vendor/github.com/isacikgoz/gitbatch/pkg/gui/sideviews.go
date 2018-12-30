package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

var (
	confirmationViewFeature = viewFeature{Name: "confirmation", Title: " Confirmation "}
	sideViews               = []viewFeature{remoteViewFeature, remoteBranchViewFeature, branchViewFeature, commitViewFeature}
)

// refreshes the side views of the application for given git.RepoEntity struct
func (gui *Gui) renderSideViews(e *git.RepoEntity) error {
	if e == nil {
		return nil
	}
	var err error
	if err = gui.renderRemotes(e); err != nil {
		return err
	}
	if err = gui.renderBranch(e); err != nil {
		return err
	}
	if err = gui.renderRemoteBranches(e); err != nil {
		return err
	}
	if err = gui.renderCommits(e); err != nil {
		return err
	}
	return err
}

// updates the remotesview for given entity
func (gui *Gui) renderRemotes(e *git.RepoEntity) error {
	var err error
	out, err := gui.g.View(remoteViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	totalRemotes := len(e.Remotes)
	if totalRemotes > 0 {
		for i, r := range e.Remotes {
			_, shortURL := trimRemoteURL(r.URL[0])
			if r.Name == e.Remote.Name {
				currentindex = i
				fmt.Fprintln(out, selectionIndicator+r.Name+": "+shortURL)
				continue
			}
			fmt.Fprintln(out, tab+r.Name+": "+shortURL)
		}
		if err = gui.smartAnchorRelativeToLine(out, currentindex, totalRemotes); err != nil {
			return err
		}
	}
	return nil
}

// updates the remotebranchview for given entity
func (gui *Gui) renderRemoteBranches(e *git.RepoEntity) error {
	var err error
	out, err := gui.g.View(remoteBranchViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	trb := len(e.Remote.Branches)
	if trb > 0 {
		for i, r := range e.Remote.Branches {
			if r.Name == e.Remote.Branch.Name {
				currentindex = i
				fmt.Fprintln(out, selectionIndicator+r.Name)
				continue
			}
			fmt.Fprintln(out, tab+r.Name)
		}
		if err = gui.smartAnchorRelativeToLine(out, currentindex, trb); err != nil {
			return err
		}
	}
	return nil
}

// updates the branchview for given entity
func (gui *Gui) renderBranch(e *git.RepoEntity) error {
	var err error
	out, err := gui.g.View(branchViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	totalbranches := len(e.Branches)
	for i, b := range e.Branches {
		if b.Name == e.Branch.Name {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator+b.Name)
			continue
		}
		fmt.Fprintln(out, tab+b.Name)
	}
	err = gui.smartAnchorRelativeToLine(out, currentindex, totalbranches)
	return err
}

// updates the commitsview for given entity
func (gui *Gui) renderCommits(e *git.RepoEntity) error {
	var err error
	out, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	totalcommits := len(e.Commits)
	for i, c := range e.Commits {
		if c.Hash == e.Commit.Hash {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator+commitLabel(c))
			continue
		}
		fmt.Fprintln(out, tab+commitLabel(c))
	}
	if err = gui.smartAnchorRelativeToLine(out, currentindex, totalcommits); err != nil {
		return err
	}
	return err
}

// cursor down variant for sideviews
func (gui *Gui) sideViewsNextItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	e := gui.getSelectedRepository()
	switch viewName := v.Name(); viewName {
	case remoteBranchViewFeature.Name:
		return e.Remote.NextRemoteBranch(e)
	case remoteViewFeature.Name:
		return e.NextRemote()
	case branchViewFeature.Name:
		if err = e.Checkout(e.NextBranch()); err != nil {
			err = gui.openErrorView(g, err.Error(),
				"You should manually resolve this issue",
				branchViewFeature.Name)
			return err
		}
	case commitViewFeature.Name:
		e.NextCommit()
		return gui.renderCommits(e)
	}
	return err
}

// cursor up variant for sideviews
func (gui *Gui) sideViewsPreviousItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	e := gui.getSelectedRepository()
	switch viewName := v.Name(); viewName {
	case remoteBranchViewFeature.Name:
		return e.Remote.PreviousRemoteBranch(e)
	case remoteViewFeature.Name:
		return e.PreviousRemote()
	case branchViewFeature.Name:
		if err = e.Checkout(e.PreviousBranch()); err != nil {
			err = gui.openErrorView(g, err.Error(),
				"You should manually resolve this issue",
				branchViewFeature.Name)
			return err
		}
	case commitViewFeature.Name:
		e.PreviousCommit()
		return gui.renderCommits(e)
	}
	return err
}

// basically does fetch --prune
func (gui *Gui) syncRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	e := gui.getSelectedRepository()
	return git.Fetch(e, git.FetchOptions{
		RemoteName: e.Remote.Name,
		Prune:      true,
	})
}

// opens a confirmation view for setting default merge branch
func (gui *Gui) setUpstreamToBranch(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	e := gui.getSelectedRepository()
	v, err := g.SetView(confirmationViewFeature.Name, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "branch."+e.Branch.Name+"."+"remote"+"="+e.Remote.Name)
		fmt.Fprintln(v, "branch."+e.Branch.Name+"."+"merge"+"="+e.Branch.Reference.Name().String())
	}
	return gui.focusToView(confirmationViewFeature.Name)
}

// add config for upstream merge
func (gui *Gui) confirmSetUpstreamToBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	e := gui.getSelectedRepository()
	if err = git.AddConfig(e, git.ConfigOptions{
		Section: "branch." + e.Branch.Name,
		Option:  "remote",
		Site:    git.ConfigSiteLocal,
	}, e.Remote.Name); err != nil {
		return err
	}
	if err = git.AddConfig(e, git.ConfigOptions{
		Section: "branch." + e.Branch.Name,
		Option:  "merge",
		Site:    git.ConfigSiteLocal,
	}, e.Branch.Reference.Name().String()); err != nil {
		return err
	}
	e.Refresh()
	return gui.closeConfirmationView(g, v)
}

// close confirmation view
func (gui *Gui) closeConfirmationView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	return gui.closeViewCleanup(branchViewFeature.Name)
}
