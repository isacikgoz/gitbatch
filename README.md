[![Build Status](https://travis-ci.com/isacikgoz/gitbatch.svg?branch=master)](https://travis-ci.com/isacikgoz/gitbatch) [![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/isacikgoz/gitbatch)](https://goreportcard.com/report/github.com/isacikgoz/gitbatch)

## gitbatch
Managing multiple git repositories is easier than ever. I (*was*) often end up walking on many directories and manually pulling updates etc. To make this routine faster, I created a simple tool to handle this job. Although the focus is batch jobs, you can still do de facto micro management of your git repositories (e.g *add/reset, stash, commit etc.*)

Check out the screencast of the app:
[![asciicast](https://asciinema.org/a/lxoZT6Z8fSliIEebWSPVIY8ct.svg)](https://asciinema.org/a/lxoZT6Z8fSliIEebWSPVIY8ct)

## Installation
To install with go, run the following command;
```bash
go get -u github.com/isacikgoz/gitbatch
```
For other options see [installation page](https://github.com/isacikgoz/gitbatch/wiki/Installation)

## Use
run the `gitbatch` command from the parent of your git repositories. For start-up options simply `gitbatch --help`

For more information see the [wiki pages](https://github.com/isacikgoz/gitbatch/wiki)

## Further goals
- improve testing
- add push
- full src-d/go-git integration (*having some performance issues in large repos*)
  - fetch, config, rev-list, add, reset, commit, status and diff commands are supported but not fully utilized, still using git occasionally
  - merge, stash are not supported yet by go-git

## Credits
- [go-git](https://github.com/src-d/go-git) for git interface (partially)
- [gocui](https://github.com/jroimartin/gocui) for user interface
- [logrus](https://github.com/sirupsen/logrus) for logging
- [viper](https://github.com/spf13/viper) for configuration management
- [color](https://github.com/fatih/color) for colored text
- [lazygit](https://github.com/jesseduffield/lazygit) for inspiration
- [kingpin](https://github.com/alecthomas/kingpin) for command-line flag&options

