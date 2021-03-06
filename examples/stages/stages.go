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
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "GOSGE Stages Game",
	BackGround: color.Black,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed:   true,
	// Width:      2048,
	// Height:     1536,
}

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
	eng.AddGameStage("menu", menuStage)
	eng.AddGameStage("main", mainStage)

	eng.World().Signal(events.ChangeGameStage{Stage: "menu"})
	return nil
}

// game constants
const (
	buttonExtraWidth       = 0.15                        // the additional width for a button si it is not only the text size
	buttonExtraHeight      = 0.15                        // the additional width for a button si it is not only the text size
	shadowExtraWidth       = 5                           // the x offset for the buttons shadow
	shadowExtraHeight      = 5                           // the y offset for the buttons shadow
	fontTittle             = 100                         // tittle font size
	fontBig                = 100                         // big buttons font size
	fontName               = "resources/go_regular.fnt"  // game font
	spriteSheetName        = "resources/stages.json"     // game sprite sheet
	gameSprite             = "go-fuzz.png"               // game sprite
	menuSprite             = "gamer.png"                 // menu sprite
	buttonExitNormalSprite = "button_exit_normal.png"    // exit button sprite normal state
	buttonExitHoverSprite  = "button_exit_hover.png"     // exit button sprite hover state
	clickSound             = "resources/audio/click.wav" // click sound
)

func mainStage(eng *gosge.Engine) error {
	var err error

	eng.DisableExitKey()

	// Preload font
	if err = eng.LoadFont(fontName); err != nil {
		return err
	}

	// pre load sprites
	if err = eng.LoadSpriteSheet(spriteSheetName); err != nil {
		return err
	}

	// pre load sounds
	if err = eng.LoadSound(clickSound); err != nil {
		return err
	}

	// get the world
	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// each stage will have it own background color
	eng.SetBackgroundColor(color.SkyBlue)

	// add the centered text
	world.AddEntity(
		ui.Text{
			String:     "Main Stage",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.TopVAlignment,
			Size:       fontTittle * gameScale.Max,
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
	)

	// add the center sprite
	world.AddEntity(
		sprite.Sprite{
			Sheet: spriteSheetName,
			Name:  gameSprite,
			Scale: 1 * gameScale.Max,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		effects.Layer{Depth: 1},
	)

	spriteScale := float32(0.5)
	var spriteSize geometry.Size
	if spriteSize, err = eng.GetSpriteSize(spriteSheetName, buttonExitNormalSprite); err != nil {
		return err
	}
	spriteSize.Width *= spriteScale
	spriteSize.Height *= spriteScale

	world.AddEntity(
		ui.SpriteButton{
			Sheet:   spriteSheetName,
			Normal:  buttonExitNormalSprite,
			Hover:   buttonExitHoverSprite,
			Clicked: buttonExitNormalSprite,
			Scale:   gameScale.Max * spriteScale,
			Sound:   clickSound,
			Volume:  1,
			Event: events.DelaySignal{
				Signal: events.ChangeGameStage{Stage: "menu"},
				Time:   0.15,
			},
		},
		geometry.Point{
			X: (designResolution.Width - (spriteSize.Width / 2)) * gameScale.Point.X,
			Y: (spriteSize.Height / 2) * gameScale.Point.Y,
		},
		effects.Layer{Depth: 0},
	)

	world.AddListener(mainStageKeyListener, events.TYPE.KeyUpEvent)

	return nil
}

func mainStageKeyListener(world *goecs.World, signal goecs.Component, _ float32) error {
	switch e := signal.(type) {
	case events.KeyUpEvent:
		if e.Key == device.KeyEscape {
			world.Signal(events.ChangeGameStage{Stage: "menu"})
		}
	}
	return nil
}

func menuStage(eng *gosge.Engine) error {
	var err error

	eng.SetExitKey(device.KeyEscape)

	// Preload font
	if err = eng.LoadFont(fontName); err != nil {
		return err
	}

	// pre load sprites
	if err = eng.LoadSpriteSheet(spriteSheetName); err != nil {
		return err
	}

	// pre load sounds
	if err = eng.LoadSound(clickSound); err != nil {
		return err
	}

	// get the world
	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// customize background color
	eng.SetBackgroundColor(color.Gopher)

	// add the centered text
	world.AddEntity(
		ui.Text{
			String:     "Menu",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.TopVAlignment,
			Size:       fontTittle * gameScale.Max,
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
	)

	// add the center sprite
	world.AddEntity(
		sprite.Sprite{
			Sheet: spriteSheetName,
			Name:  menuSprite,
			Scale: 1 * gameScale.Max,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		effects.Layer{Depth: 1},
	)

	// measuring the biggest text for size all the buttons equally
	var measure geometry.Size
	if measure, err = eng.MeasureText(fontName, "Play !", fontBig); err != nil {
		return err
	}

	measure.Width += measure.Width * buttonExtraWidth
	measure.Height += measure.Height * buttonExtraHeight

	// add the play button, it will sent a event to change to the main stage
	world.AddEntity(
		ui.FlatButton{
			Shadow: geometry.Size{Width: shadowExtraWidth, Height: shadowExtraHeight},
			Sound:  clickSound,
			Volume: 1,
			Event: events.DelaySignal{
				Signal: events.ChangeGameStage{Stage: "main"},
				Time:   0.15,
			},
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
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		ui.Text{
			String:     "Play!",
			Size:       fontBig * gameScale.Max,
			Font:       fontName,
			VAlignment: ui.MiddleVAlignment,
			HAlignment: ui.CenterHAlignment,
		},
		ui.ButtonColor{
			Gradient: color.Gradient{
				From: color.Blue,
				To:   color.SkyBlue,
			},
			Border: color.White,
			Text:   color.White,
		},
		effects.Layer{Depth: 0},
	)

	// add the exit button, it will trigger the event to close the game
	world.AddEntity(
		ui.FlatButton{
			Shadow: geometry.Size{Width: shadowExtraWidth, Height: shadowExtraHeight},
			Sound:  clickSound,
			Volume: 1,
			Event: events.DelaySignal{
				Signal: events.GameCloseEvent{},
				Time:   0.15,
			},
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
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		ui.Text{
			String:     "Exit",
			Size:       fontBig * gameScale.Max,
			Font:       fontName,
			VAlignment: ui.MiddleVAlignment,
			HAlignment: ui.CenterHAlignment,
		},
		ui.ButtonColor{
			Gradient: color.Gradient{
				From: color.Blue,
				To:   color.SkyBlue,
			},
			Border: color.White,
			Text:   color.White,
		},
		effects.Layer{Depth: 0},
	)

	world.AddListener(menuStageKeyListener, events.TYPE.KeyUpEvent)

	return nil
}

func menuStageKeyListener(world *goecs.World, signal goecs.Component, _ float32) error {
	switch e := signal.(type) {
	case events.KeyUpEvent:
		if e.Key == device.KeyReturn {
			world.Signal(events.ChangeGameStage{Stage: "main"})
		}
	}
	return nil
}
