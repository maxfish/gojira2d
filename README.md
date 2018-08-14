# gojira2d

[![Build Status](https://travis-ci.org/maxfish/gojira2d.svg?branch=master)](https://travis-ci.org/maxfish/gojira2d)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxfish/gojira2d)](https://goreportcard.com/report/github.com/maxfish/gojira2d)
[![Join the chat at https://gitter.im/gojira2d/Lobby](https://badges.gitter.im/gojira2d/Lobby.svg)](https://gitter.im/gojira2d/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Simple 2D game library based on modern OpenGL

## Installation

Install Golang and GLFW:

    $ brew install go glfw

Setup your [`$GOPATH`](https://golang.org/doc/code.html#GOPATH) and clone the
repository into `$GOPATH/src` folder:

    $ go get -u maxfish/gojira2d
    $ cd $GOPATH/src/maxfish/gojira2d
    $ git remote set-url origin git@github.com:maxfish/gojira2d.git

Install Golang [`dep tool`](https://github.com/golang/dep) and use it to fetch dependencies:

    $ dep ensure

Try running some examples:

    $ go run examples/quad/main.go
    ...
