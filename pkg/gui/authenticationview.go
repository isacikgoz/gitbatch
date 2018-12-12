package gui

import (
	"fmt"
	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

var (
	authenticationReturnView  string
	authenticationViewFeature = viewFeature{Name: "authentication", Title: " Authentication "}
	authUserFeature           = viewFeature{Name: "authuser", Title: " User "}
	authPasswordViewFeature   = viewFeature{Name: "authpasswd", Title: " Password "}
	authUserLabelFeature      = viewFeature{Name: "authuserlabel", Title: " User: "}
	authPswdLabelViewFeature  = viewFeature{Name: "authpasswdlabel", Title: " Password: "}

	authViews  = []viewFeature{authUserFeature, authPasswordViewFeature}
	authLabels = []viewFeature{authenticationViewFeature, authUserLabelFeature, authPswdLabelViewFeature}

	jobRequiresAuth *git.Job
)

// open an error view to inform user with a message and a useful note
func (gui *Gui) openAuthenticationView(g *gocui.Gui, jobQueue *git.JobQueue, job *git.Job, returnViewName string) error {
	maxX, maxY := g.Size()
	// lets remove this job from the queue so that it won't block anything
	// also it is already unsuccessfully finished
	jobRequiresAuth = job
	if job.Entity.State != git.Fail {
		if err := jobQueue.RemoveFromQueue(job.Entity); err != nil {
			log.Fatal(err.Error())
			return err
		}
	}

	authenticationReturnView = returnViewName
	v, err := g.SetView(authenticationViewFeature.Name, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, " Enter credentials for: "+red.Sprint(job.Entity.Remote.URL[0]))
	}
	g.Cursor = true
	if err := gui.openUserView(g); err != nil {
		return nil
	}
	if err := gui.openPasswordView(g); err != nil {
		return nil
	}
	return nil
}

// close the opened error view
func (gui *Gui) closeAuthenticationView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	for _, vf := range authLabels {
		if err := g.DeleteView(vf.Name); err != nil {
			return nil
		}
	}
	for _, vf := range authViews {
		if err := g.DeleteView(vf.Name); err != nil {
			return nil
		}
	}
	if _, err := g.SetCurrentView(authenticationReturnView); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, authenticationReturnView)
	return nil
}

// close the opened error view
func (gui *Gui) submitAuthenticationView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	v_user, err := g.View(authUserFeature.Name)
	v_pswd, err := g.View(authPasswordViewFeature.Name)
	creduser := v_user.ViewBuffer()
	credpswd := v_pswd.ViewBuffer()
	// Maybe pause implementation can be added
	switch mode := jobRequiresAuth.JobType; mode {
	case git.FetchJob:
		jobRequiresAuth.Options = git.FetchOptions{
			RemoteName: jobRequiresAuth.Entity.Remote.Name,
			Credentials: git.Credentials{
				User:     creduser,
				Password: credpswd,
			},
		}
	}
	err = gui.State.Queue.AddJob(jobRequiresAuth)
	if err != nil {
		return err
	}
	jobRequiresAuth.Entity.State = git.Queued
	gui.refreshMain(g)
	gui.refreshViews(g, jobRequiresAuth.Entity)
	gui.closeAuthenticationView(g, v)
	return nil
}

// open an error view to inform user with a message and a useful note
func (gui *Gui) openUserView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	vlabel, err := g.SetView(authUserLabelFeature.Name, maxX/2-30, maxY/2-1, maxX/2-19, maxY/2+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(vlabel, authUserLabelFeature.Title)
		vlabel.Frame = false
	}
	v, err := g.SetView(authUserFeature.Name, maxX/2-18, maxY/2-1, maxX/2+29, maxY/2+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = authUserFeature.Title
		v.Editable = true
		v.Frame = false
	}
	gui.updateKeyBindingsView(g, authUserFeature.Name)
	if _, err := g.SetCurrentView(authUserFeature.Name); err != nil {
		return err
	}
	return nil
}

// open an error view to inform user with a message and a useful note
func (gui *Gui) openPasswordView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	vlabel, err := g.SetView(authPswdLabelViewFeature.Name, maxX/2-30, maxY/2, maxX/2-19, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(vlabel, authPswdLabelViewFeature.Title)
		vlabel.Frame = false
	}
	v, err := g.SetView(authPasswordViewFeature.Name, maxX/2-18, maxY/2, maxX/2+29, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = authPasswordViewFeature.Title
		v.Editable = true
		v.Mask ^= '*'
		v.Frame = false
	}
	return nil
}

// focus to next view
func (gui *Gui) nextAuthView(g *gocui.Gui, v *gocui.View) error {
	err := gui.nextViewOfGroup(g, v, authViews)
	return err
}

// focus to previous view
func (gui *Gui) previousAuthView(g *gocui.Gui, v *gocui.View) error {
	err := gui.previousViewOfGroup(g, v, authViews)
	return err
}
