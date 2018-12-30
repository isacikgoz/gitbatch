package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

var (
	statusHeaderViewFeature = viewFeature{Name: "status-header", Title: " Status Header "}
	stageViewFeature        = viewFeature{Name: "staged", Title: " Staged "}
	unstageViewFeature      = viewFeature{Name: "unstaged", Title: " Not Staged "}
	stashViewFeature        = viewFeature{Name: "stash", Title: " Stash "}

	statusViews = []viewFeature{stageViewFeature, unstageViewFeature, stashViewFeature}

	commitMesageReturnView string
	stagedFiles            []*git.File
	unstagedFiles          []*git.File
)

// open the status layout
func (gui *Gui) openStatusView(g *gocui.Gui, v *gocui.View) error {
	if err := reloadFiles(gui.getSelectedRepository()); err != nil {
		return err
	}
	gui.openStatusHeaderView(g)
	gui.openStageView(g)
	gui.openUnStagedView(g)
	gui.openStashView(g)
	return nil
}

func reloadFiles(e *git.RepoEntity) error {
	_, _, err := populateFileLists(e)
	return err
}

// focus to next view
func (gui *Gui) nextStatusView(g *gocui.Gui, v *gocui.View) error {
	return gui.nextViewOfGroup(g, v, statusViews)
}

// focus to previous view
func (gui *Gui) previousStatusView(g *gocui.Gui, v *gocui.View) error {
	return gui.previousViewOfGroup(g, v, statusViews)
}

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) statusCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

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
	e := gui.getSelectedRepository()
	return refreshStatusView(v.Name(), g, e, false)
}

// moves the cursor upwards for the main view
func (gui *Gui) statusCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	e := gui.getSelectedRepository()
	return refreshStatusView(v.Name(), g, e, false)
}

// header og the status layout
func (gui *Gui) openStatusHeaderView(g *gocui.Gui) error {
	maxX, _ := g.Size()
	e := gui.getSelectedRepository()
	v, err := g.SetView(statusHeaderViewFeature.Name, 6, 2, maxX-6, 4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, e.AbsPath)
		// v.Frame = false
		v.Wrap = true
	}
	return nil
}

// close the opened stat views
func (gui *Gui) closeStatusView(g *gocui.Gui, v *gocui.View) error {
	for _, view := range statusViews {
		if err := g.DeleteView(view.Name); err != nil {
			return err
		}
	}
	if err := g.DeleteView(statusHeaderViewFeature.Name); err != nil {
		return err
	}
	stagedFiles = make([]*git.File, 0)
	unstagedFiles = make([]*git.File, 0)

	return gui.closeViewCleanup(mainViewFeature.Name)
}

// generate file lists by git status command
func populateFileLists(e *git.RepoEntity) (staged, unstaged []*git.File, err error) {
	files, err := git.Status(e)
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

func refreshStatusView(viewName string, g *gocui.Gui, e *git.RepoEntity, reload bool) error {
	if reload {
		reloadFiles(e)
	}
	var err error
	switch viewName {
	case stageViewFeature.Name:
		err = refreshStagedView(g)
	case unstageViewFeature.Name:
		err = refreshUnstagedView(g)
	case stashViewFeature.Name:
		err = refreshStashView(g, e)
	}
	return err
}

func refreshAllStatusView(g *gocui.Gui, e *git.RepoEntity, reload bool) error {
	for _, v := range statusViews {
		if err := refreshStatusView(v.Name, g, e, reload); err != nil {
			return err
		}
	}
	return nil
}
