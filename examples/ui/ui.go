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
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "GOSGE UI Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed:   true,
	// Width:      2048,
	// Height:     1536,
}

const (
	fontName  = "resources/go_regular.fnt"
	fontSmall = 30
	fontBig   = 60
	columnGap = 350
	rowGap    = 50
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	message          *goecs.Entity
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

	// element pos
	pos := geometry.Point{
		X: 10 * gameScale.Max,
		Y: 10 * gameScale.Max,
	}

	// add the flat button
	addFlatButton(world, gameScale, pos, false)
	pos.Y += rowGap * gameScale.Max
	addFlatButton(world, gameScale, pos, true)

	// add the progress bar
	pos.Y += rowGap * gameScale.Max
	addProgressBar(world, gameScale, pos, false)
	pos.Y += rowGap * gameScale.Max
	addProgressBar(world, gameScale, pos, true)

	// add the bottom text
	message = world.AddEntity(
		ui.Text{
			String:     "press <ESC> to close",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       fontBig * gameScale.Max,
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

	world.AddListener(uiEvents)
	return nil
}

func addFlatButton(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point, gradient bool) {
	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max),
		Y: labelPos.Y,
	}

	text := "FlatButton [Solid]"

	clr := ui.ButtonColor{
		Solid:  color.SkyBlue,
		Border: color.White,
		Text:   color.White,
	}

	if gradient {
		clr.Gradient = color.Gradient{
			From:      color.SkyBlue,
			To:        color.DarkBlue,
			Direction: color.GradientHorizontal,
		}
		text = "FlatButton [Gradient]"
	}

	// add a label
	world.AddEntity(
		ui.Text{
			String:     text,
			HAlignment: ui.LeftHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		labelPos,
		color.White,
	)

	// add a control : flat button
	world.AddEntity(
		ui.FlatButton{
			Shadow: geometry.Size{
				Width:  2 * gameScale.Max,
				Height: 2 * gameScale.Max,
			},
			Event: uiDemoEvent{Message: text + " clicked"},
		},
		clr,
		ui.Text{
			String:     "Click Me",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.MiddleVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  60 * gameScale.Max,
				Height: 20 * gameScale.Max,
			},
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		controlPos,
	)
}

func addProgressBar(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point, gradient bool) {
	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max),
		Y: labelPos.Y,
	}

	text := "ProgressBar [Solid]"

	clr := ui.ProgressBarColor{
		Solid:  color.SkyBlue,
		Border: color.White,
	}

	if gradient {
		clr.Gradient = color.Gradient{
			From:      color.SkyBlue,
			To:        color.DarkBlue,
			Direction: color.GradientHorizontal,
		}
		text = "ProgressBar [Gradient]"
	}

	// add a label
	world.AddEntity(
		ui.Text{
			String:     text,
			HAlignment: ui.LeftHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		labelPos,
		color.White,
	)

	world.AddEntity(
		ui.ProgressBar{
			Min:     0,
			Max:     100,
			Current: 50,
			Shadow: geometry.Size{
				Width:  2 * gameScale.Max,
				Height: 2 * gameScale.Max,
			},
			Event: uiDemoEvent{Message: text + " clicked"},
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  200 * gameScale.Max,
				Height: 20 * gameScale.Max,
			},
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		clr,
		controlPos,
	)
}

func uiEvents(_ *goecs.World, signal interface{}, _ float32) error {
	switch e := signal.(type) {
	case uiDemoEvent:
		text := ui.Get.Text(message)
		text.String = e.Message + ", press <ESC> to close"
		message.Set(text)
	}
	return nil
}

type uiDemoEvent struct {
	Message string
}
