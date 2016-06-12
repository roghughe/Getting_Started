Getting Started With Go
=======================

This project contains a few notes on getting started with GO, based around setting up a development environment 
using eclipse - which just happends to be my favourite.

1) This has been done using a new eclipse install. Therefore install java for eclipse. 
   I put it in /users/Roger/eclipse-Go
2) Install Go. Download and install the go libraries as the default unix location of /usr/local/go from
   https://golang.org/dl/
3) Set up a GPPATH variable: /users/Roger/eclipse-Go/GOPATH This was created here instead of the default locatiobn
   to keep track of what's being downloaded.
4) Install the goclipse IDE plugin: https://github.com/GoClipse/goclipse Install by selecting Install New Software and 
   creating an entry for: Go - http://goclipse.github.io/releases/. When selecting options, select the Go and basic CDT 
   options (selecting all causes the install to fail).
5) Make sure you have the git commandline tools installed.


Handy hints can be found at:

http://bark4mark.blogspot.co.uk/2015/11/getting-going-with-go.html