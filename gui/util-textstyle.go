package gui

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/isacikgoz/gitbatch/core/git"
	"github.com/isacikgoz/gitbatch/core/job"
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

	maxBranchLength     = 15
	maxRepositoryLength = 20
	hashLength          = 7

	ws       = " "
	pushable = string(blue.Sprint("↖"))
	pullable = string(blue.Sprint("↘"))
	dirty    = string(yellow.Sprint("✗"))

	queuedSymbol  = "•"
	workingSymbol = "•"
	successSymbol = "✔"
	pauseSymbol   = "॥"
	failSymbol    = "✗"

	fetchSymbol = "↓"
	pullSymbol  = "↓↳"
	mergeSymbol = "↳"

	keySymbol = ws + yellow.Sprint("🔑") + ws

	modeSeperator       = ""
	keyBindingSeperator = "░"

	selectionIndicator = ws + string(green.Sprint("→")) + ws
	tab                = ws
)

// this function handles the render and representation of the repository
// TODO: cleanup is required, right now it looks too complicated
func (gui *Gui) repositoryLabel(r *git.Repository) string {

	var prefix string
	b := r.State.Branch
	if b.Pushables != "?" {
		prefix = prefix + pushable + ws + b.Pushables +
			ws + pullable + ws + b.Pullables
	} else {
		prefix = prefix + pushable + ws + yellow.Sprint(b.Pushables) +
			ws + pullable + ws + yellow.Sprint(b.Pullables)
	}

	var repoName string
	sr := gui.getSelectedRepository()
	if sr == r {
		prefix = prefix + selectionIndicator
		repoName = green.Sprint(r.Name)
	} else {
		prefix = prefix + ws
		repoName = r.Name
	}
	// some branch names can be really long, in that times I hope the first
	// characters are important and meaningful
	branch := adjustTextLength(b.Name, maxBranchLength)
	prefix = prefix + string(cyan.Sprint(branch))

	if !b.Clean {
		prefix = prefix + ws + dirty + ws
	} else {
		prefix = prefix + ws
	}

	var suffix string
	// rendering the satus according to repository's state
	if r.WorkStatus() == git.Queued {
		if inQueue, j := gui.State.Queue.IsInTheQueue(r); inQueue {
			switch mode := j.JobType; mode {
			case job.FetchJob:
				suffix = blue.Sprint(queuedSymbol)
			case job.PullJob:
				suffix = magenta.Sprint(queuedSymbol)
			case job.MergeJob:
				suffix = cyan.Sprint(queuedSymbol)
			default:
				suffix = green.Sprint(queuedSymbol)
			}
		}
		return prefix + repoName + ws + suffix
	} else if r.WorkStatus() == git.Working {
		// TODO: maybe the type of the job can be written while its working?
		return prefix + repoName + ws + green.Sprint(workingSymbol)
	} else if r.WorkStatus() == git.Success {
		return prefix + repoName + ws + green.Sprint(successSymbol)
	} else if r.WorkStatus() == git.Paused {
		return prefix + repoName + ws + yellow.Sprint("authentication required (u)")
	} else if r.WorkStatus() == git.Fail {
		return prefix + repoName + ws + red.Sprint(failSymbol) + ws + red.Sprint(r.State.Message)
	}
	return prefix + repoName
}

func commitLabel(c *git.Commit) string {
	var body string
	switch c.CommitType {
	case git.EvenCommit:
		body = cyan.Sprint(c.Hash[:hashLength]) + " " + c.Message
	case git.LocalCommit:
		body = blue.Sprint(c.Hash[:hashLength]) + " " + c.Message
	case git.RemoteCommit:
		if len(c.Hash) > hashLength {
			body = yellow.Sprint(c.Hash[:hashLength]) + " " + c.Message
		} else {
			body = yellow.Sprint(c.Hash[:len(c.Hash)]) + " " + c.Message
		}
	default:
		body = c.Hash[:hashLength] + " " + c.Message
	}
	return body
}

// limit the text length for visual concerns
func adjustTextLength(text string, maxLength int) string {
	if len(text) > maxLength {
		return text[:maxLength-2] + ".."
	}
	return text
}

// colorize the plain diff text collected from system output
// the style is near to original diff command
func colorizeDiff(original string) (colorized []string) {
	colorized = strings.Split(original, "\n")
	re := regexp.MustCompile(`@@ .+ @@`)
	for i, line := range colorized {
		if len(line) > 0 {
			if line[0] == '-' {
				colorized[i] = red.Sprint(line)
			} else if line[0] == '+' {
				colorized[i] = green.Sprint(line)
			} else if re.MatchString(line) {
				s := re.FindString(line)
				colorized[i] = cyan.Sprint(s) + line[len(s):]
			} else {
				continue
			}
		} else {
			continue
		}
	}
	return colorized
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

	// separate the protocol and remote link
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
