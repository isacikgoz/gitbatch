package gui

import (
    "github.com/jroimartin/gocui"
)

func (gui *Gui) keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("", 'q', gocui.ModNone, gui.quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, gui.execute); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, gui.cursorDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, gui.cursorUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, gui.markRepository); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'a', gocui.ModNone, gui.markAllRepositories); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'r', gocui.ModNone, gui.unMarkAllRepositories); err != nil {
        return err
    }
    if err := g.SetKeybinding("", 'w', gocui.ModNone,
                func(g *gocui.Gui, v *gocui.View) error {
                    return gui.delView(g)
                }); err != nil {
                return err
    }
    return nil
}