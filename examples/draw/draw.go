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
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "GOSGE Draw Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed:   true,
	// Width:      2048,
	// Height:     1536,
}

const (
	fontName  = "resources/go_regular.fnt"
	fontSmall = 60
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}

	// currentLine is the line that we are currently drawing
	currentLine goecs.EntityID

	// pos is the current mouse position
	pos geometry.Point

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale geometry.Scale
)

func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng *gosge.Engine) error {

	// Preload font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale = eng.GetScreenSize().CalculateScale(designResolution)

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     "press mouse and move to drawn lines, press <ESC> to close",
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

	// add the mouse listener
	world.AddListener(mouseListener, events.TYPE.MouseDownEvent, events.TYPE.MouseUpEvent, events.TYPE.MouseMoveEvent)

	return nil
}

func mouseListener(world *goecs.World, signal goecs.Component, _ float32) error {
	switch e := signal.(type) {
	// when the mouse move
	case events.MouseMoveEvent:
		pos = e.Point // store the mouse pos
		// if we have a line
		if currentLine != 0 {
			var err error
			var lineEnt *goecs.Entity
			if lineEnt, err = world.Get(currentLine); err != nil {
				return err
			}
			// get the line component
			line := shapes.Get.Line(lineEnt)
			// update the to
			line.To = pos
			lineEnt.Set(line)
		}
	// when we press the mouse
	case events.MouseDownEvent:
		// if we don't have a line
		if currentLine == 0 {
			// create one and save the reference
			currentLine = world.AddEntity(
				shapes.Line{
					To:        pos,
					Thickness: 5 * gameScale.Max,
				},
				pos,
			)
		}
	// when we release the mouse
	case events.MouseUpEvent:
		// remove reference
		currentLine = 0
	}
	return nil
}
