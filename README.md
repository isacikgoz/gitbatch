[![Build Status](https://travis-ci.com/isacikgoz/gitbatch.svg?branch=master)](https://travis-ci.com/isacikgoz/gitbatch) [![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/isacikgoz/gitbatch)](https://goreportcard.com/report/github.com/isacikgoz/gitbatch)

## gitbatch
Aim of this simple application to make your local repositories syncrhonized with remotes easily. 

**Disclamier** This is still a work in progress project.

Here is the intial look of the project: 
[![asciicast](https://asciinema.org/a/eYfR9eWC4VGjiAyBUE7hpaZph.svg)](https://asciinema.org/a/eYfR9eWC4VGjiAyBUE7hpaZph)

### Use
run the command the parent of your git repositories. Or simply:
`gitbatch --help`

## installation
installation guide will be provided after in-house tests and after implementation of the unit tests just get less headache. And maybe later I can distribute binaries from releases page.

## Further goals
- full src-d/go-git integration
- implement modal base ux like vim (fetch/pull maybe even push)
- Resolve authentication issues
- Handle conflicts

## Credits
[go-git](https://github.com/src-d/go-git)
[gocui](https://github.com/jroimartin/gocui)
[color](https://github.com/fatih/color)
[lazygit](https://github.com/jesseduffield/lazygit)