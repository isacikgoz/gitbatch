package gui

import (
	"errors"
	"fmt"
	// "strconv"

	"github.com/jroimartin/gocui"
)

func (gui *Gui) openExecConfirmationView(g *gocui.Gui, v *gocui.View) error {
	// maxX, maxY := g.Size()
	go func(gui_go *Gui, g_go *gocui.Gui) {
		for {
			finished, err := gui_go.State.Queue.StartNext()
			g.Update(func(gu *gocui.Gui) error {
				gui_go.refreshMain(gu)
				return nil
			})
			if err != nil {
				return
			}
			if finished {
				return
			} 
		}
	}(gui, g)
	return nil
	// if mrs, _ := gui.getMarkedEntities(); len(mrs) < 1 {
	// 	return nil
	// }

	// v, err := g.SetView(execViewFeature.Name, maxX/2-35, maxY/2-5, maxX/2+35, maxY/2+5)
	// if err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	v.Title = gui.State.Mode.DisplayString + " Confirmation"
	// 	v.Wrap = true
	// 	mrs, _ := gui.getMarkedEntities()
	// 	jobs := strconv.Itoa(len(mrs)) + " " + gui.State.Mode.ExecString + ":"
	// 	fmt.Fprintln(v, jobs)
	// 	for _, r := range mrs {
	// 		line := " - " + green.Sprint(r.Name) + ": " + r.Remote.Name + green.Sprint(" → ") + r.Branch.Name
	// 		fmt.Fprintln(v, line)
	// 	}
	// 	ps := red.Sprint("Note:") + " When " + gui.State.Mode.CommandString + " operation is completed, you will be notified."
	// 	fmt.Fprintln(v, "\n"+ps)
	// }
	// gui.updateKeyBindingsView(g, execViewFeature.Name)
	// if _, err := g.SetCurrentView(execViewFeature.Name); err != nil {
	// 	return err
	// }
	return nil
}

func (gui *Gui) closeExecView(g *gocui.Gui, v *gocui.View) error {
	go g.Update(func(g *gocui.Gui) error {
		mainView, _ := g.View(mainViewFeature.Name)
		entity, err := gui.getSelectedRepository(g, mainView)
		if err != nil {
			return err
		}
		gui.updateCommits(g, entity)
		return nil
	})
	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

func (gui *Gui) execute(g *gocui.Gui, v *gocui.View) error {
	// somehow this fucntion called after this method returns, strange?
	go g.Update(func(g *gocui.Gui) error {
		err := updateKeyBindingsViewForExecution(g)
		if err != nil {
			return err
		}
		return nil
	})

	mrs, _ := gui.getMarkedEntities()
	for _, mr := range mrs {
		// here we will be waiting
		switch mode := gui.State.Mode.ModeID; mode {
		case FetchMode:
			err := mr.Fetch()
			if err != nil {
				cv, _ := g.View(execViewFeature.Name)
				gui.closeExecView(g, cv)
				gui.openErrorView(g, err.Error(), "An error occured, manual resolving reuqired")
				return nil
			}
		case PullMode:
			err := mr.Pull()
			if err != nil {
				cv, _ := g.View(execViewFeature.Name)
				gui.closeExecView(g, cv)
				gui.openErrorView(g, err.Error(), "It maybe a conflict, manual resolving reuqired")
				return nil
			}
		default:
			return errors.New("No mode is selected")
		}
		mr.Unmark()
		gui.refreshMain(g)
	}
	return nil
}

func updateKeyBindingsViewForExecution(g *gocui.Gui) error {
	v, err := g.View(keybindingsViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	v.BgColor = gocui.ColorGreen
	v.FgColor = gocui.ColorBlack
	v.Frame = false
	fmt.Fprintln(v, " Operation Completed; c: close/cancel")
	return nil
}
