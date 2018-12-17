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

// updates the remotesview for given entity
func (gui *Gui) updateRemotes(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	out, err := g.View(remoteViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()

	currentindex := 0
	totalRemotes := len(entity.Remotes)
	if totalRemotes > 0 {
		for i, r := range entity.Remotes {
			// TODO: maybe the text styling can be moved to textstyle.go file
			_, shortURL := trimRemoteURL(r.URL[0])
			suffix := shortURL
			if r.Name == entity.Remote.Name {
				currentindex = i
				fmt.Fprintln(out, selectionIndicator+r.Name+": "+suffix)
				continue
			}
			fmt.Fprintln(out, tab+r.Name+": "+suffix)
		}
		if err = gui.smartAnchorRelativeToLine(out, currentindex, totalRemotes); err != nil {
			return err
		}
	}
	return nil
}

// updates the remotebranchview for given entity
func (gui *Gui) updateRemoteBranches(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	out, err := g.View(remoteBranchViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()
	currentindex := 0
	trb := len(entity.Remote.Branches)
	if trb > 0 {
		for i, r := range entity.Remote.Branches {
			rName := r.Name
			if r.Deleted {
				rName = rName + ws + dirty
			}
			if r.Name == entity.Remote.Branch.Name {
				currentindex = i
				fmt.Fprintln(out, selectionIndicator+rName)
				continue
			}
			fmt.Fprintln(out, tab+rName)
		}
		if err = gui.smartAnchorRelativeToLine(out, currentindex, trb); err != nil {
			return err
		}
	}
	return nil
}

// updates the branchview for given entity
func (gui *Gui) updateBranch(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	out, err := g.View(branchViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()

	currentindex := 0
	totalbranches := len(entity.Branches)
	for i, b := range entity.Branches {
		if b.Name == entity.Branch.Name {
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
func (gui *Gui) updateCommits(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	out, err := g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()

	currentindex := 0
	totalcommits := len(entity.Commits)
	for i, c := range entity.Commits {
		var body string
		if c.CommitType == git.EvenCommit {
			body = cyan.Sprint(c.Hash[:hashLength]) + " " + c.Message
		} else if c.CommitType == git.LocalCommit {
			body = blue.Sprint(c.Hash[:hashLength]) + " " + c.Message
		} else {
			body = yellow.Sprint(c.Hash[:hashLength]) + " " + c.Message
		}
		if c.Hash == entity.Commit.Hash {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator+body)
			continue
		}
		fmt.Fprintln(out, tab+body)
	}
	if err = gui.smartAnchorRelativeToLine(out, currentindex, totalcommits); err != nil {
		return err
	}
	return err
}

func (gui *Gui) sideViewsNextItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	switch viewName := v.Name(); viewName {
	case remoteBranchViewFeature.Name:
		if err = entity.Remote.NextRemoteBranch(); err != nil {
			return err
		}
		err = gui.updateRemoteBranches(g, entity)
	case remoteViewFeature.Name:
		if err = entity.NextRemote(); err != nil {
			return err
		}
		err = gui.remoteChangeFollowUp(g, entity)
	case branchViewFeature.Name:
		if err = entity.Checkout(entity.NextBranch()); err != nil {
			err = gui.openErrorView(g, err.Error(),
				"You should manually resolve this issue",
				branchViewFeature.Name)
			return err
		}
		err = gui.checkoutFollowUp(g, entity)
	case commitViewFeature.Name:
		if err = entity.NextCommit(); err != nil {
			return err
		}
		err = gui.updateCommits(g, entity)
	}
	return err
}

func (gui *Gui) sideViewsPreviousItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	switch viewName := v.Name(); viewName {
	case remoteBranchViewFeature.Name:
		if err = entity.Remote.PreviousRemoteBranch(); err != nil {
			return err
		}
		err = gui.updateRemoteBranches(g, entity)
	case remoteViewFeature.Name:
		if err = entity.PreviousRemote(); err != nil {
			return err
		}
		err = gui.remoteChangeFollowUp(g, entity)
	case branchViewFeature.Name:
		if err = entity.Checkout(entity.PreviousBranch()); err != nil {
			err = gui.openErrorView(g, err.Error(),
				"You should manually resolve this issue",
				branchViewFeature.Name)
			return err
		}
		err = gui.checkoutFollowUp(g, entity)
	case commitViewFeature.Name:
		if err = entity.PreviousCommit(); err != nil {
			return err
		}
		err = gui.updateCommits(g, entity)
	}
	return err
}

// basically does fetch --prune
func (gui *Gui) syncRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = git.Fetch(entity, git.FetchOptions{
		RemoteName: entity.Remote.Name,
		Prune:      true,
	}); err != nil {
		return err
	}
	vr, err := g.View(remoteViewFeature.Name)
	if err != nil {
		return err
	}
	// have no idea why this works..
	// some time need to fix, movement aint bad huh?
	gui.sideViewsNextItem(g, vr)
	gui.sideViewsPreviousItem(g, vr)
	err = gui.updateRemoteBranches(g, entity)
	return err
}

// basically does fetch --prune
func (gui *Gui) setUpstreamToBranch(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	entity := gui.getSelectedRepository()
	v, err := g.SetView(confirmationViewFeature.Name, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "branch."+entity.Branch.Name+"."+"remote"+"="+entity.Remote.Name)
		fmt.Fprintln(v, "branch."+entity.Branch.Name+"."+"merge"+"="+entity.Branch.Reference.Name().String())
	}
	return gui.focusToView(confirmationViewFeature.Name)
}

// basically does fetch --prune
func (gui *Gui) confirmSetUpstreamToBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = git.AddConfig(entity, git.ConfigOptions{
		Section: "branch." + entity.Branch.Name,
		Option:  "remote",
		Site:    git.ConfigSiteLocal,
	}, entity.Remote.Name); err != nil {
		return err
	}
	if err = git.AddConfig(entity, git.ConfigOptions{
		Section: "branch." + entity.Branch.Name,
		Option:  "merge",
		Site:    git.ConfigSiteLocal,
	}, entity.Branch.Reference.Name().String()); err != nil {
		return err
	}
	entity.Refresh()
	gui.refreshMain(g)
	return gui.closeConfirmationView(g, v)
}

func (gui *Gui) closeConfirmationView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	return gui.closeViewCleanup(branchViewFeature.Name)
}

// after checkout a remote some refreshments needed
func (gui *Gui) remoteChangeFollowUp(g *gocui.Gui, entity *git.RepoEntity) (err error) {
	if err = gui.updateRemotes(g, entity); err != nil {
		return err
	}
	err = gui.updateRemoteBranches(g, entity)
	return err
}

// after checkout a branch some refreshments needed
func (gui *Gui) checkoutFollowUp(g *gocui.Gui, entity *git.RepoEntity) (err error) {
	if err = gui.updateBranch(g, entity); err != nil {
		return err
	}
	if err = gui.updateCommits(g, entity); err != nil {
		return err
	}
	if err = gui.updateRemoteBranches(g, entity); err != nil {
		return err
	}
	err = gui.refreshMain(g)
	return err
}
