package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/core/command"
	"github.com/jroimartin/gocui"
)

// staged view
func (gui *Gui) openStageView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(stageViewFeature.Name, 6, 5, maxX/2-1, int(0.75*float32(maxY))-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stageViewFeature.Title
	}
	if err := refreshStagedView(g); err != nil {
		return err
	}
	return gui.focusToView(stageViewFeature.Name)
}

func (gui *Gui) resetChanges(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()

	_, cy := v.Cursor()
	_, oy := v.Origin()
	if len(stagedFiles) <= 0 || len(stagedFiles) <= cy+oy {
		return nil
	}
	if err := command.Reset(r, stagedFiles[cy+oy], &command.ResetOptions{}); err != nil {
		return err
	}
	return refreshAllStatusView(g, r, true)
}

func (gui *Gui) resetAllChanges(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	ref, err := r.Repo.Head()
	if err != nil {
		return err
	}
	if err := command.ResetAll(r, &command.ResetOptions{
		Hash:  ref.Hash().String(),
		Rtype: command.ResetMixed,
	}); err != nil {
		return err
	}
	return refreshAllStatusView(g, r, true)
}

// refresh the main view and re-render the repository representations
func refreshStagedView(g *gocui.Gui) error {
	stageView, err := g.View(stageViewFeature.Name)
	if err != nil {
		return err
	}
	stageView.Clear()
	_, cy := stageView.Cursor()
	_, oy := stageView.Origin()
	for i, file := range stagedFiles {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stageView, "%s%s%s %s\n", prefix, green.Sprint(string(file.X)), red.Sprint(string(file.Y)), file.Name)
	}
	return nil
}
