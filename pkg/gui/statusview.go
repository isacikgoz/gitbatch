package gui

import (
	"fmt"
	"time"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
	ggt "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	statusHeaderViewFeature  = viewFeature{Name: "status-header", Title: " Status Header "}
	stageViewFeature         = viewFeature{Name: "staged", Title: " Staged "}
	unstageViewFeature       = viewFeature{Name: "unstaged", Title: " Not Staged "}
	stashViewFeature         = viewFeature{Name: "stash", Title: " Stash "}
	commitMessageViewFeature = viewFeature{Name: "commitmessage", Title: " Commit Mesage "}

	statusViews            = []viewFeature{stageViewFeature, unstageViewFeature, stashViewFeature}
	commitMesageReturnView string
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
	for _, view := range statusViews {
		if err := g.DeleteView(view.Name); err != nil {
			return err
		}
	}
	if err := g.DeleteView(statusHeaderViewFeature.Name); err != nil {
		return err
	}
	entity := gui.getSelectedRepository()
	if err := gui.refreshMain(g); err != nil {
		return err
	}
	if err := gui.refreshViews(g, entity); err != nil {
		return err
	}
	return gui.closeViewCleanup(mainViewFeature.Name)
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

// open the commit message views
func (gui *Gui) openCommitMessageView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	commitMesageReturnView = v.Name()
	v, err := g.SetView(commitMessageViewFeature.Name, maxX/2-30, maxY/2-3, maxX/2+30, maxY/2+3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = commitMessageViewFeature.Title
		v.Wrap = true
		v.Editable = true
		v.Editor = gocui.DefaultEditor
		v.Highlight = true
		g.Cursor = true
	}
	gui.updateKeyBindingsView(g, commitMessageViewFeature.Name)
	if _, err := g.SetCurrentView(commitMessageViewFeature.Name); err != nil {
		return err
	}
	return nil
}

// close the opened commite mesage view
func (gui *Gui) submitCommitMessageView(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()
	w, err := entity.Repository.Worktree()
	if err != nil {
		return err
	}
	// WIP: This better be removed to git pkg
	// TODO: read config and get name & e-mail
	_, err = w.Commit(v.ViewBuffer(), &ggt.CommitOptions{
		Author: &object.Signature{
			Name:  "İbrahim Serdar Açıkgöz",
			Email: "serdaracikgoz86@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	entity.Refresh()
	err = gui.closeCommitMessageView(g, v)
	return err
}

// close the opened commite mesage view
func (gui *Gui) closeCommitMessageView(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()
	g.Cursor = false
	if err := g.DeleteView(commitMessageViewFeature.Name); err != nil {
		return err
	}
	if err := gui.refreshMain(g); err != nil {
		return err
	}
	if err := gui.refreshViews(g, entity); err != nil {
		return err
	}
	if err := refreshAllStatusView(g, entity); err != nil {
		return err
	}
	return gui.closeViewCleanup(commitMesageReturnView)
}
