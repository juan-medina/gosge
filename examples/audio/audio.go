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
	"github.com/juan-medina/gosge/pkg/components/audio"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/ui"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"github.com/rs/zerolog/log"
)

var opt = options.Options{
	Title:      "Audio Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
}

const (
	fontName    = "resources/go_regular.fnt"
	fontSmall   = 60
	musicFile   = "resources/audio/loop.ogg"
	spriteSheet = "resources/audio.json"
	buttonsGap  = 5
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
)

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng engine.Engine) error {
	var err error
	// Preload font
	if err = eng.LoadFont(fontName); err != nil {
		return err
	}

	// Preload music
	if err = eng.LoadMusic(musicFile); err != nil {
		return err
	}

	// Preload sprites
	if err = eng.LoadSpriteSheet(spriteSheet); err != nil {
		return err
	}

	wld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// add the bottom text
	wld.Add(entity.New(
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
	))

	spriteScale := float32(0.25)
	var spriteSize geometry.Size
	if spriteSize, err = eng.GetSpriteSize(spriteSheet, "play_button_normal.png"); err != nil {
		return err
	}
	spriteSize.Width *= spriteScale
	spriteSize.Height *= spriteScale

	wld.Add(entity.New(
		ui.SpriteButton{
			Sheet:  spriteSheet,
			Normal: "play_button_normal.png",
			Hover:  "play_button_hover.png",
			Scale:  gameScale.Min * spriteScale,
			Event: events.PlayMusicEvent{
				Name:  musicFile,
				Loops: audio.LoopForever,
			},
		},
		geometry.Point{
			X: ((designResolution.Width / 2) - (spriteSize.Width / 2) - buttonsGap) * gameScale.Point.X,
			Y: (designResolution.Height / 2) * gameScale.Point.Y,
		},
		effects.Layer{Depth: 0},
	))

	wld.Add(entity.New(
		ui.SpriteButton{
			Sheet:  spriteSheet,
			Normal: "stop_button_normal.png",
			Hover:  "stop_button_hover.png",
			Scale:  gameScale.Min * spriteScale,
			Event: events.StopMusicEvent{
				Name: musicFile,
			},
		},
		geometry.Point{
			X: ((designResolution.Width / 2) + (spriteSize.Width / 2) + buttonsGap) * gameScale.Point.X,
			Y: (designResolution.Height / 2) * gameScale.Point.Y,
		},
		effects.Layer{Depth: 0},
	))

	return err
}
