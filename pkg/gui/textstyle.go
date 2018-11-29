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

	branch := adjustTextLength(entity.Branch.Name, maxBranchLength)
	prefix = prefix + string(cyan.Sprint(branch))

	if !entity.Branch.Clean {
		prefix = prefix + ws + dirty + ws 
	} else {
		prefix = prefix + ws 
	}

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
		return prefix + entity.Name + ws + green.Sprint(workingSymbol)
	} else if entity.State == git.Success {
		return prefix + entity.Name + ws + green.Sprint(successSymbol)
	} else if entity.State == git.Fail {
		return prefix + entity.Name + ws + red.Sprint(failSymbol)
	} else {
		return prefix + entity.Name
	}
}

func adjustTextLength(text string, maxLength int) (adjusted string) {
	if len(text) > maxLength {
		adjusted := text[:maxLength-2] + ".."
		return adjusted
	} else {
		return text
	}
}

func trimRemoteURL(url string) (urltype string, shorturl string) {
	regit := regexp.MustCompile(`.git`)
	if regit.MatchString(url[len(url)-4:]) {
		url = url[:len(url)-4]
	}
	ressh := regexp.MustCompile(`git@`)
	rehttp := regexp.MustCompile(`http://`)
	rehttps := regexp.MustCompile(`https://`)

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
