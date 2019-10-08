package gui

import (
	"fmt"
	"regexp"

	"github.com/isacikgoz/gitbatch/internal/command"
	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/isacikgoz/gitbatch/internal/job"
	"github.com/jroimartin/gocui"
)

var (
	// this is required so we can know where we can return
	authenticationReturnView string

	// these views used as a label for git repository address and credential views
	authenticationViewFeature = viewFeature{Name: "authentication", Title: " Authentication "}
	authUserLabelFeature      = viewFeature{Name: "authuserlabel", Title: " User: "}
	authPswdLabelViewFeature  = viewFeature{Name: "authpasswdlabel", Title: " Password: "}

	// these views used as a input for the credentials
	authUserFeature         = viewFeature{Name: "authuser", Title: " User "}
	authPasswordViewFeature = viewFeature{Name: "authpasswd", Title: " Password "}

	// these are the view groups, so that we can assign common keybindings
	authViews  = []viewFeature{authUserFeature, authPasswordViewFeature}
	authLabels = []viewFeature{authenticationViewFeature, authUserLabelFeature, authPswdLabelViewFeature}

	// we can hold the job that is required to authenticate
	jobRequiresAuth *job.Job
)

// open an auth view to get user credentials
func (gui *Gui) openAuthenticationView(g *gocui.Gui, jobQueue *job.Queue, job *job.Job, returnViewName string) error {
	maxX, maxY := g.Size()
	// lets add this job since it is removed from the queue
	// also it is already unsuccessfully finished
	if err := jobQueue.AddJob(job); err != nil {
		return err
	}
	jobRequiresAuth = job
	if job.Repository.WorkStatus() != git.Fail {
		if err := jobQueue.RemoveFromQueue(job.Repository); err != nil {
			return err
		}
	}
	authenticationReturnView = returnViewName
	v, err := g.SetView(authenticationViewFeature.Name, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, keySymbol+selectionIndicator+red.Sprint(jobRequiresAuth.Repository.State.Remote.URL[0]))
	}
	g.Cursor = true
	if err := gui.openUserView(g); err != nil {
		return err
	}
	return gui.openPasswordView(g)
}

// close the opened auth views
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
	return gui.closeViewCleanup(authenticationReturnView)
}

// close the opened auth views and submit the credentials
func (gui *Gui) submitAuthenticationView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false

	// in order to read buffer of the views, first we need to find them
	vUser, err := g.View(authUserFeature.Name)
	if err != nil {
		return err // should return??
	}

	vPswd, err := g.View(authPasswordViewFeature.Name)
	if err != nil {
		return err // should return??
	}

	// the return string of the views contain trailing new lines
	re := regexp.MustCompile(`\r?\n`)
	creduser := re.ReplaceAllString(vUser.ViewBuffer(), "")
	credpswd := re.ReplaceAllString(vPswd.ViewBuffer(), "")

	// since the git ops require different types of options we better switch
	switch mode := jobRequiresAuth.JobType; mode {
	case job.FetchJob:
		jobRequiresAuth.Options = &command.FetchOptions{
			RemoteName: jobRequiresAuth.Repository.State.Remote.Name,
			Credentials: &git.Credentials{
				User:     creduser,
				Password: credpswd,
			},
		}
	case job.PullJob:
		// we handle pull as fetch&merge so same rule applies
		jobRequiresAuth.Options = &command.PullOptions{
			RemoteName: jobRequiresAuth.Repository.State.Remote.Name,
			Credentials: &git.Credentials{
				User:     creduser,
				Password: credpswd,
			},
		}
	}
	jobRequiresAuth.Repository.SetWorkStatus(git.Queued)

	// add this job to the last of the queue
	if err := gui.State.Queue.AddJob(jobRequiresAuth); err != nil {
		return err
	}

	return gui.closeAuthenticationView(g, v)
}

// open an error view to inform user with a message and a useful note
func (gui *Gui) openUserView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// first, create the label for user
	vlabel, err := g.SetView(authUserLabelFeature.Name, maxX/2-30, maxY/2-1, maxX/2-19, maxY/2+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(vlabel, authUserLabelFeature.Title)
		vlabel.Frame = false
	}
	// second, crete the user input
	v, err := g.SetView(authUserFeature.Name, maxX/2-18, maxY/2-1, maxX/2+29, maxY/2+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = authUserFeature.Title
		v.Editable = true
		v.Frame = false
	}
	return gui.focusToView(authUserFeature.Name)
}

// open an error view to inform user with a message and a useful note
func (gui *Gui) openPasswordView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// first, create the label for password
	vlabel, err := g.SetView(authPswdLabelViewFeature.Name, maxX/2-30, maxY/2, maxX/2-19, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(vlabel, authPswdLabelViewFeature.Title)
		vlabel.Frame = false
	}
	// second, crete the masked password input
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
