[![Build Status](https://travis-ci.com/isacikgoz/gitbatch.svg?branch=master)](https://travis-ci.com/isacikgoz/gitbatch) [![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/isacikgoz/gitbatch)](https://goreportcard.com/report/github.com/isacikgoz/gitbatch)

## gitbatch
Aim of this tool to make your local repositories synchronized with remotes easily. Since my daily work is tied to many repositories I often end up walking on many directories and manually pulling updates etc. To make this routine more elegant, I made a decision to create a simple tool to handle this job in a seamless way. Actually I am not a golang developer but I thought it would be cool to create this tool with a new language for me. While developing this project, I am getting experience with golang. As a result, I really enjoy working on this project and I hope it will be a useful tool. I hope this tool can be useful for others too. The satisfaction of saving someones precious time is priceless.

**Disclamier** This is still a work in progress project.

Here is the initial look of the project: 
[![asciicast](https://asciinema.org/a/eYfR9eWC4VGjiAyBUE7hpaZph.svg)](https://asciinema.org/a/eYfR9eWC4VGjiAyBUE7hpaZph)

### Use
run the command the parent of your git repositories. Or simply:
`gitbatch --help`

## installation
the project is at very very early version but if you like new adventures;
```
go get github.com/isacikgoz/gitbatch
cd $GOPATH/src/github.com/isacikgoz/gitbatch
go run main.go ".." # or simply go build && mv gitbatch $GOPATH/bin or any where you use as path
```

## Further goals
- recursive repository search from the filesystem
- full src-d/go-git integration (*waiting for merge abilities*)
- implement modal base ux like vim (fetch/merge *maybe* even push)
- resolve authentication issues
- *maybe* handle conflicts with [fac](https://github.com/mkchoi212/fac) integration

## Credits
- [go-git](https://github.com/src-d/go-git) for git interface
- [gocui](https://github.com/jroimartin/gocui) for user interface
- [color](https://github.com/fatih/color) for colored text
- [lazygit](https://github.com/jesseduffield/lazygit) for reference
- [kingpin](https://github.com/alecthomas/kingpin) for command-line flag&options
