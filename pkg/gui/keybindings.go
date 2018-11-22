package gui

import (
    "github.com/jroimartin/gocui"
)

func (gui *Gui) keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'q', gocui.ModNone, gui.quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, gui.openPullView); err != nil {
        return err
    }
    if err := g.SetKeybinding("pull", 'c', gocui.ModNone, gui.closePullView); err != nil {
        return err
    }
    if err := g.SetKeybinding("pull", gocui.KeyEnter, gocui.ModNone, gui.executePull); err != nil {
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
    if err := g.SetKeybinding("main", 'c', gocui.ModNone, gui.openCheatSheetView); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'a', gocui.ModNone, gui.markAllRepositories); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'd', gocui.ModNone, gui.unMarkAllRepositories); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'b', gocui.ModNone, gui.nextBranch); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'r', gocui.ModNone, gui.nextRemote); err != nil {
        return err
    }
    if err := g.SetKeybinding("cheatsheet", 'c', gocui.ModNone, gui.closeCheatSheetView); err != nil {
        return err
    }
    return nil
}