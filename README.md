# gojira2d

[![Build Status](https://travis-ci.org/maxfish/gojira2d.svg?branch=master)](https://travis-ci.org/maxfish/gojira2d)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxfish/gojira2d)](https://goreportcard.com/report/github.com/maxfish/gojira2d)
[![Coverage Status](https://coveralls.io/repos/github/maxfish/gojira2d/badge.svg?branch=master)](https://coveralls.io/github/maxfish/gojira2d?branch=master)
[![Join the chat at https://gitter.im/gojira2d/Lobby](https://badges.gitter.im/gojira2d/Lobby.svg)](https://gitter.im/gojira2d/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Simple 2D game library written in Go and based on modern OpenGL.

Gojira2D is licensed under the terms of the MIT License. See LICENSE for details of the usage license granted to you for this code.

## HIGHLIGHTS

* It uses modern OpenGL - Core Profile 4.1
* Super easy setup: create a window and a game loop in less than 10 lines of code
* Shader support: Vertex, Geometry, Fragment
* Input handling:
  * Basic support for keyboard and mouse
  * Joystick support, including emulation via keyboard
* Basic support for shapes: lines, polylines and [approximated] circles
* Fonts support:
  * Bitmap fonts in [BMFont format](http://www.angelcode.com/products/bmfont/doc/file_format.html)
  * Distance field rendering of TTF fonts
* Physics support:
  * 2D rigid body physics via Box2D
  * [R.U.B.E](https://www.iforce2d.net/rube/) ([JSON format](https://www.iforce2d.net/rube/json-structure)) scene loader
* Developed and tested on MacOS

## Dependencies

* [GLFW](https://github.com/go-gl/glfw) as window manager
* [MathGL](https://github.com/go-gl/mathgl) as math library
* [Box2D](https://github.com/ByteArena/box2d) for 2D physics
* [B2DJson](https://github.com/maxfish/go-b2dJson) for loading R.U.B.E scenes

## Installation

Install Golang and GLFW:

    $ brew install go dep glfw

Setup your [`$GOPATH`](https://golang.org/doc/code.html#GOPATH) and clone the
repository into `$GOPATH/src` folder:

    $ go get -u maxfish/gojira2d
    $ cd $GOPATH/src/maxfish/gojira2d
    $ git remote set-url origin git@github.com:maxfish/gojira2d.git

Use `dep` to fetch dependencies:

    $ dep ensure

Try running some examples:

    $ go run examples/quad/main.go
    ...
