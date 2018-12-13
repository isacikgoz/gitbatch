package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// KeyBinding structs is helpful for not re-writinh the same function over and
// over again. it hold useful values to generate a controls view
type KeyBinding struct {
	View        string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         interface{}
	Modifier    gocui.Modifier
	Display     string
	Description string
	Vital       bool
}

// generate the gui's controls a.k.a. keybindings
func (gui *Gui) generateKeybindings() error {
	// Mainviews common keybindings
	for _, view := range mainViews {
		mainKeybindings := []*KeyBinding{
			{
				View:        view.Name,
				Key:         'q',
				Modifier:    gocui.ModNone,
				Handler:     gui.quit,
				Display:     "q",
				Description: "Quit",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         'f',
				Modifier:    gocui.ModNone,
				Handler:     gui.switchToFetchMode,
				Display:     "f",
				Description: "Fetch mode",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         'p',
				Modifier:    gocui.ModNone,
				Handler:     gui.switchToPullMode,
				Display:     "p",
				Description: "Pull mode",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         'm',
				Modifier:    gocui.ModNone,
				Handler:     gui.switchToMergeMode,
				Display:     "m",
				Description: "Merge mode",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         gocui.KeyTab,
				Modifier:    gocui.ModNone,
				Handler:     gui.nextMainView,
				Display:     "tab",
				Description: "Next Panel",
				Vital:       false,
			},
		}
		gui.KeyBindings = append(gui.KeyBindings, mainKeybindings...)
	}
	// Statusviews common keybindings
	for _, view := range statusViews {
		statusKeybindings := []*KeyBinding{
			{
				View:        view.Name,
				Key:         gocui.KeyEsc,
				Modifier:    gocui.ModNone,
				Handler:     gui.closeStatusView,
				Display:     "esc",
				Description: "Close/Cancel",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         gocui.KeyTab,
				Modifier:    gocui.ModNone,
				Handler:     gui.nextStatusView,
				Display:     "tab",
				Description: "Next Panel",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         gocui.KeyArrowUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorUp,
				Display:     "↑",
				Description: "Up",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         gocui.KeyArrowDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorDown,
				Display:     "↓",
				Description: "Down",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         'k',
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorUp,
				Display:     "k",
				Description: "Up",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         'j',
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorDown,
				Display:     "j",
				Description: "Down",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         't',
				Modifier:    gocui.ModNone,
				Handler:     gui.stashChanges,
				Display:     "t",
				Description: "Save to Stash",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         'm',
				Modifier:    gocui.ModNone,
				Handler:     gui.openCommitMessageView,
				Display:     "m",
				Description: "Commit Changes",
				Vital:       true,
			},
		}
		gui.KeyBindings = append(gui.KeyBindings, statusKeybindings...)
	}
	for _, view := range authViews {
		authKeybindings := []*KeyBinding{
			{
				View:        view.Name,
				Key:         gocui.KeyEsc,
				Modifier:    gocui.ModNone,
				Handler:     gui.closeAuthenticationView,
				Display:     "esc",
				Description: "close/cancel",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         gocui.KeyTab,
				Modifier:    gocui.ModNone,
				Handler:     gui.nextAuthView,
				Display:     "tab",
				Description: "Next Panel",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         gocui.KeyEnter,
				Modifier:    gocui.ModNone,
				Handler:     gui.submitAuthenticationView,
				Display:     "enter",
				Description: "Submit",
				Vital:       true,
			},
		}
		gui.KeyBindings = append(gui.KeyBindings, authKeybindings...)
	}
	for _, view := range commitViews {
		commitKeybindings := []*KeyBinding{
			{
				View:        view.Name,
				Key:         gocui.KeyEsc,
				Modifier:    gocui.ModNone,
				Handler:     gui.closeCommitMessageView,
				Display:     "esc",
				Description: "close/cancel",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         gocui.KeyTab,
				Modifier:    gocui.ModNone,
				Handler:     gui.nextCommitView,
				Display:     "tab",
				Description: "Next Panel",
				Vital:       true,
			}, {
				View:        view.Name,
				Key:         gocui.KeyEnter,
				Modifier:    gocui.ModNone,
				Handler:     gui.submitCommitMessageView,
				Display:     "enter",
				Description: "Submit",
				Vital:       true,
			},
		}
		gui.KeyBindings = append(gui.KeyBindings, commitKeybindings...)
	}
	individualKeybindings := []*KeyBinding{
		// stash view
		{
			View:        stashViewFeature.Name,
			Key:         'p',
			Modifier:    gocui.ModNone,
			Handler:     gui.popStash,
			Display:     "p",
			Description: "Pop Item",
			Vital:       true,
		}, {
			View:        stashViewFeature.Name,
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.showStash,
			Display:     "d",
			Description: "Show diff",
			Vital:       true,
		},
		// staged view
		{
			View:        stageViewFeature.Name,
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.resetChanges,
			Display:     "r",
			Description: "Reset Item",
			Vital:       true,
		}, {
			View:        stageViewFeature.Name,
			Key:         gocui.KeyCtrlR,
			Modifier:    gocui.ModNone,
			Handler:     gui.resetAllChanges,
			Display:     "ctrl+r",
			Description: "Reset All Items",
			Vital:       true,
		}, {
			View:        stageViewFeature.Name,
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.openFileDiffView,
			Display:     "d",
			Description: "Show diff",
			Vital:       true,
		},
		// unstaged view
		{
			View:        unstageViewFeature.Name,
			Key:         'a',
			Modifier:    gocui.ModNone,
			Handler:     gui.addChanges,
			Display:     "a",
			Description: "Add Item",
			Vital:       true,
		}, {
			View:        unstageViewFeature.Name,
			Key:         gocui.KeyCtrlA,
			Modifier:    gocui.ModNone,
			Handler:     gui.addAllChanges,
			Display:     "ctrl+a",
			Description: "Add All Items",
			Vital:       true,
		}, {
			View:        unstageViewFeature.Name,
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.openFileDiffView,
			Display:     "d",
			Description: "Show diff",
			Vital:       true,
		},
		// Main view controls
		{
			View:        mainViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.cursorUp,
			Display:     "↑",
			Description: "Up",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.cursorDown,
			Display:     "↓",
			Description: "Down",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.cursorUp,
			Display:     "k",
			Description: "Up",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.cursorDown,
			Display:     "j",
			Description: "Down",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.markRepository,
			Display:     "space",
			Description: "Select",
			Vital:       true,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.startQueue,
			Display:     "enter",
			Description: "Start",
			Vital:       true,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyCtrlSpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.markAllRepositories,
			Display:     "ctrl + space",
			Description: "Select All",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyBackspace2,
			Modifier:    gocui.ModNone,
			Handler:     gui.unmarkAllRepositories,
			Display:     "backspace",
			Description: "Deselect All",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.openCheatSheetView,
			Display:     "c",
			Description: "Controls",
			Vital:       true,
		}, {
			View:        mainViewFeature.Name,
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.sortByName,
			Display:     "n",
			Description: "Sort repositories by Name",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.sortByMod,
			Display:     "r",
			Description: "Sort repositories by Modification date",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.openStatusView,
			Display:     "s",
			Description: "Open Status",
			Vital:       true,
		}, {
			View:        "",
			Key:         gocui.KeyCtrlC,
			Modifier:    gocui.ModNone,
			Handler:     gui.quit,
			Display:     "ctrl + c",
			Description: "Force application to quit",
			Vital:       false,
		},
		// Branch View Controls
		{
			View:        branchViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.nextBranch,
			Display:     "↓",
			Description: "Iterate over branches",
			Vital:       false,
		}, {
			View:        branchViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.previousBranch,
			Display:     "↑",
			Description: "Iterate over branches",
			Vital:       false,
		}, {
			View:        branchViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.nextBranch,
			Display:     "j",
			Description: "Down",
			Vital:       false,
		}, {
			View:        branchViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.previousBranch,
			Display:     "k",
			Description: "Up",
			Vital:       false,
		},
		// Remote View Controls
		{
			View:        remoteViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.nextRemote,
			Display:     "↓",
			Description: "Iterate over remotes",
			Vital:       false,
		}, {
			View:        remoteViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.previousRemote,
			Display:     "↑",
			Description: "Iterate over remotes",
			Vital:       false,
		}, {
			View:        remoteViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.nextRemote,
			Display:     "j",
			Description: "Down",
			Vital:       false,
		}, {
			View:        remoteViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.previousRemote,
			Display:     "k",
			Description: "Up",
			Vital:       false,
		},
		// Remote Branch View Controls
		{
			View:        remoteBranchViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.nextRemoteBranch,
			Display:     "↓",
			Description: "Iterate over remote branches",
			Vital:       false,
		}, {
			View:        remoteBranchViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.previousRemoteBranch,
			Display:     "↑",
			Description: "Iterate over remote branches",
			Vital:       false,
		}, {
			View:        remoteBranchViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.nextRemoteBranch,
			Display:     "j",
			Description: "Down",
			Vital:       false,
		}, {
			View:        remoteBranchViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.previousRemoteBranch,
			Display:     "k",
			Description: "Up",
			Vital:       false,
		}, {
			View:        remoteBranchViewFeature.Name,
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.syncRemoteBranch,
			Display:     "s",
			Description: "Synch with Remote",
			Vital:       true,
		},
		// Commit View Controls
		{
			View:        commitViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.nextCommit,
			Display:     "↓",
			Description: "Iterate over commits",
			Vital:       false,
		}, {
			View:        commitViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.prevCommit,
			Display:     "↑",
			Description: "Iterate over commits",
			Vital:       false,
		}, {
			View:        commitViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.nextCommit,
			Display:     "j",
			Description: "Down",
			Vital:       false,
		}, {
			View:        commitViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.prevCommit,
			Display:     "k",
			Description: "Up",
			Vital:       false,
		}, {
			View:        commitViewFeature.Name,
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.openCommitDiffView,
			Display:     "d",
			Description: "Show commit diff",
			Vital:       true,
		},
		// Diff View Controls
		{
			View:        diffViewFeature.Name,
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.closeCommitDiffView,
			Display:     "esc",
			Description: "close/cancel",
			Vital:       true,
		}, {
			View:        diffViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorUp,
			Display:     "↑",
			Description: "Page up",
			Vital:       true,
		}, {
			View:        diffViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorDown,
			Display:     "↓",
			Description: "Page down",
			Vital:       true,
		}, {
			View:        diffViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorUp,
			Display:     "k",
			Description: "Page up",
			Vital:       false,
		}, {
			View:        diffViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorDown,
			Display:     "j",
			Description: "Page down",
			Vital:       false,
		},
		// Application Controls
		{
			View:        cheatSheetViewFeature.Name,
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.closeCheatSheetView,
			Display:     "esc",
			Description: "close/cancel",
			Vital:       true,
		}, {
			View:        cheatSheetViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorUp,
			Display:     "↑",
			Description: "Up",
			Vital:       true,
		}, {
			View:        cheatSheetViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorDown,
			Display:     "↓",
			Description: "Down",
			Vital:       true,
		}, {
			View:        cheatSheetViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorUp,
			Display:     "k",
			Description: "Up",
			Vital:       false,
		}, {
			View:        cheatSheetViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorDown,
			Display:     "j",
			Description: "Down",
			Vital:       false,
		},
		// Error View
		{
			View:        errorViewFeature.Name,
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.closeErrorView,
			Display:     "esc",
			Description: "close/cancel",
			Vital:       true,
		}, {
			View:        errorViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorUp,
			Display:     "↑",
			Description: "Up",
			Vital:       true,
		}, {
			View:        errorViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorDown,
			Display:     "↓",
			Description: "Down",
			Vital:       true,
		}, {
			View:        errorViewFeature.Name,
			Key:         'k',
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorUp,
			Display:     "k",
			Description: "Up",
			Vital:       false,
		}, {
			View:        errorViewFeature.Name,
			Key:         'j',
			Modifier:    gocui.ModNone,
			Handler:     gui.fastCursorDown,
			Display:     "j",
			Description: "Down",
			Vital:       false,
		},
	}
	gui.KeyBindings = append(gui.KeyBindings, individualKeybindings...)
	return nil
}

// set the guis by iterating over a slice of the gui's keybindings struct
func (gui *Gui) keybindings(g *gocui.Gui) error {
	for _, k := range gui.KeyBindings {
		if err := g.SetKeybinding(k.View, k.Key, k.Modifier, k.Handler); err != nil {
			return err
		}
	}
	return nil
}

// the bottom line of the gui is mode indicator and keybindings view. Only the
// important controls (marked as vital) are shown
func (gui *Gui) updateKeyBindingsView(g *gocui.Gui, viewName string) error {
	v, err := g.View(keybindingsViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	v.BgColor = gocui.ColorWhite
	v.FgColor = gocui.ColorBlack
	v.Frame = false
	fmt.Fprint(v, ws)
	modeLabel := ""
	switch mode := gui.State.Mode.ModeID; mode {
	case FetchMode:
		v.BgColor = gocui.ColorBlue
		modeLabel = fetchSymbol + ws + "FETCH"
	case PullMode:
		v.BgColor = gocui.ColorMagenta
		modeLabel = pullSymbol + ws + "PULL"
	case MergeMode:
		v.BgColor = gocui.ColorCyan
		modeLabel = mergeSymbol + ws + "MERGE"
	default:
		modeLabel = "No mode selected"
	}

	fmt.Fprint(v, ws+modeLabel+ws+modeSeperator)

	for _, k := range gui.KeyBindings {
		if k.View == viewName && k.Vital {
			binding := keyBindingSeperator + ws + k.Display + ":" + ws + k.Description + ws
			fmt.Fprint(v, binding)
		}
	}
	return nil
}
