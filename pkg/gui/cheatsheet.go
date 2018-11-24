package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) openCheatSheetView(g *gocui.Gui, v *gocui.View) error {
    maxX, maxY := g.Size()

    v, err := g.SetView("cheatsheet", maxX/2-25, maxY/2-10, maxX/2+25, maxY/2+10)
    if err != nil {
            if err != gocui.ErrUnknownView {
                    return err
            }
            v.Title = " " + "Application Controls" + " "
            fmt.Fprintln(v, " ")
            fmt.Fprintln(v, " a: select all")
            fmt.Fprintln(v, " b: iterate over branch")
            fmt.Fprintln(v, " c: close window")
            fmt.Fprintln(v, " d: deselect all")
            fmt.Fprintln(v, " r: iterate over remote")
            fmt.Fprintln(v, " s: iterate over commit")
            fmt.Fprintln(v, " x: show commit detail")
            fmt.Fprintln(v, " enter: execute")
            fmt.Fprintln(v, " q: quit")
            fmt.Fprintln(v, " ctrl+c: force quit")
    }
    gui.updateKeyBindingsViewForCheatSheetView(g)
    if _, err := g.SetCurrentView("cheatsheet"); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) closeCheatSheetView(g *gocui.Gui, v *gocui.View) error {

        if err := g.DeleteView(v.Name()); err != nil {
            return nil
        }
        if _, err := g.SetCurrentView("main"); err != nil {
            return err
        }
        gui.updateKeyBindingsViewForMainView(g)

    return nil
}

func (gui *Gui) updateKeyBindingsViewForCheatSheetView(g *gocui.Gui) error {

    v, err := g.View("keybindings")
    if err != nil {
        return err
    }
    v.Clear()
    v.BgColor = gocui.ColorWhite
    v.FgColor = gocui.ColorBlack
    v.Frame = false
    fmt.Fprintln(v, "c: cancel")
    return nil
}