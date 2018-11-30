package gui

import (
	"regexp"

	"github.com/fatih/color"
	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/isacikgoz/gitbatch/pkg/job"
)

var (
	black   = color.New(color.FgBlack)
	blue    = color.New(color.FgBlue)
	green   = color.New(color.FgGreen)
	red     = color.New(color.FgRed)
	cyan    = color.New(color.FgCyan)
	yellow  = color.New(color.FgYellow)
	white   = color.New(color.FgWhite)
	magenta = color.New(color.FgMagenta)

	bold = color.New(color.Bold)

	maxBranchLength = 15
	maxRepositoryLength = 20
	hashLength = 7

	ws = " "
	pushable = string(blue.Sprint("↖"))
	pullable = string(blue.Sprint("↘"))
	confidentArrow = string(magenta.Sprint("→"))
	unconfidentArrow = string(yellow.Sprint("→"))
	dirty = string(yellow.Sprint("✗"))
	unkown = magenta.Sprint("?")

	queuedSymbol = "•"
	workingSymbol = "•"
	successSymbol = "✔"
	failSymbol = "✗"

	fetchSymbol = "↓"
	pullSymbol = "↓↳"
	mergeSymbol = "↳"

	modeSeperator = ""
	keyBindingSeperator = "░"

	selectionIndicator = string(green.Sprint("→")) + ws
	tab = ws + ws
)

// this fucntion handles the render and representation of the repository
// TODO: cleanup is required, right now it looks too complicated
func (gui *Gui) displayString(entity *git.RepoEntity) string {
	suffix := ""
	prefix := ""

	if entity.Branch.Pushables != "?" {
		prefix = prefix + pushable + ws + entity.Branch.Pushables + ws +
			pullable + ws + entity.Branch.Pullables + ws + confidentArrow + ws
	} else {
		prefix = prefix + pushable + ws + yellow.Sprint(entity.Branch.Pushables) + ws +
			pullable + ws + yellow.Sprint(entity.Branch.Pullables) + ws + unconfidentArrow + ws
	}

	// some branch names can be really long, in that times I hope the first
	// characters are important and meaningful
	branch := adjustTextLength(entity.Branch.Name, maxBranchLength)
	prefix = prefix + string(cyan.Sprint(branch))

	if !entity.Branch.Clean {
		prefix = prefix + ws + dirty + ws 
	} else {
		prefix = prefix + ws 
	}

	// rendering the satus according to repository's state
	if entity.State == git.Queued {
		if inQueue, ty := gui.State.Queue.IsInTheQueue(entity); inQueue {
		switch mode := ty; mode {
			case job.Fetch:
				suffix = blue.Sprint(queuedSymbol)
			case job.Pull:
				suffix = magenta.Sprint(queuedSymbol)
			case job.Merge:
				suffix = cyan.Sprint(queuedSymbol)
			default:
				suffix = green.Sprint(queuedSymbol)
			}
		}
		return prefix + entity.Name + ws + suffix
	} else if entity.State == git.Working {
		// TODO: maybe the type of the job can be written while its working?
		return prefix + entity.Name + ws + green.Sprint(workingSymbol)
	} else if entity.State == git.Success {
		return prefix + entity.Name + ws + green.Sprint(successSymbol)
	} else if entity.State == git.Fail {
		return prefix + entity.Name + ws + red.Sprint(failSymbol)
	} else {
		return prefix + entity.Name
	}
}

// limit the text length for visual concerns
func adjustTextLength(text string, maxLength int) (adjusted string) {
	if len(text) > maxLength {
		adjusted := text[:maxLength-2] + ".."
		return adjusted
	} else {
		return text
	}
}

// the remote link can be too verbose sometimes, so it is good to trim it
func trimRemoteURL(url string) (urltype string, shorturl string) {
	// lets trim the unnecessary .git extension of the url
	regit := regexp.MustCompile(`.git`)
	if regit.MatchString(url[len(url)-4:]) {
		url = url[:len(url)-4]
	}

	// find out the protocol
	ressh := regexp.MustCompile(`git@`)
	rehttp := regexp.MustCompile(`http://`)
	rehttps := regexp.MustCompile(`https://`)

	// seperate the protocol and remote link
	if ressh.MatchString(url) {
		shorturl = ressh.Split(url, 5)[1]
		urltype = "ssh"
	} else if rehttp.MatchString(url) {
		shorturl = rehttp.Split(url, 5)[1]
		urltype = "http"
	} else if rehttps.MatchString(url) {
		shorturl = rehttps.Split(url, 5)[1]
		urltype = "https"
	}
	return urltype, shorturl
}
