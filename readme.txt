Getting Started With Go
=======================

This project contains a few notes on getting started with GO, based around setting up a development environment 
using eclipse - which used to be favourite Golang IDE.

1) This has been done using a new eclipse install. Therefore install java for eclipse. 
   I put it in /users/Roger/eclipse-Go
2) Install Go. Download and install the go libraries as the default unix location of /usr/local/go from
   https://golang.org/dl/
3) Set up a GPPATH variable: /users/Roger/eclipse-Go/GOPATH This was created here instead of the default location
   to keep track of what's being downloaded.
4) Install the goclipse IDE plugin: https://github.com/GoClipse/goclipse Install by selecting Install New Software and 
   creating an entry for: Go - http://goclipse.github.io/releases/. When selecting options, select the Go and basic CDT 
   options (selecting all causes the install to fail).
5) Make sure you have the git command line tools installed as the go command uses them.
6) Make sure you have the xcode command line tools available: http://osxdaily.com/2014/02/12/install-command-line-tools-mac-os-x/
7) Install the gocode Go completion tool: "go get -u github.com/nsf/gocode"
8) Configure the eclipse preferences for go - this is fiddly. 



Handy hints can be found at:

http://bark4mark.blogspot.co.uk/2015/11/getting-going-with-go.html

Eclipse Go Preferences
----------------------

Main:

GOROOT: /usr/local/go
GOPATH: /Users/Roger/eclipse-Go/GOPATH
-- Also ensure that the project is part of the GOPATH

Tools:

gocode: /Users/Roger/eclipse-Go/GOPATH/bin/gocode

