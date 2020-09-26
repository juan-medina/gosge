/*
 * Copyright (c) 2020 Juan Medina.
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documgopher1tion files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

package main

import (
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "GOSGE Collision Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed:   true,
	// Width:      2048,
	// Height:     1536,
}

const (
	fontName    = "resources/go_regular.fnt"
	fontSmall   = 60
	gopherSpeed = 500
	spriteScale = 0.25
)

type demoFactors struct {
	factor1 geometry.Point
	factor2 geometry.Point
}

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	gEng             *gosge.Engine
	gopher1          *goecs.Entity
	gopher2          *goecs.Entity
	move             geometry.Point
	gameScale        geometry.Scale
	area1            *goecs.Entity
	area2            *goecs.Entity
	spriteSize       geometry.Size
	currentFactor    int
	factors          = []demoFactors{
		{
			factor1: geometry.Point{X: 1, Y: 1},
			factor2: geometry.Point{X: 1, Y: 1},
		},
		{
			factor1: geometry.Point{X: 0.5, Y: 1},
			factor2: geometry.Point{X: 1, Y: 1},
		},
		{
			factor1: geometry.Point{X: 0.5, Y: 1},
			factor2: geometry.Point{X: 0.5, Y: 1},
		},
		{
			factor1: geometry.Point{X: 0.5, Y: 0.5},
			factor2: geometry.Point{X: 0.5, Y: 1},
		},
		{
			factor1: geometry.Point{X: 0.5, Y: 0.5},
			factor2: geometry.Point{X: 0.5, Y: 0.5},
		},
		{
			factor1: geometry.Point{X: 1, Y: 0.5},
			factor2: geometry.Point{X: 1, Y: 1},
		},
		{
			factor1: geometry.Point{X: 1, Y: 0.5},
			factor2: geometry.Point{X: 1, Y: 0.5},
		},
	}
)

func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng *gosge.Engine) error {
	var err error
	gEng = eng

	// Preload font
	if err = eng.LoadFont(fontName); err != nil {
		return err
	}

	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale = eng.GetScreenSize().CalculateScale(designResolution)

	// preload sprite sheet
	if err = eng.LoadSpriteSheet("resources/gamer.json"); err != nil {
		return err
	}

	// get sprite size
	if spriteSize, err = eng.GetSpriteSize("resources/gamer.json", "gamer.png"); err != nil {
		return err
	}

	gopherPos := geometry.Point{
		X: designResolution.Width / 4 * gameScale.Point.X,
		Y: designResolution.Height / 2 * gameScale.Point.Y,
	}

	areaPos := geometry.Point{
		X: gopherPos.X - (spriteSize.Width * spriteScale / 2 * gameScale.Max),
		Y: gopherPos.Y - (spriteSize.Height * spriteScale / 2 * gameScale.Max),
	}

	// add our moving gopher
	gopher1 = world.AddEntity(
		sprite.Sprite{
			Sheet: "resources/gamer.json",
			Name:  "gamer.png",
			Scale: gameScale.Max * spriteScale,
			FlipX: false,
			FlipY: false},
		gopherPos,
	)

	area1 = world.AddEntity(
		shapes.Box{
			Size:      spriteSize,
			Scale:     gameScale.Max * spriteScale,
			Thickness: 1,
		},
		areaPos,
		color.Red,
	)

	gopherPos = geometry.Point{
		X: designResolution.Width / 2 * gameScale.Point.X,
		Y: designResolution.Height / 2 * gameScale.Point.Y,
	}

	areaPos = geometry.Point{
		X: gopherPos.X - (spriteSize.Width * spriteScale / 2 * gameScale.Max),
		Y: gopherPos.Y - (spriteSize.Height * spriteScale / 2 * gameScale.Max),
	}

	// add an static gopher
	gopher2 = world.AddEntity(
		sprite.Sprite{
			Sheet: "resources/gamer.json",
			Name:  "gamer.png",
			Scale: gameScale.Max * spriteScale,
			FlipX: false,
			FlipY: false},
		gopherPos,
	)

	area2 = world.AddEntity(
		shapes.Box{
			Size:      spriteSize,
			Scale:     gameScale.Max * spriteScale,
			Thickness: 1,
		},
		areaPos,
		color.Red,
	)

	currentFactor = 0

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     "move with cursors, space change collision, press <ESC> to close",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height * gameScale.Point.Y,
		},
		effects.AlternateColor{
			Time: .25,
			From: color.White,
			To:   color.White.Alpha(0),
		},
	)

	// add the collide system
	world.AddSystem(collideSystem)

	// add the movement system
	world.AddSystem(moveSystem)

	// update areas system
	world.AddSystem(updateAreasSystem)

	// add the key listener
	world.AddListener(keyListener)

	return nil
}

func updateAreasSystem(_ *goecs.World, _ float32) error {
	updateArea(gopher1, area1, factors[currentFactor].factor1)
	updateArea(gopher2, area2, factors[currentFactor].factor2)
	return nil
}

func updateArea(sprite *goecs.Entity, area *goecs.Entity, factor geometry.Point) {
	spritePos := geometry.Get.Point(sprite)
	size := geometry.Size{
		Width:  spriteSize.Width * factor.X,
		Height: spriteSize.Height * factor.Y,
	}
	areaPos := geometry.Point{
		X: spritePos.X - (spriteSize.Width * spriteScale / 2 * gameScale.Max * factor.X),
		Y: spritePos.Y - (spriteSize.Height * spriteScale / 2 * gameScale.Max * factor.Y),
	}
	box := shapes.Get.Box(area)
	box.Size = size
	area.Set(box)
	area.Set(areaPos)
}

// move the gopher using the current move
func moveSystem(_ *goecs.World, delta float32) error {
	pos := geometry.Get.Point(gopher1)
	pos.X += move.X * delta
	pos.Y += move.Y * delta
	gopher1.Set(pos)
	return nil
}

// calculate move on key press
func keyListener(_ *goecs.World, signal interface{}, _ float32) error {
	switch e := signal.(type) {
	case events.KeyDownEvent:
		if e.Key == device.KeyUp {
			move.Y = -gopherSpeed * gameScale.Point.Y
		} else if e.Key == device.KeyDown {
			move.Y = gopherSpeed * gameScale.Point.Y
		} else if e.Key == device.KeyLeft {
			move.X = -gopherSpeed * gameScale.Point.X
		} else if e.Key == device.KeyRight {
			move.X = gopherSpeed * gameScale.Point.X
		}
	case events.KeyUpEvent:
		if e.Key == device.KeyUp {
			move.Y = 0
		} else if e.Key == device.KeyDown {
			move.Y = 0
		} else if e.Key == device.KeyLeft {
			move.X = 0
		} else if e.Key == device.KeyRight {
			move.X = 0
		} else if e.Key == device.KeySpace {
			currentFactor++
			if currentFactor >= len(factors) {
				currentFactor = 0
			}
		}
	}

	return nil
}

// color in red sprites that collides
func collideSystem(_ *goecs.World, _ float32) error {
	// no collision is color white
	color1 := color.White
	color2 := color.White

	// get a components
	pos1 := geometry.Get.Point(gopher1)
	spr1 := sprite.Get(gopher1)

	pos2 := geometry.Get.Point(gopher2)
	spr2 := sprite.Get(gopher2)

	// if they collide
	if gEng.SpritesCollidesFactor(spr1, pos1, spr2, pos2, factors[currentFactor].factor1, factors[currentFactor].factor2) {
		// tim them in red
		color1 = color.Red
		color2 = color.Red
	}

	// update colors
	gopher1.Set(color1)
	gopher2.Set(color2)

	return nil
}
