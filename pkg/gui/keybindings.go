package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
)

type KeyBinding struct {
    View        string
    Handler     func(*gocui.Gui, *gocui.View) error
    Key         interface{}
    Modifier    gocui.Modifier
    Display     string
    Description string
    Vital       bool
}

func (gui *Gui) generateKeybindings() error {
    gui.KeyBindings = []*KeyBinding{
        {
            View:        "",
            Key:         gocui.KeyCtrlC,
            Modifier:    gocui.ModNone,
            Handler:     gui.quit,
            Display:     "ctrl + c",
            Description: "Force application to quit",
            Vital:       false,
        },{
            View:        mainViewFeature.Name,
            Key:         'q',
            Modifier:    gocui.ModNone,
            Handler:     gui.quit,
            Display:     "q",
            Description: "Quit from application",
            Vital:       true,
        },{
            View:        mainViewFeature.Name,
            Key:         gocui.KeyArrowUp,
            Modifier:    gocui.ModNone,
            Handler:     gui.cursorUp,
            Display:     "↑",
            Description: "Cursor up",
            Vital:       true,
        },{
            View:        mainViewFeature.Name,
            Key:         gocui.KeyArrowDown,
            Modifier:    gocui.ModNone,
            Handler:     gui.cursorDown,
            Display:     "↓",
            Description: "Cursor down",
            Vital:       true,
        },{
            View:        mainViewFeature.Name,
            Key:         'b',
            Modifier:    gocui.ModNone,
            Handler:     gui.nextBranch,
            Display:     "b",
            Description: "Iterate over branches",
            Vital:       false,
        },{
            View:        mainViewFeature.Name,
            Key:         'r',
            Modifier:    gocui.ModNone,
            Handler:     gui.nextRemote,
            Display:     "r",
            Description: "Iterate over remotes",
            Vital:       false,
        },{
            View:        mainViewFeature.Name,
            Key:         's',
            Modifier:    gocui.ModNone,
            Handler:     gui.nextCommit,
            Display:     "s",
            Description: "Iterate over commits",
            Vital:       false,
        },{
            View:        mainViewFeature.Name,
            Key:         'x',
            Modifier:    gocui.ModNone,
            Handler:     gui.showCommitDetail,
            Display:     "x",
            Description: "Show commit diff",
            Vital:       false,
        },{
            View:        mainViewFeature.Name,
            Key:         'c',
            Modifier:    gocui.ModNone,
            Handler:     gui.openCheatSheetView,
            Display:     "c",
            Description: "Open cheatsheet window",
            Vital:       false,
        },{
            View:        mainViewFeature.Name,
            Key:         gocui.KeyEnter,
            Modifier:    gocui.ModNone,
            Handler:     gui.openPullView,
            Display:     "enter",
            Description: "Execute jobs",
            Vital:       true,
        },{
            View:        mainViewFeature.Name,
            Key:         gocui.KeySpace,
            Modifier:    gocui.ModNone,
            Handler:     gui.markRepository,
            Display:     "space",
            Description: "Select repository",
            Vital:       true,
        },{
            View:        mainViewFeature.Name,
            Key:         'a',
            Modifier:    gocui.ModNone,
            Handler:     gui.markAllRepositories,
            Display:     "a",
            Description: "Select all repositories",
            Vital:       false,
        },{
            View:        mainViewFeature.Name,
            Key:         'd',
            Modifier:    gocui.ModNone,
            Handler:     gui.unMarkAllRepositories,
            Display:     "d",
            Description: "Deselect all repositories",
            Vital:       false,
        },{
            View:        commitdetailViewFeature.Name,
            Key:         'c',
            Modifier:    gocui.ModNone,
            Handler:     gui.closeCommitDetailView,
            Display:     "c",
            Description: "close/cancel",
            Vital:       true,
        },{
            View:        commitdetailViewFeature.Name,
            Key:         gocui.KeyArrowUp,
            Modifier:    gocui.ModNone,
            Handler:     gui.commitCursorUp,
            Display:     "↑",
            Description: "Page up",
            Vital:       true,
        },{
            View:        commitdetailViewFeature.Name,
            Key:         gocui.KeyArrowDown,
            Modifier:    gocui.ModNone,
            Handler:     gui.commitCursorDown,
            Display:     "↓",
            Description: "Page down",
            Vital:       true,
        },{
            View:        cheatSheetViewFeature.Name,
            Key:         'c',
            Modifier:    gocui.ModNone,
            Handler:     gui.closeCheatSheetView,
            Display:     "c",
            Description: "close/cancel",
            Vital:       true,
        },{
            View:        pullViewFeature.Name,
            Key:         'c',
            Modifier:    gocui.ModNone,
            Handler:     gui.closePullView,
            Display:     "c",
            Description: "close/cancel",
            Vital:       true,
        },{
            View:        pullViewFeature.Name,
            Key:         gocui.KeyEnter,
            Modifier:    gocui.ModNone,
            Handler:     gui.executePull,
            Display:     "enter",
            Description: "Execute",
            Vital:       true,
        },{
            View:        errorViewFeature.Name,
            Key:         'c',
            Modifier:    gocui.ModNone,
            Handler:     gui.closeErrorView,
            Display:     "c",
            Description: "close/cancel",
            Vital:       true,
        },
    }
    return nil
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
    for _, k := range gui.KeyBindings {
        if err := g.SetKeybinding(k.View, k.Key, k.Modifier, k.Handler); err != nil {
            return err
        }
    }
    return nil
}

func (gui *Gui) updateKeyBindingsView(g *gocui.Gui, viewName string) error {
    v, err := g.View(keybindingsViewFeature.Name)
    if err != nil {
        return err
    }
    v.Clear()
    v.BgColor = gocui.ColorWhite
    v.FgColor = gocui.ColorBlack
    v.Frame = false
    for _, k := range gui.KeyBindings {
        if k.View == viewName && k.Vital {
            binding := " " + k.Display + ": " + k.Description + " |"
            fmt.Fprint(v, binding)
        }
    }
    return nil
}