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
	Title:      "GOSGE Tiled Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed:   true,
	// Width:      2048,
	// Height:     1536,
}

const (
	fontName       = "resources/go_regular.fnt"              // our game font
	fontSmall      = 60                                      // font size
	mapFile        = "resources/maps/gameart2d-desert.tmx"   // tiled map file
	mapMoveSpeed   = 850                                     // map move speed, pixel / second
	defaultTopText = "click a tile with a \"name\" property" // default top text
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

	// move indicate how we are moving
	move geometry.Point

	// our game scale
	gameScale geometry.Scale
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
	gameScale = eng.GetScreenSize().CalculateScale(designResolution)

	// add a gradient background
	world.AddEntity(
		shapes.SolidBox{
			Size: geometry.Size{
				Width:  designResolution.Width,
				Height: designResolution.Height,
			},
			Scale: gameScale.Max,
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

	// calculate minimum map position for scrolling
	minMapPos = geometry.Point{
		X: 0,
		Y: (designResolution.Height - mapSize.Height) * gameScale.Point.Y,
	}

	// calculate maximum map position for scrolling
	maxMapPos = geometry.Point{
		X: (mapSize.Width - designResolution.Width) * gameScale.Point.X,
		Y: 0,
	}

	// add the map
	mapEnt = world.AddEntity(
		tiled.Map{
			Name:  mapFile,
			Scale: gameScale.Max,
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
			Size:       fontSmall * gameScale.Max,
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
		effects.Layer{Depth: 0},
	)

	// add our key listener
	world.AddListener(keyListener)
	// add our mouse listener
	world.AddListener(mouseListener)
	// add our move system
	world.AddSystem(moveSystem)

	return nil
}

// listen to mouse clicks
func mouseListener(world *goecs.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	case events.MouseUpEvent:
		// the block info of a clicked block
		var clickedInfo tiled.BlockInfo
		// get ll the blocks
		for it := world.Iterator(tiled.TYPE.BlockInfo, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
			// get values
			ent := it.Value()
			info := tiled.Get.BlockInfo(ent)
			spr := sprite.Get(ent)
			pos := geometry.Get.Point(ent)
			// if we click it with the mouse
			if gen.SpriteAtContains(spr, pos, e.Point) {
				// if we have name property
				if _, ok := info.Properties["name"]; ok {
					// we click this block
					clickedInfo = info
				}
			}
		}
		// get the top text component
		text := ui.Get.Text(topText)
		// if we don't have a layer we haven't click nothing
		if clickedInfo.Layer == "" {
			// set default text
			text.String = defaultTopText
		} else {
			// set the text with the block information
			text.String = fmt.Sprintf("Tiled clicked has name = %q, layer = %q", clickedInfo.Properties["name"], clickedInfo.Layer)
		}
		// update the top entity
		topText.Set(text)
	}
	return nil
}

// listen to keys
func keyListener(_ *goecs.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	case events.KeyDownEvent:
		switch e.Key {
		case device.KeyLeft:
			move.X = -mapMoveSpeed
		case device.KeyRight:
			move.X = mapMoveSpeed
		case device.KeyUp:
			move.Y = -mapMoveSpeed
		case device.KeyDown:
			move.Y = mapMoveSpeed
		}
	case events.KeyUpEvent:
		switch e.Key {
		case device.KeyLeft:
			move.X = 0
		case device.KeyRight:
			move.X = 0
		case device.KeyUp:
			move.Y = 0
		case device.KeyDown:
			move.Y = 0
		}
	}
	return nil
}

func moveSystem(_ *goecs.World, delta float32) error {
	// if we need to move anything
	if move.X != 0 || move.Y != 0 {
		// get the current position
		pos := geometry.Get.Point(mapEnt)

		// move it
		pos.X += move.X * delta * gameScale.Point.X
		pos.Y -= move.Y * delta * gameScale.Point.Y

		// clamp to min and max scroll pos
		pos.Clamp(minMapPos, maxMapPos)

		// update entity
		mapEnt.Set(pos)
	}
	return nil
}
