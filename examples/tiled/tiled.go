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
	"fmt"
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/tiled"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "Tiled Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	Windowed:   true,
}

const (
	fontName       = "resources/go_regular.fnt"
	fontSmall      = 60
	mapFile        = "resources/maps/gameart2d-desert.tmx"
	mapMoveSpeed   = 850 // map move speed design resolution pixel / second
	defaultTopText = "click a tile with a 'name' property"
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

	// gen is a reference to the gosge.Engine
	gen *gosge.Engine

	// topText is the text on top of the screen
	topText *goecs.Entity
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

	// Get the world
	world := eng.World()

	// save the engine reference
	gen = eng

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

	var textSize geometry.Size
	if textSize, err = eng.MeasureText(fontName, defaultTopText, fontSmall); err != nil {
		return err
	}

	// add the top text
	topText = world.AddEntity(
		ui.Text{
			String:     defaultTopText,
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Min,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: textSize.Height,
		},
		color.Gopher,
		effects.Layer{Depth: 0},
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
	world.AddListener(mouseListener)

	return nil
}

func mouseListener(world *goecs.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	case events.MouseUpEvent:
		for it := world.Iterator(tiled.TYPE.Properties, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
			ent := it.Value()
			properties := tiled.Get.Properties(ent)
			spr := sprite.Get(ent)
			pos := geometry.Get.Point(ent)
			if gen.SpriteAtContains(spr, pos, e.Point) {
				if v, ok := properties.Values["name"]; ok {
					text := ui.Get.Text(topText)
					text.String = fmt.Sprintf("Tiled clicked has 'name' = %q", v)
					topText.Set(text)
					return nil
				}
			}
		}
		text := ui.Get.Text(topText)
		text.String = defaultTopText
		topText.Set(text)
	}
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
