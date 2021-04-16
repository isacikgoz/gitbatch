package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) updateDynamicKeybindings() error {
	vd, err := gui.g.View(dynamicViewFeature.Name)
	if err != nil {
		return err
	}
	t := DynamicViewMode(vd.Title)

	// the order is important here
	_ = gui.generateKeybindings()
	gui.g.DeleteKeybindings(vd.Name())

	keybindings := []*KeyBinding{
		{
			View:        dynamicViewFeature.Name,
			Key:         gocui.KeyTab,
			Modifier:    gocui.ModNone,
			Handler:     gui.focusBackToMain,
			Display:     "tab",
			Description: "Back to Overview",
			Vital:       false,
		}, {
			View:        dynamicViewFeature.Name,
			Key:         gocui.KeyArrowRight,
			Modifier:    gocui.ModNone,
			Handler:     gui.nextFocusView,
			Display:     "→",
			Description: "Next Panel",
			Vital:       false,
		}, {
			View:        dynamicViewFeature.Name,
			Key:         'l',
			Modifier:    gocui.ModNone,
			Handler:     gui.nextFocusView,
			Display:     "l",
			Description: "Next Panel",
			Vital:       false,
		}, {
			View:        dynamicViewFeature.Name,
			Key:         gocui.KeyArrowLeft,
			Modifier:    gocui.ModNone,
			Handler:     gui.previousFocusView,
			Display:     "←",
			Description: "Prev Panel",
			Vital:       false,
		}, {
			View:        dynamicViewFeature.Name,
			Key:         'h',
			Modifier:    gocui.ModNone,
			Handler:     gui.previousFocusView,
			Display:     "h",
			Description: "Prev Panel",
			Vital:       false,
		},
	}

	switch t {
	case CommitStatMode:
		caseBindings := []*KeyBinding{
			{
				View:        dynamicViewFeature.Name,
				Key:         'd',
				Modifier:    gocui.ModNone,
				Handler:     gui.commitDiff,
				Display:     "d",
				Description: "diff",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgup,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageUp,
				Display:     "pg up",
				Description: "Page up",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       true,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case CommitDiffMode:
		caseBindings := []*KeyBinding{
			{
				View:        dynamicViewFeature.Name,
				Key:         's',
				Modifier:    gocui.ModNone,
				Handler:     gui.commitStat,
				Display:     "s",
				Description: "stats",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgup,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageUp,
				Display:     "pg up",
				Description: "Page up",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       true,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case StashDiffMode:
		caseBindings := []*KeyBinding{
			{
				View:        dynamicViewFeature.Name,
				Key:         's',
				Modifier:    gocui.ModNone,
				Handler:     gui.commitDiff,
				Display:     "s",
				Description: "stats",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgup,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageUp,
				Display:     "pg up",
				Description: "Page up",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       true,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case StashStatMode:
		caseBindings := []*KeyBinding{
			{
				View:        dynamicViewFeature.Name,
				Key:         'd',
				Modifier:    gocui.ModNone,
				Handler:     gui.commitDiff,
				Display:     "d",
				Description: "diff",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgup,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageUp,
				Display:     "pg up",
				Description: "Page up",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       true,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case StatusMode:
		caseBindings := []*KeyBinding{
			{
				View:        dynamicViewFeature.Name,
				Key:         'd',
				Modifier:    gocui.ModNone,
				Handler:     gui.statusDiff,
				Display:     "d",
				Description: "diff",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         'c',
				Modifier:    gocui.ModNone,
				Handler:     gui.openCommitMessageView,
				Display:     "c",
				Description: "commit",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         't',
				Modifier:    gocui.ModNone,
				Handler:     gui.stashChanges,
				Display:     "t",
				Description: "stash",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyArrowDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorDown,
				Display:     "↓",
				Description: "Down",
				Vital:       false,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyArrowUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorUp,
				Display:     "↑",
				Description: "Up",
				Vital:       false,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         'j',
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorDown,
				Display:     "j",
				Description: "Down",
				Vital:       false,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         'k',
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorUp,
				Display:     "k",
				Description: "Up",
				Vital:       false,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeySpace,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusAddReset,
				Display:     "space",
				Description: "add/reset",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyCtrlA,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusAddAll,
				Display:     "c-a",
				Description: "add all",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyCtrlR,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusResetAll,
				Display:     "c-r",
				Description: "reset all",
				Vital:       true,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case FileDiffMode:
		caseBindings := []*KeyBinding{
			{
				View:        dynamicViewFeature.Name,
				Key:         's',
				Modifier:    gocui.ModNone,
				Handler:     gui.statusStat,
				Display:     "s",
				Description: "stats",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgup,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageUp,
				Display:     "pg up",
				Description: "Page up",
				Vital:       true,
			}, {
				View:        dynamicViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       true,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	default:

	}
	gui.KeyBindings = append(gui.KeyBindings, keybindings...)
	for _, k := range keybindings {
		if err := gui.g.SetKeybinding(k.View, k.Key, k.Modifier, k.Handler); err != nil {
			return err
		}
	}
	v := gui.g.CurrentView()
	if v.Name() == dynamicViewFeature.Name {
		_ = gui.updateKeyBindingsView(gui.g, dynamicViewFeature.Name)
	}
	return nil
}
