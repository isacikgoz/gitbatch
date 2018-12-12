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
			g_go.Update(func(gu *gocui.Gui) error {
				gui_go.refreshMain(gu)
				return nil
			})
			defer gui.updateKeyBindingsView(g, mainViewFeature.Name)
			if err != nil {
				if err == git.ErrAuthenticationRequired {
					err := gui_go.openAuthenticationView(g, gui_go.State.Queue, job, v.Name())
					if err != nil {
						log.Warn(err.Error())
						return
					}
				}
				return
			}
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
