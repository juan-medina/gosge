/*
 * Copyright (c) 2020 Juan Medina.
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
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
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "Collision Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
}

const (
	fontName    = "resources/go_regular.fnt"
	fontSmall   = 60
	gopherSpeed = 500
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	gEng             *gosge.Engine
	gopher           *goecs.Entity
	move             geometry.Point
	gameScale        geometry.Scale
)

func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng *gosge.Engine) error {
	gEng = eng

	// Preload font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale = eng.GetScreenSize().CalculateScale(designResolution)

	// preload sprite sheet
	if err := eng.LoadSpriteSheet("resources/gamer.json"); err != nil {
		return err
	}

	// add an static gopher
	world.AddEntity(
		sprite.Sprite{
			Sheet: "resources/gamer.json",
			Name:  "gamer.png",
			Scale: gameScale.Min * 0.25,
			FlipX: false,
			FlipY: false},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
	)

	// add our moving gopher
	gopher = world.AddEntity(
		sprite.Sprite{
			Sheet: "resources/gamer.json",
			Name:  "gamer.png",
			Scale: gameScale.Min * 0.25,
			FlipX: false,
			FlipY: false},
		geometry.Point{
			X: designResolution.Width / 4 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
	)

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     "press <ESC> to close",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Min,
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

	// add the key listener
	world.AddListener(keyListener)

	return nil
}

// move the gopher using the current move
func moveSystem(_ *goecs.World, delta float32) error {
	pos := geometry.Get.Point(gopher)
	pos.X += move.X * delta
	pos.Y += move.Y * delta
	gopher.Set(pos)
	return nil
}

// calculate move on key press
func keyListener(_ *goecs.World, signal interface{}, _ float32) error {
	switch e := signal.(type) {
	case events.KeyDownEvent:
		if e.Key == device.KeyUp {
			move.Y = -gopherSpeed * gameScale.Point.Y
		}
		if e.Key == device.KeyDown {
			move.Y = gopherSpeed * gameScale.Point.Y
		}
		if e.Key == device.KeyLeft {
			move.X = -gopherSpeed * gameScale.Point.X
		}
		if e.Key == device.KeyRight {
			move.X = gopherSpeed * gameScale.Point.X
		}
	case events.KeyUpEvent:
		if e.Key == device.KeyUp {
			move.Y = 0
		}
		if e.Key == device.KeyDown {
			move.Y = 0
		}
		if e.Key == device.KeyLeft {
			move.X = 0
		}
		if e.Key == device.KeyRight {
			move.X = 0
		}
	}

	return nil
}

// color in red sprites that collides
func collideSystem(world *goecs.World, _ float32) error {
	// go trough all the sprites
	for itA := world.Iterator(sprite.TYPE, geometry.TYPE.Point); itA != nil; itA = itA.Next() {
		// no collision is color white
		colorA := color.White
		// get a components
		entA := itA.Value()
		posA := geometry.Get.Point(entA)
		sprA := sprite.Get(entA)
		// go trough all the sprites
		for itB := world.Iterator(sprite.TYPE, geometry.TYPE.Point); itB != nil; itB = itB.Next() {
			// no collision is color white
			colorB := color.White
			entB := itB.Value()
			posB := geometry.Get.Point(entB)
			sprB := sprite.Get(entB)
			// if a and b are not the same sprite and collide
			if entA.ID() != entB.ID() && gEng.SpritesCollides(sprA, posA, sprB, posB) {
				// tim them in red
				colorA = color.Red
				colorB = color.Red
			}
			// update b color
			entB.Set(colorB)
		}
		// update a color
		entA.Set(colorA)
	}

	return nil
}
