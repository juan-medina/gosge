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
	buttonExtraWidth  = 0.15                       // the additional width for a button si it is not only the text size
	buttonExtraHeight = 0.15                       // the additional width for a button si it is not only the text size
	shadowExtraWidth  = 5                          // the x offset for the buttons shadow
	shadowExtraHeight = 5                          // the y offset for the buttons shadow
	fontTittle        = 100                        // tittle font size
	fontBig           = 100                        // big buttons font size
	fontSmall         = 50                         // small buttons font size
	fontName          = "resources/go_regular.fnt" // game font
)

func mainStage(eng engine.Engine) error {
	// Preload font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	// pre load sprites
	if err := eng.LoadSpriteSheet("resources/stages.json"); err != nil {
		return err
	}

	// get the world
	wld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// each stage will have it own background color
	eng.SetBackgroundColor(color.SkyBlue)

	// add the centered text
	wld.Add(entity.New(
		ui.Text{
			String:     "Main Stage",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.TopVAlignment,
			Size:       fontTittle * gameScale.Min,
			Font:       fontName,
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

	// add the center sprite
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

	// measure the text so we could make a button just for it
	var measure geometry.Size
	var err error
	if measure, err = eng.MeasureText(fontName, "< back", fontSmall); err != nil {
		return err
	}

	measure.Width += measure.Width * buttonExtraWidth
	measure.Height += measure.Height * buttonExtraHeight

	// add the back button in the top right corner
	wld.Add(entity.New(
		ui.FlatButton{
			Shadow: geometry.Size{Width: shadowExtraWidth, Height: shadowExtraHeight},
			Event:  events.ChangeGameStage{Stage: "menu"},
		},
		geometry.Point{
			X: (designResolution.Width - (measure.Width) - shadowExtraWidth) * gameScale.Point.X,
			Y: (shadowExtraHeight) * gameScale.Point.Y,
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  measure.Width,
				Height: measure.Height,
			},
			Scale: gameScale.Min,
		},
		ui.Text{
			String:     "< back",
			Size:       fontSmall * gameScale.Min,
			Font:       fontName,
			VAlignment: ui.MiddleVAlignment,
			HAlignment: ui.CenterHAlignment,
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
	// Preload font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	// pre load sprites
	if err := eng.LoadSpriteSheet("resources/stages.json"); err != nil {
		return err
	}

	// get the world
	wld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// customize background color
	eng.SetBackgroundColor(color.Gopher)

	// add the centered text
	wld.Add(entity.New(
		ui.Text{
			String:     "Menu",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.TopVAlignment,
			Size:       fontTittle * gameScale.Min,
			Font:       fontName,
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

	// add the center sprite
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

	// measuring the biggest text for size all the buttons equally
	var measure geometry.Size
	var err error
	if measure, err = eng.MeasureText(fontName, "Play !", fontBig); err != nil {
		return err
	}

	measure.Width += measure.Width * buttonExtraWidth
	measure.Height += measure.Height * buttonExtraHeight

	// add the play button, it will sent a event to change to the main stage
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
		ui.Text{
			String:     "Play!",
			Size:       fontBig * gameScale.Min,
			Font:       fontName,
			VAlignment: ui.MiddleVAlignment,
			HAlignment: ui.CenterHAlignment,
		},
		color.Gradient{
			From: color.Blue,
			To:   color.SkyBlue,
		},
		effects.Layer{Depth: 0},
	))

	// add the exit button, it will trigger the event to close the game
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
		ui.Text{
			String:     "Exit",
			Size:       fontBig * gameScale.Min,
			Font:       fontName,
			VAlignment: ui.MiddleVAlignment,
			HAlignment: ui.CenterHAlignment,
		},
		color.Gradient{
			From: color.Blue,
			To:   color.SkyBlue,
		},
		effects.Layer{Depth: 0},
	))

	return nil
}
