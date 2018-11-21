package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
//    "sync"
)

func (gui *Gui) openPullView(g *gocui.Gui, v *gocui.View) error {
    maxX, maxY := g.Size()

    v, err := g.SetView("pull", maxX/2-35, maxY/2-5, maxX/2+35, maxY/2+5)
    if err != nil {
            if err != gocui.ErrUnknownView {
                    return err
            }
            v.Title = " " + "Execution Parameters" + " "
            v.Wrap = false
            mrs, _ := gui.getMarkedEntities()
            for _, r := range mrs {
                line := r.Name + " : " + r.GetActiveRemote() + "/" + r.Branch + " → " + r.GetActiveBranch()
                fmt.Fprintln(v, line)
            }
    }
    gui.updateKeyBindingsViewForPullView(g)
    if _, err := g.SetCurrentView("pull"); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) closePullView(g *gocui.Gui, v *gocui.View) error {
 
        if err := g.DeleteView(v.Name()); err != nil {
            return nil
        }
        if _, err := g.SetCurrentView("main"); err != nil {
            return err
        }
        gui.updateKeyBindingsViewForMainView(g)
    
    return nil
}

func (gui *Gui) executePull(g *gocui.Gui, v *gocui.View) error {
    gui.updateKeyBindingsViewForExecution(g)
    mrs, _ := gui.getMarkedEntities()

    gui.updateKeyBindingsViewForExecution(g)

    //var wg sync.WaitGroup

    for _, mr := range mrs {
        //wg.Add(1)
        go func(g *gocui.Gui, v *gocui.View) {
            
            gui.updatePullViewWithExec(g)

            // here we will be waiting
            mr.PullTest()
            gui.updateCommits(g, mr)
            mr.Unmark()

            //not working somehow
            gui.closePullView(g, v)
            //defer wg.Done()
        }(g, v)
    }

    //wg.Wait()
    gui.refreshMain(g)
    gui.updateSchedule(g)
    return nil
}

func (gui *Gui) updateKeyBindingsViewForPullView(g *gocui.Gui) error {

    v, err := g.View("keybindings")
    if err != nil {
        return err
    }
    v.Clear()
    v.BgColor = gocui.ColorGreen
    v.FgColor = gocui.ColorBlack
    v.Frame = false
    fmt.Fprintln(v, "c: cancel | ↑ ↓: navigate | enter: execute")
    return nil
}


func (gui *Gui) updateKeyBindingsViewForExecution(g *gocui.Gui) error {

    v, err := g.View("keybindings")
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
    
    v, err := g.View("pull")
    if err != nil {
        return
    }
    
    g.Update(func(g *gocui.Gui) error {
        v.Clear()
        fmt.Fprintln(v, "Pulling...")
        return nil
    })
}
