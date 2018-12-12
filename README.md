[![Build Status](https://travis-ci.com/isacikgoz/gitbatch.svg?branch=master)](https://travis-ci.com/isacikgoz/gitbatch) [![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/isacikgoz/gitbatch)](https://goreportcard.com/report/github.com/isacikgoz/gitbatch)

## gitbatch
This tool is beening built to make your local repositories synchronized with remotes easily. Although the focus is batch jobs, you can still do de facto micro management of your git repositories (e.g *add/reset, stash, commit etc.*)

Here is the screencast of the app:
[![asciicast](https://asciinema.org/a/eXgXpzZfuHxMpZqGMVODUipyc.svg)](https://asciinema.org/a/eXgXpzZfuHxMpZqGMVODUipyc)

## Installation
For now, installation requires golang compiler and minimum golang 1.10 is recommended. (binary distribution will be provided on minimum viable product)
- If you don't have golang installed refer to [golang.org](https://golang.org/dl/).
- You should have $GOPATH env variable set and your $PATH should include $GOPATH/bin to run app from anywhere.

To install run the following command;
```bash
go get -u github.com/isacikgoz/gitbatch
```
I prefer using gopath like this in my zshrc or bashrc;
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## Use
run the `gitbatch` command from the parent of your git repositories. For start-up options simply `gitbatch --help`

For more information;
- To see keybindings refer to [Controls page](https://github.com/isacikgoz/gitbatch/wiki/Controls)
- Learn how to config app at [Configuration page](https://github.com/isacikgoz/gitbatch/wiki/Configuration)
- Wonder what mode does what? see [Modes page](https://github.com/isacikgoz/gitbatch/wiki/Modes)
- What are those arrows, colors etc mean? Answer is here at [Display page](https://github.com/isacikgoz/gitbatch/wiki/Display)

## Further goals
- add testing
- full src-d/go-git integration (*having some performance issues*)
- add commit and maybe push?

## Known issues
Please refer to [Known issues page](https://github.com/isacikgoz/gitbatch/wiki/Known-issues)

## Credits
- [go-git](https://github.com/src-d/go-git) for git interface (partially)
- [gocui](https://github.com/jroimartin/gocui) for user interface
- [logrus](https://github.com/sirupsen/logrus) for logging
- [viper](https://github.com/spf13/viper) for configuration management
- [color](https://github.com/fatih/color) for colored text
- [lazygit](https://github.com/jesseduffield/lazygit) as app template and reference
- [kingpin](https://github.com/alecthomas/kingpin) for command-line flag&options

I love [lazygit](https://github.com/jesseduffield/lazygit), with that inspiration, decided to build this project to be even more lazy. The rationale was; my daily work is tied to many repositories and I often end up walking on many directories and manually pulling updates etc. To make this routine faster, I created a simple tool to handle this job. I really enjoy working on this project and I hope it will be a useful tool.
