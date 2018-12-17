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

			if err != nil {
				if err == git.ErrAuthenticationRequired {
					// pause the job, so it will be indicated to being blocking
					job.Entity.SetState(git.Paused)
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
		}
	}(gui, g)
	return nil
}
