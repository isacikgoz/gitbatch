package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
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
	entity := gui.getSelectedRepository()
	if err := refreshStagedView(g, entity); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, stageViewFeature.Name)
	if _, err := g.SetCurrentView(stageViewFeature.Name); err != nil {
		return err
	}
	return nil
}

func  (gui *Gui) resetChanges(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()
	files, _, err := generateFileLists(entity)
	if err != nil {
		return err
	}
	if len(files) <= 0 {
		return nil
	}
	_, cy := v.Cursor()
	_, oy := v.Origin()
	if err := files[cy+oy].Reset(git.ResetOptions{

	}); err != nil {
		return err
	}
	if err := refreshAllStatusView(g, entity); err != nil {
		return err
	}
	return nil
}

func  (gui *Gui) resetAllChanges(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()
	if err := entity.ResetAll(git.ResetOptions{
		
	}); err != nil {
		return err
	}
	if err := refreshAllStatusView(g, entity); err != nil {
		return err
	}
	return nil
}

// refresh the main view and re-render the repository representations
func refreshStagedView(g *gocui.Gui, entity *git.RepoEntity) error {
	stageView, err := g.View(stageViewFeature.Name)
	if err != nil {
		return err
	}
	stageView.Clear()
	_, cy := stageView.Cursor()
	_, oy := stageView.Origin()
	files, _, err := generateFileLists(entity)
	if err != nil {
		return err
	}
	for i, file := range files {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stageView, "%s%s%s %s\n", prefix, green.Sprint(string(file.X)), red.Sprint(string(file.Y)), file.Name)
	}
	return nil
}
