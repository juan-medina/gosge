# gosge
Go Simple Game Engine using an [ECS](https://github.com/juan-medina/goecs)

[![License: Apache2](https://img.shields.io/badge/license-Apache%202-blue.svg)](/LICENSE)
[![go version](https://img.shields.io/github/v/tag/juan-medina/gosge?label=version)](https://pkg.go.dev/mod/github.com/juan-medina/gosge)
[![godoc](https://godoc.org/github.com/juan-medina/gosge?status.svg)](https://pkg.go.dev/mod/github.com/juan-medina/gosge)
[![Build Status](https://travis-ci.com/juan-medina/gosge.svg?branch=main)](https://travis-ci.com/juan-medina/gosge)
[![conduct](https://img.shields.io/badge/code%20of%20conduct-contributor%20covenant%202.0-purple.svg?style=flat-square)](https://www.contributor-covenant.org/version/2/0/code_of_conduct/)

## Info

gosge is an opinionated 2D only game engine the uses [GOECS](https://github.com/juan-medina/goecs) for _easily_ develop games
with an ECS paradigm.

Internally uses the go [port](https://github.com/gen2brain/raylib-go) of [raylib](https://www.raylib.com/) for most of the device functionalities, including rendering.

## Simple Hello World

```go
package main

import (
	"github.com/juan-medina/gosge"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/options"
	"log"
)

// game options
var opt = options.Options{
	Title:      "Hello Game",
	BackGround: color.Black,
}

const (
	fontName = "resources/go_regular.fnt"
	fontSize = 100
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
)

// Simple Usage
func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

func loadGame(eng *gosge.Engine) error {
	// Preload font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	// get the ECS world
	world := eng.World()

	// gameScale from the real screen size to our design resolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// add the centered text
	world.AddEntity(
		ui.Text{
			String:     "Hello World",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.MiddleVAlignment,
			Font:       fontName,
			Size:       fontSize * gameScale.Min,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		color.White,
	)
	return nil
}
```

## Run examples

The examples are available on [this folder](/examples), and can be run using make:

```bash
$ make run example=hello
$ make run example=eyes
$ make run example=layers
$ make run example=stages
$ make run example=animation
$ make run example=audio
$ make run example=tiled
```

Alternatively you could run them with :

```bash
$ go run examples/hello/hello.go
$ go run examples/eyes/eyes.go
$ go run examples/layers/layers.go
$ go run examples/stages/stages.go
$ go run examples/animation/animation.go
$ go run examples/audio/audio.go
$ go run examples/tiled/tiled.go
```

## Requirements

### Ubuntu

#### X11

    apt-get install libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev

#### Wayland

    apt-get install libgl1-mesa-dev libwayland-dev libxkbcommon-dev

### Fedora

#### X11

    dnf install mesa-libGL-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel

#### Wayland

    dnf install mesa-libGL-devel wayland-devel libxkbcommon-devel

### macOS

On macOS, you need Xcode or Command Line Tools for Xcode.

### Windows

On Windows, you need C compiler, like [Mingw-w64](https://mingw-w64.org) or [TDM-GCC](http://tdm-gcc.tdragon.net/).
You can also build binary in [MSYS2](https://msys2.github.io/) shell.

## Installation

```bash
go get -v -u github.com/juan-medina/gosge
```

## Build Tags

- `opengl21` : uses OpenGL 2.1 backend (default is 3.3)
- `wayland` : builds against Wayland libraries

## Examples Resources
- Gopher Graphics
    - https://awesomeopensource.com/project/egonelbre/gophers
- Game art 2D:
    - https://www.gameart2d.com
-  Mobile Game Graphics
    - https://mobilegamegraphics.com/product/free-parallax-backgrounds
- Of Far Different Nature Music Loops
    - https://fardifferent.itch.io/loops
- freesound.org
    - https://freesound.org/people/SKKreativ/sounds/456255/
