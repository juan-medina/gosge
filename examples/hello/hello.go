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
	"github.com/juan-medina/gosge"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "Hello Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
}

const (
	fontName  = "resources/go_regular.fnt"
	fontBig   = 100
	fontSmall = 60
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
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
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// add the centered text
	world.AddEntity(
		ui.Text{
			String:     "Hello World",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.MiddleVAlignment,
			Font:       fontName,
			Size:       fontBig * gameScale.Min,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		effects.AlternateColor{
			Time:  1,
			Delay: 0.5,
			From:  color.Red,
			To:    color.Yellow,
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
	return nil
}
