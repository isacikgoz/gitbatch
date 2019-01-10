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
	for _, view := range sideViews {
		sideViewKeybindings := []*KeyBinding{
			{
				View:        view.Name,
				Key:         gocui.KeyArrowDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.sideCursorDown,
				Display:     "↓",
				Description: "Iterate over branches",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         gocui.KeyArrowUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.sideCursorUp,
				Display:     "↑",
				Description: "Iterate over branches",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         'j',
				Modifier:    gocui.ModNone,
				Handler:     gui.sideCursorDown,
				Display:     "j",
				Description: "Down",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         'k',
				Modifier:    gocui.ModNone,
				Handler:     gui.sideCursorUp,
				Display:     "k",
				Description: "Up",
				Vital:       false,
			}, {
				View:        view.Name,
				Key:         gocui.KeySpace,
				Modifier:    gocui.ModNone,
				Handler:     gui.selectSideItem,
				Display:     "space",
				Description: "select",
				Vital:       true,
			},
		}
		gui.KeyBindings = append(gui.KeyBindings, sideViewKeybindings...)
	}
	for _, view := range authViews {
		authKeybindings := []*KeyBinding{
			{
				View:        view.Name,
				Key:         gocui.KeyEsc,
				Modifier:    gocui.ModNone,
				Handler:     gui.closeAuthenticationView,
				Display:     "esc",
				Description: "Close/Cancel",
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
				Description: "Close/Cancel",
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
		// Main view controls
		{
			View:        mainViewFeature.Name,
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.focusToRepository,
			Display:     "f",
			Description: "Focus",
			Vital:       true,
		}, {
			View:        mainViewFeature.Name,
			Key:         'u',
			Modifier:    gocui.ModNone,
			Handler:     gui.submitCredentials,
			Display:     "u",
			Description: "Submit Credentials",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.cursorUp,
			Display:     "↑",
			Description: "Up",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyPgup,
			Modifier:    gocui.ModNone,
			Handler:     gui.pageUp,
			Display:     "page up",
			Description: "Page up",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyHome,
			Modifier:    gocui.ModNone,
			Handler:     gui.cursorTop,
			Display:     "home",
			Description: "Home",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyPgdn,
			Modifier:    gocui.ModNone,
			Handler:     gui.pageDown,
			Display:     "page down",
			Description: "Page Down",
			Vital:       false,
		}, {
			View:        mainViewFeature.Name,
			Key:         gocui.KeyEnd,
			Modifier:    gocui.ModNone,
			Handler:     gui.cursorEnd,
			Display:     "end",
			Description: "End",
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
			View:        "",
			Key:         gocui.KeyCtrlC,
			Modifier:    gocui.ModNone,
			Handler:     gui.quit,
			Display:     "ctrl + c",
			Description: "Force application to quit",
			Vital:       false,
		}, {
			View:        remoteBranchViewFeature.Name,
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.syncRemoteBranch,
			Display:     "s",
			Description: "Synch with Remote",
			Vital:       true,
		}, {
			View:        branchViewFeature.Name,
			Key:         'u',
			Modifier:    gocui.ModNone,
			Handler:     gui.setUpstreamToBranch,
			Display:     "u",
			Description: "Set Upstream",
			Vital:       true,
		},
		// CommitView
		{
			View:        commitViewFeature.Name,
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.focusBackToMain,
			Display:     "f",
			Description: "Back to Overview",
			Vital:       true,
		}, {
			View:        commitViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.commitCursorDown,
			Display:     "↓",
			Description: "Iterate over branches",
			Vital:       false,
		}, {
			View:        commitViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.commitCursorUp,
			Display:     "↑",
			Description: "Iterate over branches",
			Vital:       false,
		}, {
			View:        commitViewFeature.Name,
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.commitDiff,
			Display:     "d",
			Description: "Show commit diff",
			Vital:       true,
		},
		// Detailview
		{
			View:        detailViewFeature.Name,
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.commitDiff,
			Display:     "d",
			Description: "diff",
			Vital:       true,
		}, {
			View:        detailViewFeature.Name,
			Key:         gocui.KeyPgup,
			Modifier:    gocui.ModNone,
			Handler:     gui.dpageUp,
			Display:     "pg up",
			Description: "Page up",
			Vital:       true,
		}, {
			View:        detailViewFeature.Name,
			Key:         gocui.KeyPgdn,
			Modifier:    gocui.ModNone,
			Handler:     gui.dpageDown,
			Display:     "pg down",
			Description: "Page Down",
			Vital:       true,
		}, {
			View:        detailViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.focusStatusCursorDown,
			Display:     "↓",
			Description: "Down",
			Vital:       false,
		}, {
			View:        detailViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.focusStatusCursorUp,
			Display:     "↑",
			Description: "Up",
			Vital:       false,
		}, {
			View:        detailViewFeature.Name,
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.addreset,
			Display:     "space",
			Description: "add/reset",
			Vital:       false,
		},
		// stashview
		{
			View:        stashViewFeature.Name,
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.stashCursorDown,
			Display:     "↓",
			Description: "Iterate over branches",
			Vital:       false,
		}, {
			View:        stashViewFeature.Name,
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.stashCursorUp,
			Display:     "↑",
			Description: "Iterate over branches",
			Vital:       false,
		},
		// upstream confirmation
		{
			View:        confirmationViewFeature.Name,
			Key:         'q',
			Modifier:    gocui.ModNone,
			Handler:     gui.closeConfirmationView,
			Display:     "q",
			Description: "Close/Cancel",
			Vital:       true,
		}, {
			View:        confirmationViewFeature.Name,
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.confirmSetUpstreamToBranch,
			Display:     "enter",
			Description: "Set Upstream",
			Vital:       true,
		},
		// Diff View Controls
		{
			View:        diffViewFeature.Name,
			Key:         'q',
			Modifier:    gocui.ModNone,
			Handler:     gui.closeCommitDiffView,
			Display:     "q",
			Description: "Close/Cancel",
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
			Key:         'q',
			Modifier:    gocui.ModNone,
			Handler:     gui.closeCheatSheetView,
			Display:     "q",
			Description: "Close/Cancel",
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
			Key:         'q',
			Modifier:    gocui.ModNone,
			Handler:     gui.closeErrorView,
			Display:     "q",
			Description: "Close/Cancel",
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
