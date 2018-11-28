package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/isacikgoz/gitbatch/pkg/job"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) fillMain(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(mainViewFeature.Name)
		if err != nil {
			return err
		}
		for _, r := range gui.State.Repositories {
			fmt.Fprintln(v, gui.displayString(r))
		}
		err = g.DeleteView(loadingViewFeature.Name)
		if err != nil {
			return err
		}
		if _, err = gui.setCurrentViewOnTop(g, mainViewFeature.Name); err != nil {
			return err
		}
		if entity, err := gui.getSelectedRepository(g, v); err != nil {
			return err
		} else {
			gui.refreshViews(g, entity)
		}
		return nil
	})
	return nil
}

func (gui *Gui) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		ly := len(gui.State.Repositories) - 1

		// if we are at the end we just return
		if cy+oy == ly {
			return nil
		}
		if err := v.SetCursor(cx, cy+1); err != nil {

			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		if entity, err := gui.getSelectedRepository(g, v); err != nil {
			return err
		} else {
			gui.refreshViews(g, entity)
		}
	}
	return nil
}

func (gui *Gui) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		if entity, err := gui.getSelectedRepository(g, v); err != nil {
			return err
		} else {
			gui.refreshViews(g, entity)
		}
	}
	return nil
}

func (gui *Gui) getSelectedRepository(g *gocui.Gui, v *gocui.View) (*git.RepoEntity, error) {
	var r *git.RepoEntity

	_, cy := v.Cursor()
	if _, err := v.Line(cy); err != nil {
		return r, err
	}
	return gui.State.Repositories[cy], nil
}

func (gui *Gui) markRepository(g *gocui.Gui, v *gocui.View) error {
	if r, err := gui.getSelectedRepository(g, v); err != nil {
		return err
	} else {
		if err != nil {
			return err
		}
		if r.State == git.Available || r.State == git.Success {
			var jt job.JobType
			switch mode := gui.State.Mode.ModeID; mode {
			case FetchMode:
				jt = job.Fetch
			case PullMode:
				jt = job.Pull
			case MergeMode:
				jt = job.Merge
			default:
				return nil
			}
			err := gui.State.Queue.AddJob(&job.Job{
				JobType: jt,
				Entity:  r,
			})
			if err != nil {
				return err
			}
			r.State = git.Queued
		} else if r.State == git.Queued {
			err := gui.State.Queue.RemoveFromQueue(r)
			if err != nil {
				return err
			}
			r.State = git.Available
		} else {
			return nil
		}
		gui.refreshMain(g)
	}
	return nil
}

func (gui *Gui) refreshMain(g *gocui.Gui) error {

	mainView, err := g.View(mainViewFeature.Name)
	if err != nil {
		return err
	}
	mainView.Clear()
	for _, r := range gui.State.Repositories {
		fmt.Fprintln(mainView, gui.displayString(r))
	}
	return nil
}
