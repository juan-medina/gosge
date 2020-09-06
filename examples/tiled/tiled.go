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
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/tiled"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "Hello Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
}

const (
	fontName     = "resources/go_regular.fnt"
	fontSmall    = 60
	mapFile      = "resources/maps/gameart2d-desert.tmx"
	mapMoveSpeed = 850 // map move speed design resolution pixel / second
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	// mapEnt is our map entity
	mapEnt *goecs.Entity

	// minMapPos is the min position that we could scroll the map
	minMapPos geometry.Point

	// maxMapPos is the max position that we could scroll the map
	maxMapPos geometry.Point
)

func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng *gosge.Engine) error {
	var err error
	// Preload font
	if err = eng.LoadFont(fontName); err != nil {
		return err
	}

	// Preload map
	if err = eng.LoadTiledMap(mapFile); err != nil {
		return err
	}

	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	world.AddEntity(
		shapes.Box{
			Size: geometry.Size{
				Width:  designResolution.Width,
				Height: designResolution.Height,
			},
			Scale: gameScale.Min,
		},
		geometry.Point{},
		color.Gradient{
			From:      color.Orange.Blend(color.SkyBlue, 0.50),
			To:        color.Orange,
			Direction: color.GradientVertical,
		},
	)

	// Get map size
	var mapSize geometry.Size
	if mapSize, err = eng.GeTiledMapSize(mapFile); err != nil {
		return err
	}

	minMapPos = geometry.Point{
		X: 0,
		Y: (designResolution.Height - mapSize.Height) * gameScale.Point.Y,
	}

	maxMapPos = geometry.Point{
		X: (mapSize.Width - designResolution.Width) * gameScale.Point.X,
		Y: 0,
	}

	// add the map
	mapEnt = world.AddEntity(
		tiled.Map{
			Name:  mapFile,
			Scale: gameScale.Min,
		},
		minMapPos,
	)

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     "move the map with cursors, press <ESC> to close",
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
		effects.Layer{Depth: 0},
	)

	world.AddListener(keyListener)

	return nil
}

func keyListener(_ *goecs.World, event interface{}, delta float32) error {
	switch e := event.(type) {
	case events.KeyEvent:
		if e.Status.Down {
			move := geometry.Point{}
			switch e.Key {
			case device.KeyLeft:
				move.X = -mapMoveSpeed * delta
			case device.KeyRight:
				move.X = mapMoveSpeed * delta
			case device.KeyUp:
				move.Y = -mapMoveSpeed * delta
			case device.KeyDown:
				move.Y = mapMoveSpeed * delta
			}

			if move.X != 0 || move.Y != 0 {
				pos := geometry.Get.Point(mapEnt)

				pos.X += move.X
				pos.Y -= move.Y
				pos.Clamp(minMapPos, maxMapPos)

				mapEnt.Set(pos)
			}
		}
	}
	return nil
}
