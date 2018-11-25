package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) openPullView(g *gocui.Gui, v *gocui.View) error {
    maxX, maxY := g.Size()

    v, err := g.SetView(pullViewFeature.Name, maxX/2-35, maxY/2-5, maxX/2+35, maxY/2+5)
    if err != nil {
            if err != gocui.ErrUnknownView {
                    return err
            }
            v.Title = pullViewFeature.Title
            v.Wrap = false
            mrs, _ := gui.getMarkedEntities()
            for _, r := range mrs {
                line := r.Name + " : " + r.GetActiveRemote() + "/" + r.Branch + " → " + r.GetActiveBranch()
                fmt.Fprintln(v, line)
            }
    }
    gui.updateKeyBindingsView(g, pullViewFeature.Name)
    if _, err := g.SetCurrentView(pullViewFeature.Name); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) closePullView(g *gocui.Gui, v *gocui.View) error {

        if err := g.DeleteView(v.Name()); err != nil {
            return nil
        }
        if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
            return err
        }
        gui.updateKeyBindingsView(g, mainViewFeature.Name)

    return nil
}

func (gui *Gui) executePull(g *gocui.Gui, v *gocui.View) error {

    updateKeyBindingsViewForExecution(g)

    mrs, _ := gui.getMarkedEntities()
    for _, mr := range mrs {
        // here we will be waiting
        mr.PullTest()
        gui.updateCommits(g, mr)
        mr.Unmark()
    }
    gui.closePullView(g,v)
    gui.refreshMain(g)
    gui.updateJobs(g)

    return nil
}

func updateKeyBindingsViewForExecution(g *gocui.Gui) error {

    v, err := g.View(keybindingsViewFeature.Name)
    if err != nil {
        return err
    }
    v.Clear()
    v.BgColor = gocui.ColorRed
    v.FgColor = gocui.ColorWhite
    v.Frame = false
    fmt.Fprintln(v, " PULLING REPOSITORIES")
    return nil
}

func (gui *Gui) updatePullViewWithExec(g *gocui.Gui) {

    v, err := g.View(pullViewFeature.Name)
    if err != nil {
        return
    }

    g.Update(func(g *gocui.Gui) error {
        v.Clear()
        fmt.Fprintln(v, "Pulling...")
        return nil
    })
}
