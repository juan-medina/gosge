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
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/components/ui"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
)

var opt = options.Options{
	Title:      "Stages Game",
	BackGround: color.Black,
	Icon:       "resources/icon.png",
	Windowed:   true,
}

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
)

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

func loadGame(eng engine.Engine) error {
	eng.AddGameStage("menu", menuStage)
	eng.AddGameStage("main", mainStage)

	return eng.World().Notify(events.ChangeGameStage{Stage: "menu"})
}

// game constants
const (
	buttonExtraWidth  = 0.20
	buttonExtraHeight = 0.20
	shadowExtraWidth  = 5
	shadowExtraHeight = 5
	fontTittle        = 300
	fontBig           = 100
	fontSmall         = 50
	fontSpacing       = 10
)

func mainStage(eng engine.Engine) error {
	// pre load sprites
	if err := eng.LoadSpriteSheet("resources/stages.json"); err != nil {
		return err
	}

	wld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	eng.SetBackgroundColor(color.SkyBlue)

	// add the centered text
	wld.Add(entity.New(
		text.Text{
			String:     "Main Stage",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.TopVAlignment,
			Size:       fontTittle * gameScale.Min,
			Spacing:    fontSpacing * gameScale.Min,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: 0,
		},
		effects.AlternateColor{
			Time:  1,
			Delay: 0.5,
			From:  color.Red,
			To:    color.Yellow,
		},
		effects.Layer{Depth: 0},
	))

	wld.Add(entity.New(
		sprite.Sprite{
			Sheet: "resources/stages.json",
			Name:  "go-fuzz.png",
			Scale: 1 * gameScale.Min,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		effects.Layer{Depth: 1},
	))

	measure := eng.MeasureText("< back", fontSmall, fontSpacing)

	measure.Width += measure.Width * buttonExtraWidth
	measure.Height += measure.Height * buttonExtraHeight

	wld.Add(entity.New(
		ui.FlatButton{
			Shadow: geometry.Size{Width: shadowExtraWidth, Height: shadowExtraHeight},
			Event:  events.ChangeGameStage{Stage: "menu"},
		},
		geometry.Point{
			X: (designResolution.Width - (measure.Width) - shadowExtraWidth) * gameScale.Point.X,
			Y: ((measure.Height / 2) - shadowExtraHeight) * gameScale.Point.Y,
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  measure.Width,
				Height: measure.Height,
			},
			Scale: gameScale.Min,
		},
		text.Text{
			String:     "< back",
			Size:       fontSmall * gameScale.Min,
			Spacing:    fontSpacing * gameScale.Min,
			VAlignment: text.MiddleVAlignment,
			HAlignment: text.CenterHAlignment,
		},
		color.Gradient{
			From: color.Red,
			To:   color.Beige,
		},
		effects.Layer{Depth: 0},
	))

	return nil
}

func menuStage(eng engine.Engine) error {
	log.Printf("load menu stage")
	// pre load sprites
	if err := eng.LoadSpriteSheet("resources/stages.json"); err != nil {
		return err
	}

	wld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	eng.SetBackgroundColor(color.Gopher)

	// add the centered text
	wld.Add(entity.New(
		text.Text{
			String:     "Menu",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.TopVAlignment,
			Size:       fontTittle * gameScale.Min,
			Spacing:    fontSpacing * gameScale.Min,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: 0,
		},
		effects.AlternateColor{
			Time:  1,
			Delay: 0.5,
			From:  color.Red,
			To:    color.Yellow,
		},
		effects.Layer{Depth: 0},
	))

	wld.Add(entity.New(
		sprite.Sprite{
			Sheet: "resources/stages.json",
			Name:  "gamer.png",
			Scale: 1 * gameScale.Min,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		effects.Layer{Depth: 1},
	))

	measure := eng.MeasureText("Play!", fontBig, fontSpacing)
	measure.Width += measure.Width * buttonExtraWidth
	measure.Height += measure.Height * buttonExtraHeight

	wld.Add(entity.New(
		ui.FlatButton{
			Shadow: geometry.Size{Width: shadowExtraWidth, Height: shadowExtraHeight},
			Event:  events.ChangeGameStage{Stage: "main"},
		},
		geometry.Point{
			X: ((designResolution.Width / 2) - (measure.Width / 2)) * gameScale.Point.X,
			Y: ((designResolution.Height - 210) - (measure.Height / 2)) * gameScale.Point.Y,
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  measure.Width,
				Height: measure.Height,
			},
			Scale: gameScale.Min,
		},
		text.Text{
			String:     "Play!",
			Size:       fontBig * gameScale.Min,
			Spacing:    fontSpacing * gameScale.Min,
			VAlignment: text.MiddleVAlignment,
			HAlignment: text.CenterHAlignment,
		},
		color.Gradient{
			From: color.Blue,
			To:   color.SkyBlue,
		},
		effects.Layer{Depth: 0},
	))

	wld.Add(entity.New(
		ui.FlatButton{
			Shadow: geometry.Size{Width: shadowExtraWidth, Height: shadowExtraHeight},
			Event:  events.GameCloseEvent{},
		},
		geometry.Point{
			X: ((designResolution.Width / 2) - (measure.Width / 2)) * gameScale.Point.X,
			Y: ((designResolution.Height - 80) - (measure.Height / 2)) * gameScale.Point.Y,
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  measure.Width,
				Height: measure.Height,
			},
			Scale: gameScale.Min,
		},
		text.Text{
			String:     "Exit",
			Size:       fontBig * gameScale.Min,
			Spacing:    fontSpacing * gameScale.Min,
			VAlignment: text.MiddleVAlignment,
			HAlignment: text.CenterHAlignment,
		},
		color.Gradient{
			From: color.Blue,
			To:   color.SkyBlue,
		},
		effects.Layer{Depth: 0},
	))

	return nil
}
