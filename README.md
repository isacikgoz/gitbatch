[![Build Status](https://travis-ci.com/isacikgoz/gitbatch.svg?branch=master)](https://travis-ci.com/isacikgoz/gitbatch) [![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/isacikgoz/gitbatch)](https://goreportcard.com/report/github.com/isacikgoz/gitbatch)

## gitbatch
Aim of this tool to make your local repositories synchronized with remotes easily. Since my daily work is tied to many repositories I often end up walking on many directories and manually pulling updates etc. To make this routine more elegant, I made a decision to create a simple tool to handle this job in a seamless way. Actually I am not a golang developer but I thought it would be cool to create this tool with a new language for me. While developing this project, I am getting experience with golang. As a result, I really enjoy working on this project and I hope it will be a useful tool. I hope this tool can be useful for others too.

**Disclamier**
- This is still a work in progress project.
- Authentication reuqired repostires are **NOT SUPPORTED** using ssh is recommended if you need to authenticate to fetch/pull
  - [Connecting to GitHub with SSH](https://help.github.com/articles/connecting-to-github-with-ssh/)
  - [GitLab and SSH keys](https://docs.gitlab.com/ee/ssh/)
  - [BitBucket Set up an SSH key](https://confluence.atlassian.com/bitbucket/set-up-ssh-for-git-728138079.html)
- This project is not widely tested in various environments. (macOS and Ubuntu is okay but didn't tried on Windows)

Here is the initial look of the project:
[![asciicast](https://asciinema.org/a/qmZDhmUjwWmZvdpZGYIRW56h7.svg)](https://asciinema.org/a/qmZDhmUjwWmZvdpZGYIRW56h7)

## installation
the project is at very early version but it can be tested;
```bash
go get github.com/isacikgoz/gitbatch
```

## Use
run the command the parent of your git repositories. For start-up options simply:
`gitbatch --help`

### Controls

- **tab**: Switch mode
- **↑** or **k**: Up
- **↓** or **j**: Down
- **b**: Iterate over branches
- **r**: Iterate over remotes
- **e**: Iterate over remote branches
- **s**: Iterate over commits
- **d**: Show commit diff
- **c**: Controls or close windows if a pop-up window opened
- **enter**: Start queue
- **space**: Add to queue
- **ctrl + c**: Force application to quit
- **q**: Quit

### Modes

- **FETCH**: fetches the selected **remote** e.g. *Origin*, *Upstream*, etc.
- **PULL**: fetches the selected **remote** and merges selected **remote branch** into **active branch** e.g. origin/master → master
- **MERGE**: merges the selected **remote branch** into **active branch** e.g. origin/master → master

### Display

#### Repository Screen
↖ 0 ↘ 0 → master ✗ ips.server.slave.native •  
↖ (pushables) ↘ (pullables) → (branch) (✗ if dirty) (repository folder name) (• if queued)

- if pushables or pullables appear to be "**?**", that means no upstream configured for the active branch
- the queued indicator color represents the operation same as mode color

#### Commit Screen
- if hash color is cyan it means that commit is a local commit, if yellow it means it is a commit that will merge in to your active branch if you pull or merge
- you can see the diff by simply pressing **d** on the selected commit

## Further goals
- add testing, currently this is the M.V.P. so that people can swiftly get hands on
- select all feature
- arrange repositories to an order e.g. alphabetic, last modified, etc.
- shift keys, i.e. **s** for iterate **shift + s** for reverse iteration
- binary distrubiton over [homebrew](https://github.com/Homebrew/brew) or similar
- recursive repository search from the filesystem
- full src-d/go-git integration (*waiting for merge abilities*)
- implement config file to pre-define repo locations or some settings
- resolve authentication issues
- *maybe* handle conflicts with [fac](https://github.com/mkchoi212/fac) integration

## Credits
- [go-git](https://github.com/src-d/go-git) for git interface
- [gocui](https://github.com/jroimartin/gocui) for user interface
- [color](https://github.com/fatih/color) for colored text
- [kingpin](https://github.com/alecthomas/kingpin) for command-line flag&options
