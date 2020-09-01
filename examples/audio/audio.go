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
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/animation"
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
	fontName                = "resources/go_regular.fnt"
	fontSmall               = 60
	musicFile               = "resources/audio/loop.ogg"
	spriteSheet             = "resources/audio.json"
	buttonsGap              = 5
	playButtonNormalSprite  = "play_button_normal.png"
	playButtonHoverSprite   = "play_button_hover.png"
	stopButtonNormalSprite  = "stop_button_normal.png"
	stopButtonHoverSprite   = "stop_button_hover.png"
	pauseButtonNormalSprite = "pause_button_normal.png"
	pauseButtonHoverSprite  = "pause_button_hover.png"
	danceAnim               = "dance"
	idleAnim                = "idle"
	gopherIdle              = "gopher_coffee_%02d.png"
	gopherDance             = "gopher_dance_%02d.png"
)

var (
	playButton *entity.Entity
	gopher     *entity.Entity
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
	if spriteSize, err = eng.GetSpriteSize(spriteSheet, playButtonNormalSprite); err != nil {
		return err
	}
	spriteSize.Width *= spriteScale
	spriteSize.Height *= spriteScale

	playButton = wld.Add(entity.New(
		ui.SpriteButton{
			Sheet:  spriteSheet,
			Normal: playButtonNormalSprite,
			Hover:  playButtonHoverSprite,
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
			Normal: stopButtonNormalSprite,
			Hover:  stopButtonHoverSprite,
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

	spriteScale = float32(2)

	// add the gopher with it animations
	gopher = wld.Add(entity.New(
		animation.Animation{
			Sequences: map[string]animation.Sequence{
				danceAnim: {
					Sheet:  spriteSheet,
					Base:   gopherDance,
					Scale:  gameScale.Min * spriteScale,
					Frames: 24,
					Delay:  0.100,
				},
				idleAnim: {
					Sheet:  spriteSheet,
					Base:   gopherIdle,
					Scale:  gameScale.Min * spriteScale,
					Frames: 10,
					Delay:  0.100,
				},
			},
			Current: idleAnim, // current animation is idle
			Speed:   1,
		},
		geometry.Point{
			X: (designResolution.Width / 2) * gameScale.Point.X,
			Y: (designResolution.Height/2)*gameScale.Point.Y - 300,
		},
	))

	wld.AddSystem(UISystem())
	return err
}

type uiSystem struct{}

func (us uiSystem) Notify(_ *world.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	case events.MusicStateChangeEvent:
		sb := ui.Get.SpriteButton(playButton)
		anim := animation.Get.Animation(gopher)

		switch e.New {
		case audio.StatePlaying:
			sb.Normal = pauseButtonNormalSprite
			sb.Hover = pauseButtonHoverSprite
			sb.Event = events.PauseMusicEvent{
				Name: e.Name,
			}
			anim.Current = danceAnim
		case audio.StateStopped:
			sb.Normal = playButtonNormalSprite
			sb.Hover = playButtonHoverSprite
			sb.Event = events.PlayMusicEvent{
				Name:  e.Name,
				Loops: audio.LoopForever,
			}
			anim.Current = idleAnim
		case audio.StatePaused:
			sb.Normal = playButtonNormalSprite
			sb.Hover = playButtonHoverSprite
			sb.Event = events.ResumeMusicEvent{
				Name: e.Name,
			}
			anim.Current = idleAnim
		}
		playButton.Set(sb)
		gopher.Set(anim)
	}
	return nil
}

func (us uiSystem) Update(_ *world.World, _ float32) error {
	return nil
}

// UISystem create our example UI world.System
func UISystem() world.System {
	return &uiSystem{}
}
