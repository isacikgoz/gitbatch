package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

// not staged view
func (gui *Gui) openUnStagedView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(unstageViewFeature.Name, maxX/2+1, 5, maxX-6, int(0.75*float32(maxY))-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = unstageViewFeature.Title
	}
	entity := gui.getSelectedRepository()
	if err := refreshUnstagedView(g, entity); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) addChanges(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()
	_, files, err := generateFileLists(entity)
	if err != nil {
		return err
	}
	if len(files) <= 0 {
		return nil
	}
	_, cy := v.Cursor()
	_, oy := v.Origin()
	if err := files[cy+oy].Add(git.AddOptions{}); err != nil {
		return err
	}
	if err := refreshAllStatusView(g, entity); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) addAllChanges(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()
	if err := entity.AddAll(git.AddOptions{}); err != nil {
		return err
	}
	if err := refreshAllStatusView(g, entity); err != nil {
		return err
	}
	return nil
}

// refresh the main view and re-render the repository representations
func refreshUnstagedView(g *gocui.Gui, entity *git.RepoEntity) error {
	stageView, err := g.View(unstageViewFeature.Name)
	if err != nil {
		return err
	}
	stageView.Clear()
	_, cy := stageView.Cursor()
	_, oy := stageView.Origin()
	_, files, err := generateFileLists(entity)
	if err != nil {
		return err
	}
	for i, file := range files {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stageView, "%s%s%s %s\n", prefix, red.Sprint(string(file.X)), red.Sprint(string(file.Y)), file.Name)
	}
	return nil
}
