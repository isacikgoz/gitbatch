package gui

import (
	"github.com/jroimartin/gocui"
)

// this function starts the queue and updates the gui with the result of an
// operation
func (gui *Gui) startQueue(g *gocui.Gui, v *gocui.View) error {
	go func(gui_go *Gui, g_go *gocui.Gui) {
		indicateQueueStarted(g_go)
		for {
			job, finished, err := gui_go.State.Queue.StartNext()
			g_go.Update(func(gu *gocui.Gui) error {
				gui_go.refreshMain(gu)
				return nil
			})
			defer gui.updateKeyBindingsView(g, mainViewFeature.Name)
			if err != nil {
				return
			}
			if finished {
				return
			} else {
				selectedEntity, _ := gui_go.getSelectedRepository(g, v)
				if job.Entity == selectedEntity {
					gui_go.refreshViews(g, job.Entity)
				}
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
