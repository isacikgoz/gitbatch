package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
    "strconv"
)

func (gui *Gui) openPullView(g *gocui.Gui, v *gocui.View) error {
    maxX, maxY := g.Size()

    v, err := g.SetView(pullViewFeature.Name, maxX/2-35, maxY/2-5, maxX/2+35, maxY/2+5)
    if err != nil {
            if err != gocui.ErrUnknownView {
                    return err
            }
            v.Title = pullViewFeature.Title
            v.Wrap = true
            mrs, _ := gui.getMarkedEntities()
            jobs := strconv.Itoa(len(mrs)) + " repositories to fetch & merge:"
            fmt.Fprintln(v, jobs)
            for _, r := range mrs {
                line := " - " + green.Sprint(r.Name) + ": " + r.Remote.Name + green.Sprint(" → ") + r.Branch.Name
                fmt.Fprintln(v, line)
            }
            ps := red.Sprint("Note:") + " After execution you will be notified"
            fmt.Fprintln(v, "\n" + ps)
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
    gui.refreshMain(g)
    gui.updateKeyBindingsView(g, mainViewFeature.Name)
    return nil
}

func (gui *Gui) executePull(g *gocui.Gui, v *gocui.View) error {
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
        mr.Pull()
        gui.updateCommits(g, mr)
        mr.Unmark()
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
    fmt.Fprintln(v, " Execution Completed; c: close/cancel")
    return nil
}