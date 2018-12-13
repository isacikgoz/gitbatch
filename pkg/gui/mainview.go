package gui

import (
	"fmt"
	"sort"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

// this is the initial function for filling the values for the main view. the
// function waits a separate routine to fill the gui's repository slice
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
		if err = gui.sortByName(g, v); err != nil {
			return err
		}
		entity := gui.getSelectedRepository()
		gui.refreshViews(g, entity)
		return nil
	})
	return nil
}

// refresh the main view and re-render the repository representations
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

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
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
		entity := gui.getSelectedRepository()
		if err := gui.refreshMain(g); err != nil {
			return err
		}
		gui.refreshViews(g, entity)
	}
	return nil
}

// moves the cursor upwards for the main view
func (gui *Gui) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		entity := gui.getSelectedRepository()
		if err := gui.refreshMain(g); err != nil {
			return err
		}
		gui.refreshViews(g, entity)
	}
	return nil
}

// returns the entity at cursors position by taking its position in the gui's
// slice of repositories. Since it is not a %100 percent safe methodology it may
// rrequire a better implementation or the slice's order must be synchronized
// with the views lines
func (gui *Gui) getSelectedRepository() *git.RepoEntity {
	v, _ := gui.g.View(mainViewFeature.Name)
	_, oy := v.Origin()
	_, cy := v.Cursor()
	// if _, err := v.Line(cy); err != nil {
	// 	return r, err
	// }
	return gui.State.Repositories[cy+oy]
}

// adds given entity to job queue
func (gui *Gui) addToQueue(entity *git.RepoEntity) error {
	var jt git.JobType
	switch mode := gui.State.Mode.ModeID; mode {
	case FetchMode:
		jt = git.FetchJob
	case PullMode:
		jt = git.PullJob
	case MergeMode:
		jt = git.MergeJob
	default:
		return nil
	}
	err := gui.State.Queue.AddJob(&git.Job{
		JobType: jt,
		Entity:  entity,
	})
	if err != nil {
		return err
	}
	entity.State = git.Queued
	return nil
}

// removes given entity from job queue
func (gui *Gui) removeFromQueue(entity *git.RepoEntity) error {
	err := gui.State.Queue.RemoveFromQueue(entity)
	if err != nil {
		return err
	}
	entity.State = git.Available
	return nil
}

// marking repository is simply adding the repostirory into the queue. the
// function does take its current state into account before adding it
func (gui *Gui) markRepository(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	// maybe, failed entities may be added to queue again
	if r.State == git.Available || r.State == git.Success || r.State == git.Paused {
		if err := gui.addToQueue(r); err != nil {
			return err
		}
	} else if r.State == git.Queued {
		if err := gui.removeFromQueue(r); err != nil {
			return err
		}
	} else {
		return nil
	}
	gui.refreshMain(g)
	return nil
}

// add all remaining repositories into the queue. the function does take its
// current state into account before adding it
func (gui *Gui) markAllRepositories(g *gocui.Gui, v *gocui.View) error {
	for _, r := range gui.State.Repositories {
		if r.State == git.Available || r.State == git.Success {
			if err := gui.addToQueue(r); err != nil {
				return err
			}
		} else {
			continue
		}
	}
	gui.refreshMain(g)
	return nil
}

// remove all repositories from the queue. the function does take its
// current state into account before removing it
func (gui *Gui) unmarkAllRepositories(g *gocui.Gui, v *gocui.View) error {
	for _, r := range gui.State.Repositories {
		if r.State == git.Queued {
			if err := gui.removeFromQueue(r); err != nil {
				return err
			}
		} else {
			continue
		}
	}
	gui.refreshMain(g)
	return nil
}

// sortByName sorts the repositories by A to Z order
func (gui *Gui) sortByName(g *gocui.Gui, v *gocui.View) error {
	sort.Sort(git.Alphabetical(gui.State.Repositories))
	gui.refreshAfterSort(g)
	return nil
}

// sortByMod sorts the repositories according to last modifed date
// the top element will be the last modified
func (gui *Gui) sortByMod(g *gocui.Gui, v *gocui.View) error {
	sort.Sort(git.LastModified(gui.State.Repositories))
	gui.refreshAfterSort(g)
	return nil
}

// utility function that refreshes main and side views after that
func (gui *Gui) refreshAfterSort(g *gocui.Gui) error {
	gui.refreshMain(g)
	entity := gui.getSelectedRepository()
	gui.refreshViews(g, entity)
	return nil
}
