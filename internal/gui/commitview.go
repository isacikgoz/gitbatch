package gui

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/isacikgoz/gitbatch/internal/command"
	"github.com/jroimartin/gocui"
)

var (
	commitFrameViewFeature          = viewFeature{Name: "commitframe", Title: " Frame "}
	commitUserNameLabelFeature      = viewFeature{Name: "commitusernamelabel", Title: " Name: "}
	commitUserEmailLabelViewFeature = viewFeature{Name: "commituseremaillabel", Title: " E-Mail: "}

	// these views used as a input for the credentials
	commitMessageViewFeature   = viewFeature{Name: "commitmessage", Title: " Commit Mesage "}
	commitUserUserViewFeature  = viewFeature{Name: "commitusername", Title: " Name "}
	commitUserEmailViewFeature = viewFeature{Name: "commituseremail", Title: " E-Mail "}

	commitViews      = []viewFeature{commitMessageViewFeature, commitUserUserViewFeature, commitUserEmailViewFeature}
	commitLabelViews = []viewFeature{commitFrameViewFeature, commitUserNameLabelFeature, commitUserEmailLabelViewFeature}
)

// open the commit message views
func (gui *Gui) openCommitMessageView(g *gocui.Gui, _ *gocui.View) error {
	maxX, maxY := g.Size()
	vFrame, err := g.SetView(commitFrameViewFeature.Name, maxX/2-30, maxY/2-4, maxX/2+30, maxY/2+3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vFrame.Frame = true
		fmt.Fprintln(vFrame, " Enter your commit message:")
	}
	v, err := g.SetView(commitMessageViewFeature.Name, maxX/2-29, maxY/2-3, maxX/2+29, maxY/2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true
		v.Editable = true
		v.Editor = gocui.DefaultEditor
		g.Cursor = true
	}
	if err := gui.openCommitUserNameView(g); err != nil {
		return err
	}
	if err := gui.openCommitUserEmailView(g); err != nil {
		return err
	}
	return gui.focusToView(commitMessageViewFeature.Name)
}

// open an error view to inform user with a message and a useful note
func (gui *Gui) openCommitUserNameView(g *gocui.Gui) error {
	r := gui.getSelectedRepository()
	maxX, maxY := g.Size()
	// first, create the label for user
	vlabel, err := g.SetView(commitUserNameLabelFeature.Name, maxX/2-30, maxY/2, maxX/2-19, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(vlabel, commitUserNameLabelFeature.Title)
		vlabel.Frame = false
	}
	// second, crete the user input
	v, err := g.SetView(commitUserUserViewFeature.Name, maxX/2-18, maxY/2, maxX/2+29, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		name, err := command.Config(r, &command.ConfigOptions{
			Section:     "user",
			Option:      "name",
			CommandMode: command.ModeLegacy,
		})
		if err != nil {
			// possibly could not get the user name
			name = ""
		}
		fmt.Fprintln(v, name)
		v.Editable = true
		v.Frame = false
	}
	return nil
}

// open an error view to inform user with a message and a useful note
func (gui *Gui) openCommitUserEmailView(g *gocui.Gui) error {
	r := gui.getSelectedRepository()
	maxX, maxY := g.Size()
	// first, create the label for password
	vlabel, err := g.SetView(commitUserEmailLabelViewFeature.Name, maxX/2-30, maxY/2+1, maxX/2-19, maxY/2+3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(vlabel, commitUserEmailLabelViewFeature.Title)
		vlabel.Frame = false
	}
	// second, crete the masked password input
	v, err := g.SetView(commitUserEmailViewFeature.Name, maxX/2-18, maxY/2+1, maxX/2+29, maxY/2+3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		email, err := command.Config(r, &command.ConfigOptions{
			Section:     "user",
			Option:      "email",
			CommandMode: command.ModeLegacy,
		})
		if err != nil {
			// possibly could not get the user email
			email = ""
		}
		fmt.Fprintln(v, email)
		v.Editable = true
		v.Frame = false
	}
	return nil
}

// close the opened commite mesage view
func (gui *Gui) submitCommitMessageView(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()

	// in order to read buffer of the views, first we need to find'em
	vMsg, err := g.View(commitMessageViewFeature.Name)
	if err != nil {
		return err // should return??
	}

	vName, err := g.View(commitUserUserViewFeature.Name)
	if err != nil {
		return err // should return??
	}

	vEmail, err := g.View(commitUserEmailViewFeature.Name)
	if err != nil {
		return err // should return??
	}

	// the return string of the views contain trailing new lines
	re := regexp.MustCompile(`\r?\n`)
	// TODO: maybe intentionally added new lines?
	msg := re.ReplaceAllString(vMsg.ViewBuffer(), "")
	name := re.ReplaceAllString(vName.ViewBuffer(), "")
	email := re.ReplaceAllString(vEmail.ViewBuffer(), "")
	if len(email) <= 0 {
		return errors.New("User email needs to be provided")
	}

	err = command.Commit(r, &command.CommitOptions{
		CommitMsg:   msg,
		User:        name,
		Email:       email,
		CommandMode: command.ModeNative,
	})
	if err != nil {
		return err
	}

	return gui.closeCommitMessageView(g, v)
}

// focus to next view
func (gui *Gui) nextCommitView(g *gocui.Gui, v *gocui.View) error {
	return gui.nextViewOfGroup(g, v, commitViews)
}

// close the opened commite mesage view
func (gui *Gui) closeCommitMessageView(g *gocui.Gui, v *gocui.View) error {
	// r := gui.getSelectedRepository()
	g.Cursor = false
	for _, view := range commitViews {
		if err := g.DeleteView(view.Name); err != nil {
			return err
		}
	}
	for _, view := range commitLabelViews {
		if err := g.DeleteView(view.Name); err != nil {
			return err
		}
	}
	if err := gui.focusToRepository(g, v); err != nil {
		return err
	}
	return nil
}
