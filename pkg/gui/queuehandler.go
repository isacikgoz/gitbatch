package gui

import (
	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// this function starts the queue and updates the gui with the result of an
// operation
func (gui *Gui) startQueue(g *gocui.Gui, v *gocui.View) error {
	go func(gui_go *Gui, g_go *gocui.Gui) {
		for {
			job, finished, err := gui_go.State.Queue.StartNext()
			// for each job execution we better refresh the main
			// it would be nice if we can also refresh side views
			g_go.Update(func(gu *gocui.Gui) error {
				gui_go.refreshMain(gu)
				return nil
			})

			if err != nil {
				if err == git.ErrAuthenticationRequired {
					// pause the job, so it will be indicated to being blocking
					job.Entity.State = git.Paused
					err := gui_go.openAuthenticationView(g, gui_go.State.Queue, job, v.Name())
					if err != nil {
						log.Warn(err.Error())
						return
					}
				}
				return
				// with not returning here, we simply ignore and continue
			}
			// if queue is finished simply return from this goroutine
			if finished {
				return
			}
			selectedEntity := gui_go.getSelectedRepository()
			if job.Entity == selectedEntity {
				gui_go.refreshViews(g, job.Entity)
			}
		}
	}(gui, g)
	return nil
}

// flashes the keybinding view's backgroun with green color to indicate that
// the queue is started
func indicateQueueStarted(g *gocui.Gui) error {
	v, err := g.View(keybindingsViewFeature.Name)
	if err != nil {
		return err
	}
	v.BgColor = gocui.ColorGreen
	v.FgColor = gocui.ColorBlack
	return nil
}
