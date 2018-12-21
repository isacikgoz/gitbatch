[![Build Status](https://travis-ci.com/isacikgoz/gitbatch.svg?branch=master)](https://travis-ci.com/isacikgoz/gitbatch) [![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/isacikgoz/gitbatch)](https://goreportcard.com/report/github.com/isacikgoz/gitbatch)

## gitbatch
This tool is being built to make your local repositories synchronized with remotes easily. Although the focus is batch jobs, you can still do de facto micro management of your git repositories (e.g *add/reset, stash, commit etc.*)

Here is the screencast of the app:
[![asciicast](https://asciinema.org/a/QQPVDWVxUR3bvJhIZY3c4PTuG.svg)](https://asciinema.org/a/QQPVDWVxUR3bvJhIZY3c4PTuG)

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
- add testing
- add push
- full src-d/go-git integration (*having some performance issues in such cases*)
  - fetch, config, add, reset, commit, status and diff commands are supported but not fully utilized, still using git sometimes
  - merge, rev-list, stash are not supported yet by go-git

## Known issues
Please refer to [Known issues page](https://github.com/isacikgoz/gitbatch/wiki/Known-issues) and feel free to open an issue if you encounter with a problem.

## Credits
- [go-git](https://github.com/src-d/go-git) for git interface (partially)
- [gocui](https://github.com/jroimartin/gocui) for user interface
- [logrus](https://github.com/sirupsen/logrus) for logging
- [viper](https://github.com/spf13/viper) for configuration management
- [color](https://github.com/fatih/color) for colored text
- [lazygit](https://github.com/jesseduffield/lazygit) as app template and reference
- [kingpin](https://github.com/alecthomas/kingpin) for command-line flag&options

I love [lazygit](https://github.com/jesseduffield/lazygit), with that inspiration, decided to build this project to be even more lazy. The rationale was; my daily work is tied to many repositories and I often end up walking on many directories and manually pulling updates etc. To make this routine faster, I created a simple tool to handle this job. I really enjoy working on this project and I hope it will be a useful tool.
