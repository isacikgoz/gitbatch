[![Build Status](https://travis-ci.com/isacikgoz/gitbatch.svg?branch=master)](https://travis-ci.com/isacikgoz/gitbatch) [![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/isacikgoz/gitbatch)](https://goreportcard.com/report/github.com/isacikgoz/gitbatch)

## gitbatch
Aim of this tool to make your local repositories synchronized with remotes easily. I inspired from lazygit and decided to build this project for my needs; Since my daily work is tied to many repositories I often end up walking on many directories and manually pulling updates etc. To make this routine faster, I created a simple tool to handle this job. I really enjoy working on this project and I hope it will be a useful tool.

**Disclaimer**
- Authentication required repositories are **not supported** so using ssh is recommended if you need to authenticate to fetch/pull
  - [Connecting to GitHub with SSH](https://help.github.com/articles/connecting-to-github-with-ssh/)
  - [GitLab and SSH keys](https://docs.gitlab.com/ee/ssh/)
  - [BitBucket Set up an SSH key](https://confluence.atlassian.com/bitbucket/set-up-ssh-for-git-728138079.html)
- Feedbacks are welcome. For now, known issues are:
  - At very low probability app fails to load repositories, try again it will load next time (multithreading problem)
  - Sometimes when you scroll too fast while pulling/fetching/merging, some multithreading problem occurs and app crashes (will fix soon)
  - colors vary to your terminal theme colors, so if the contrast is not enough on some color decisions; discussions are welcome

Here is the screencast of the app:
[![asciicast](https://asciinema.org/a/B4heYReiNgqwUbWL2RYnTzt5H.svg)](https://asciinema.org/a/B4heYReiNgqwUbWL2RYnTzt5H)

## installation
for now, installation requires golang compiler and minimum golang 1.10 is recommended.
- if you don't have golang installed refer to [golang.org](https://golang.org/dl/)
- you should have $GOPATH env variable set and your $PATH should include $GOPATH/bin to run app from anywhere

to install run the command;
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
- [Controls page](https://github.com/isacikgoz/gitbatch/wiki/Controls)
- [Modes page](https://github.com/isacikgoz/gitbatch/wiki/Modes)
- [Display page](https://github.com/isacikgoz/gitbatch/wiki/Display)

## Further goals
- add testing
- select all feature ✔
- arrange repositories to an order e.g. alphabetic, last modified, etc. ✔
- shift keys, i.e. **s** for iterate **alt + s** for reverse iteration ✔
- recursive repository search from the filesystem
- full src-d/go-git integration (*having some performance issues*)
- implement config file to pre-define repo locations or some settings
- resolve authentication issues

## Credits
- [go-git](https://github.com/src-d/go-git) for git interface
- [gocui](https://github.com/jroimartin/gocui) for user interface
- [logrus](https://github.com/sirupsen/logrus) for logging
- [lazygit](https://github.com/jesseduffield/lazygit) as app template
- [color](https://github.com/fatih/color) for colored text
- [kingpin](https://github.com/alecthomas/kingpin) for command-line flag&options
