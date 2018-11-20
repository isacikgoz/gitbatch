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
    if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, gui.cursorDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, gui.cursorUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, gui.markRepository); err != nil {
        return err
    }
    return nil
}