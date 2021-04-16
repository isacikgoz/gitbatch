package gui

import (
	"fmt"
	"sort"

	"github.com/isacikgoz/gitbatch/internal/command"
	gerr "github.com/isacikgoz/gitbatch/internal/errors"
	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/isacikgoz/gitbatch/internal/job"
	"github.com/jroimartin/gocui"
)

// refresh the main view and re-render the repository representations
func (gui *Gui) renderMain() error {
	gui.mutex.Lock()
	defer gui.mutex.Unlock()

	mainView, err := gui.g.View(mainViewFeature.Name)
	if err != nil {
		return err
	}
	mainView.Clear()
	for _, r := range gui.State.Repositories {
		fmt.Fprintln(mainView, gui.repositoryLabel(r))
	}
	// while refreshing, refresh sideViews for selected entity, something may
	// be changed?
	return gui.renderSideViews(gui.getSelectedRepository())
}

// listens the event -> "repository.updated"
func (gui *Gui) repositoryUpdated(event *git.RepositoryEvent) error {
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.renderMain()
	})
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
	}
	return gui.renderMain()
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
	}
	return gui.renderMain()
}

// moves cursor to the top
func (gui *Gui) cursorTop(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, _ := v.Origin()
		cx, _ := v.Cursor()
		if err := v.SetOrigin(ox, 0); err != nil {
			return err
		}
		if err := v.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	return gui.renderMain()
}

// moves cursor to the end
func (gui *Gui) cursorEnd(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, _ := v.Origin()
		cx, _ := v.Cursor()
		_, vy := v.Size()
		lr := len(gui.State.Repositories)
		if lr <= vy {
			if err := v.SetCursor(cx, lr-1); err != nil {
				return err
			}
			return gui.renderMain()
		}
		if err := v.SetOrigin(ox, lr-vy); err != nil {
			return err
		}
		if err := v.SetCursor(cx, vy-1); err != nil {
			return err
		}
	}
	return gui.renderMain()
}

// moves cursor down for a page size
func (gui *Gui) pageDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, _ := v.Cursor()
		_, vy := v.Size()
		lr := len(gui.State.Repositories)
		if lr < vy {
			return nil
		}
		if oy+vy >= lr-vy {
			if err := v.SetOrigin(ox, lr-vy); err != nil {
				return err
			}
		} else if err := v.SetOrigin(ox, oy+vy); err != nil {
			return err
		}
		if err := v.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	return gui.renderMain()
}

// moves cursor up for a page size
func (gui *Gui) pageUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		_, vy := v.Size()
		if oy == 0 || oy+cy < vy {
			if err := v.SetOrigin(ox, 0); err != nil {
				return err
			}
		} else if oy <= vy {
			if err := v.SetOrigin(ox, oy+cy-vy); err != nil {
				return err
			}
		} else if err := v.SetOrigin(ox, oy-vy); err != nil {
			return err
		}
		if err := v.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	return gui.renderMain()
}

// returns the entity at cursors position by taking its position in the gui's
// slice of repositories. Since it is not a %100 percent safe methodology it may
// require a better implementation or the slice's order must be synchronized
// with the views lines
func (gui *Gui) getSelectedRepository() *git.Repository {
	if len(gui.State.Repositories) == 0 {
		return nil
	}
	v, _ := gui.g.View(mainViewFeature.Name)
	_, oy := v.Origin()
	_, cy := v.Cursor()
	return gui.State.Repositories[cy+oy]
}

// adds given entity to job queue
func (gui *Gui) addToQueue(r *git.Repository) error {
	j := &job.Job{
		Repository: r,
	}
	switch mode := gui.State.Mode.ModeID; mode {
	case FetchMode:
		j.JobType = job.FetchJob
	case PullMode:
		if r.State.Branch.Upstream == nil {
			return nil
		}
		j.JobType = job.PullJob
	case MergeMode:
		if r.State.Branch.Upstream == nil {
			return nil
		}
		j.JobType = job.MergeJob
	case CheckoutMode:
		j.JobType = job.CheckoutJob
		j.Options = &command.CheckoutOptions{
			TargetRef:      gui.State.targetBranch,
			CreateIfAbsent: true,
		}
	default:
		return nil
	}
	err := gui.State.Queue.AddJob(j)
	if err != nil {
		return err
	}
	r.SetWorkStatus(git.Queued)
	return nil
}

// removes given entity from job queue
func (gui *Gui) removeFromQueue(r *git.Repository) error {
	err := gui.State.Queue.RemoveFromQueue(r)
	if err != nil {
		return err
	}
	r.SetWorkStatus(git.Available)
	return nil
}

// this function starts the queue and updates the gui with the result of an
// operation
func (gui *Gui) startQueue(g *gocui.Gui, v *gocui.View) error {
	go func(gui_go *Gui) {
		fails := gui_go.State.Queue.StartJobsAsync()
		gui_go.State.Queue = job.CreateJobQueue()
		for j, err := range fails {
			if err == gerr.ErrAuthenticationRequired {
				j.Repository.SetWorkStatus(git.Paused)
				_ = gui_go.State.FailoverQueue.AddJob(j)
			}
		}
	}(gui)
	return nil
}

func (gui *Gui) submitCredentials(g *gocui.Gui, v *gocui.View) error {
	if is, j := gui.State.FailoverQueue.IsInTheQueue(gui.getSelectedRepository()); is {
		if j.Repository.WorkStatus() == git.Paused {
			if err := gui.State.FailoverQueue.RemoveFromQueue(j.Repository); err != nil {
				return err
			}
			err := gui.openAuthenticationView(g, gui.State.Queue, j, v.Name())
			if err != nil {
				return err
			}
			if isnt, _ := gui.State.Queue.IsInTheQueue(j.Repository); !isnt {
				_ = gui.State.FailoverQueue.AddJob(j)
			}
		}
	}
	return nil
}

// marking repository is simply adding the repository into the queue. the
// function does take its current state into account before adding it
func (gui *Gui) markRepository(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	// maybe, failed entities may be added to queue again
	if r == nil {
		return nil
	}
	if r.WorkStatus().Ready {
		if err := gui.addToQueue(r); err != nil {
			return err
		}
	} else if r.WorkStatus() == git.Queued {
		if err := gui.removeFromQueue(r); err != nil {
			return err
		}
	}
	return nil
}

// add all remaining repositories into the queue. the function does take its
// current state into account before adding it
func (gui *Gui) markAllRepositories(g *gocui.Gui, v *gocui.View) error {
	for _, r := range gui.State.Repositories {
		if r.WorkStatus().Ready {
			if err := gui.addToQueue(r); err != nil {
				return err
			}
		} else {
			continue
		}
	}
	return nil
}

// remove all repositories from the queue. the function does take its
// current state into account before removing it
func (gui *Gui) unmarkAllRepositories(g *gocui.Gui, v *gocui.View) error {
	for _, r := range gui.State.Repositories {
		if r.WorkStatus() == git.Queued {
			if err := gui.removeFromQueue(r); err != nil {
				return err
			}
		} else {
			continue
		}
	}
	return nil
}

// sortByName sorts the repositories by A to Z order
func (gui *Gui) sortByName(g *gocui.Gui, v *gocui.View) error {
	sort.Sort(git.Alphabetical(gui.State.Repositories))
	_ = gui.renderMain()
	return nil
}

// sortByMod sorts the repositories according to last modifed date
// the top element will be the last modified
func (gui *Gui) sortByMod(g *gocui.Gui, v *gocui.View) error {
	sort.Sort(git.LastModified(gui.State.Repositories))
	_ = gui.renderMain()
	return nil
}
