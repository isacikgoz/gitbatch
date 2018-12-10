package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

var (
	statusHeaderViewFeature = viewFeature{Name: "status-header", Title: " Status Header "}
	// statusViewFeature       = viewFeature{Name: "status", Title: " Status "}
	stageViewFeature   = viewFeature{Name: "staged", Title: " Staged "}
	unstageViewFeature = viewFeature{Name: "unstaged", Title: " Not Staged "}
	stashViewFeature   = viewFeature{Name: "stash", Title: " Stash "}

	statusViews = []viewFeature{stageViewFeature, unstageViewFeature, stashViewFeature}
)

// open the status layout
func (gui *Gui) openStatusView(g *gocui.Gui, v *gocui.View) error {
	gui.openStatusHeaderView(g)
	gui.openStageView(g)
	gui.openUnStagedView(g)
	gui.openStashView(g)
	return nil
}

// focus to next view
func (gui *Gui) nextStatusView(g *gocui.Gui, v *gocui.View) error {
	err := gui.nextViewOfGroup(g, v, statusViews)
	return err
}

// focus to previous view
func (gui *Gui) previousStatusView(g *gocui.Gui, v *gocui.View) error {
	err := gui.previousViewOfGroup(g, v, statusViews)
	return err
}

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) statusCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		ly := len(v.BufferLines()) - 2 // why magic number? have no idea

		// if we are at the end we just return
		if cy+oy == ly {
			return nil
		}
		if err := v.SetCursor(cx, cy+1); err != nil {

			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		entity := gui.getSelectedRepository()
		if err := refreshStatusView(v.Name(), g, entity); err != nil {
			return err
		}
	}
	return nil
}

// moves the cursor upwards for the main view
func (gui *Gui) statusCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		entity := gui.getSelectedRepository()
		if err := refreshStatusView(v.Name(), g, entity); err != nil {
			return err
		}
	}
	return nil
}

// header og the status layout
func (gui *Gui) openStatusHeaderView(g *gocui.Gui) error {
	maxX, _ := g.Size()
	entity := gui.getSelectedRepository()
	v, err := g.SetView(statusHeaderViewFeature.Name, 6, 2, maxX-6, 4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, entity.AbsPath)
		// v.Frame = false
		v.Wrap = true
	}
	return nil
}

// close the opened stat views
func (gui *Gui) closeStatusView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(stashViewFeature.Name); err != nil {
		return err
	}
	if err := g.DeleteView(unstageViewFeature.Name); err != nil {
		return err
	}
	if err := g.DeleteView(stageViewFeature.Name); err != nil {
		return err
	}
	if err := g.DeleteView(statusHeaderViewFeature.Name); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	entity := gui.getSelectedRepository()
	if err := gui.refreshMain(g); err != nil {
		return err
	}
	if err := gui.refreshViews(g, entity); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

func generateFileLists(entity *git.RepoEntity) (staged, unstaged []*git.File, err error) {
	files, err := entity.LoadFiles()
	if err != nil {
		return nil, nil, err
	}
	for _, file := range files {
		if file.X != git.StatusNotupdated && file.X != git.StatusUntracked && file.X != git.StatusIgnored && file.X != git.StatusUpdated {
			staged = append(staged, file)
		}
		if file.Y != git.StatusNotupdated {
			unstaged = append(unstaged, file)
		}
	}
	return staged, unstaged, err
}

func refreshStatusView(viewName string, g *gocui.Gui, entity *git.RepoEntity) error {
	switch viewName {
	case stageViewFeature.Name:
		if err := refreshStagedView(g, entity); err != nil {
			return err
		}
	case unstageViewFeature.Name:
		if err := refreshUnstagedView(g, entity); err != nil {
			return err
		}
	case stashViewFeature.Name:
		if err := refreshStashView(g, entity); err != nil {
			return err
		}
	}
	return nil
}

func refreshAllStatusView(g *gocui.Gui, entity *git.RepoEntity) error {
	for _, v := range statusViews {
		if err := refreshStatusView(v.Name, g, entity); err != nil {
			return err
		}
	}
	return nil
}
