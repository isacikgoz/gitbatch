package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) generateKeybindingsForDetailView(mode string) error {
	t := DetailViewMode(mode)
	keybindings := []*KeyBinding{
		{
			View:        detailViewFeature.Name,
			Key:         'q',
			Modifier:    gocui.ModNone,
			Handler:     gui.quit,
			Display:     "q",
			Description: "Quit",
			Vital:       true,
		}, {
			View:        detailViewFeature.Name,
			Key:         gocui.KeyTab,
			Modifier:    gocui.ModNone,
			Handler:     gui.nextMainView,
			Display:     "tab",
			Description: "Next Panel",
			Vital:       false,
		},
	}

	switch t {
	case CommitStatMode:
		caseBindings := []*KeyBinding{
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
				Vital:       false,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       false,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case CommitDiffMode:
		caseBindings := []*KeyBinding{
			{
				View:        detailViewFeature.Name,
				Key:         's',
				Modifier:    gocui.ModNone,
				Handler:     gui.commitStat,
				Display:     "s",
				Description: "stats",
				Vital:       true,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyPgup,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageUp,
				Display:     "pg up",
				Description: "Page up",
				Vital:       false,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       false,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case StashDiffMode:
		caseBindings := []*KeyBinding{
			{
				View:        detailViewFeature.Name,
				Key:         's',
				Modifier:    gocui.ModNone,
				Handler:     gui.commitDiff,
				Display:     "s",
				Description: "stats",
				Vital:       true,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyPgup,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageUp,
				Display:     "pg up",
				Description: "Page up",
				Vital:       false,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       false,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case StashStatMode:
		caseBindings := []*KeyBinding{
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
				Vital:       false,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       false,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case StatusMode:
		caseBindings := []*KeyBinding{
			{
				View:        detailViewFeature.Name,
				Key:         'd',
				Modifier:    gocui.ModNone,
				Handler:     gui.statusDiff,
				Display:     "d",
				Description: "diff",
				Vital:       true,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyArrowDown,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorDown,
				Display:     "↓",
				Description: "Down",
				Vital:       false,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyArrowUp,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusCursorUp,
				Display:     "↑",
				Description: "Up",
				Vital:       false,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeySpace,
				Modifier:    gocui.ModNone,
				Handler:     gui.statusAddReset,
				Display:     "space",
				Description: "add/reset",
				Vital:       true,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	case FileDiffMode:
		caseBindings := []*KeyBinding{
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
				Vital:       false,
			}, {
				View:        detailViewFeature.Name,
				Key:         gocui.KeyPgdn,
				Modifier:    gocui.ModNone,
				Handler:     gui.dpageDown,
				Display:     "pg down",
				Description: "Page Down",
				Vital:       false,
			},
		}
		keybindings = append(keybindings, caseBindings...)
	default:

	}

	// gui.KeyBindings = append(gui.KeyBindings, keybindings...)
	for _, k := range keybindings {
		if err := gui.g.SetKeybinding(k.View, k.Key, k.Modifier, k.Handler); err != nil {
			return err
		}
	}
	return nil
}
