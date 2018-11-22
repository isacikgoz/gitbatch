package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
    "strconv"
    "github.com/fatih/color"
)

func (gui *Gui) updateJobs(g *gocui.Gui) error {
    var err error

    out, err := g.View("jobs")
    if err != nil {
        return err
    }
    out.Clear()
    pullJobs := 0
    for _, r := range gui.State.Repositories {
        if r.Marked {
            pullJobs++
        }
    }
    fcolor := color.New(color.FgWhite)
    if pullJobs > 0 {
        fcolor = color.New(color.FgGreen)
    }
    
    jobs := strconv.Itoa(pullJobs) + " repositories to pull"
    fmt.Fprintln(out, fcolor.Sprint(jobs))
    return nil
}