package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

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

// iteration handler for the remotebranchview
func (gui *Gui) syncRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = git.Fetch(entity, git.FetchOptions{
		RemoteName: entity.Remote.Name,
		Prune:      true,
	}); err != nil {
		return err
	}
	// have no idea why this works..
	// some time need to fix, movement aint bad huh?
	gui.nextRemote(g, v)
	gui.previousRemote(g, v)
	err = gui.updateRemoteBranches(g, entity)
	return err
}

// iteration handler for the remotebranchview
func (gui *Gui) nextRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = entity.Remote.NextRemoteBranch(); err != nil {
		return err
	}
	err = gui.updateRemoteBranches(g, entity)
	return err
}

// iteration handler for the remotebranchview
func (gui *Gui) previousRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = entity.Remote.PreviousRemoteBranch(); err != nil {
		return err
	}
	err = gui.updateRemoteBranches(g, entity)
	return err
}
