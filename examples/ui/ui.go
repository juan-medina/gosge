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
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
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
	fontName    = "resources/go_regular.fnt"
	fontSmall   = 30
	fontBig     = 50
	columnGap   = 350
	rowGap      = 50
	spriteSheet = "resources/ui.json"
	spriteScale = 0.25
	hint        = "press <SPACE> to hide all, <ESC> to close"
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	message          *goecs.Entity
	gEng             *gosge.Engine
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

	// Preload sprite sheet
	if err := eng.LoadSpriteSheet(spriteSheet); err != nil {
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

	// add sprite button
	pos.Y += rowGap * gameScale.Max
	if err := addSpriteButton(world, gameScale, pos); err != nil {
		return err
	}

	// add the bottom text
	message = world.AddEntity(
		ui.Text{
			String:     hint,
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

	// listen tu ui events
	world.AddListener(uiEvents)

	// listen to keys
	world.AddListener(keyEvents)

	return nil
}

func keyEvents(world *goecs.World, signal interface{}, _ float32) error {
	switch e := signal.(type) {
	case events.KeyUpEvent:
		if e.Key == device.KeySpace {
			hideUnhide(world)
		}
	}
	return nil
}

func hideUnhide(world *goecs.World) {
	for it := world.Iterator(geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.ID() != message.ID() {
			if ent.Contains(effects.TYPE.Hide) {
				ent.Remove(effects.TYPE.Hide)
			} else {
				ent.Add(effects.Hide{})
			}
		}
	}
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

	bar := ui.ProgressBar{
		Min:     0,
		Max:     100,
		Current: 50,
		Shadow: geometry.Size{
			Width:  2 * gameScale.Max,
			Height: 2 * gameScale.Max,
		},
	}

	barEnt := world.AddEntity(
		bar,
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

	controlPos.X += 410 * gameScale.Max

	valueEnt := world.AddEntity(
		ui.Text{
			String:     "50",
			Size:       fontSmall * gameScale.Max,
			Font:       fontName,
			VAlignment: ui.TopVAlignment,
			HAlignment: ui.LeftHAlignment,
		},
		controlPos,
		color.White,
	)

	bar.Event = progressBarEvent{
		barEnt:   barEnt,
		valueEnt: valueEnt,
	}

	barEnt.Set(bar)
}

func addSpriteButton(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point) error {
	text := "SpriteButton"

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
	var err error
	var size geometry.Size
	if size, err = gEng.GetSpriteSize(spriteSheet, "normal.png"); err != nil {
		return err
	}

	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max) + (size.Width * 0.5 * gameScale.Max * spriteScale),
		Y: labelPos.Y + (size.Height * 0.5 * gameScale.Max * spriteScale),
	}

	// add the sprite button
	world.AddEntity(
		ui.SpriteButton{
			Sheet:   spriteSheet,
			Normal:  "normal.png",
			Hover:   "hover.png",
			Clicked: "click.png",
			Scale:   gameScale.Max * spriteScale,
			Event:   uiDemoEvent{Message: text + " clicked"},
		},
		controlPos,
	)
	return nil
}

func uiEvents(_ *goecs.World, signal interface{}, _ float32) error {
	switch e := signal.(type) {
	case uiDemoEvent:
		text := ui.Get.Text(message)
		text.String = e.Message + ", " + hint
		message.Set(text)
	case progressBarEvent:
		bar := ui.Get.ProgressBar(e.barEnt)
		text := ui.Get.Text(e.valueEnt)
		text.String = fmt.Sprintf("%d", int(bar.Current))
		e.valueEnt.Set(text)
	}
	return nil
}

type uiDemoEvent struct {
	Message string
}

type progressBarEvent struct {
	barEnt   *goecs.Entity
	valueEnt *goecs.Entity
}
