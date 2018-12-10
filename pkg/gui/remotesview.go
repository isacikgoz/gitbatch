package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
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

// iteration handler for the remotesview
func (gui *Gui) nextRemote(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = entity.NextRemote(); err != nil {
		return err
	}
	if err = gui.remoteChangeFollowUp(g, entity); err != nil {
		return err
	}
	return err
}

// iteration handler for the remotesview
func (gui *Gui) previousRemote(g *gocui.Gui, v *gocui.View) error {
	var err error
	entity := gui.getSelectedRepository()
	if err = entity.PreviousRemote(); err != nil {
		return err
	}
	if err = gui.remoteChangeFollowUp(g, entity); err != nil {
		return err
	}
	return err
}

// after checkout a remote some refreshments needed
func (gui *Gui) remoteChangeFollowUp(g *gocui.Gui, entity *git.RepoEntity) (err error) {
	if err = gui.updateRemotes(g, entity); err != nil {
		return err
	}
	err = gui.updateRemoteBranches(g, entity)
	return err
}
