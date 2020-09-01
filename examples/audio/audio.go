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
	fontName                = "resources/go_regular.fnt" // the font for the game
	fontSmall               = 60                         // the font size
	musicFile               = "resources/audio/loop.ogg" // music file
	spriteSheet             = "resources/audio.json"     // sprite sheet with our graphics
	buttonsGap              = 5                          // separation between buttons
	playButtonNormalSprite  = "play_button_normal.png"   // the normal sprite for the play button
	playButtonHoverSprite   = "play_button_hover.png"    // the hover sprite for the play button
	stopButtonNormalSprite  = "stop_button_normal.png"   // the normal sprite for the stop button
	stopButtonHoverSprite   = "stop_button_hover.png"    // the hover sprite for the stop button
	pauseButtonNormalSprite = "pause_button_normal.png"  // the normal sprite for the pause button
	pauseButtonHoverSprite  = "pause_button_hover.png"   // the hover sprite for the pause button
	gopherVerticalGap       = 300                        // vertical gap from the center for our gopher
	danceAnim               = "dance"                    // the dance animation name
	idleAnim                = "idle"                     // the idle animation name
	gopherIdle              = "gopher_coffee_%02d.png"   // sprite frame base name for the idle animation
	gopherIdleFrames        = 10                         // number of frames for the gopher idle animation
	gopherDance             = "gopher_dance_%02d.png"    // sprite frame base name for the dance animation
	gopherDanceFrames       = 24                         // number of frames for the gopher dance animation
	gopherFrameDelay        = 0.100                      // delay between frames for the gopher animations
)

var (
	playButton *entity.Entity // the play button entity
	gopher     *entity.Entity // the gopher sprite entity
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

	// Get the ECS world
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

	// calculate the size of the buttons
	spriteScale := float32(0.25)
	var spriteSize geometry.Size
	if spriteSize, err = eng.GetSpriteSize(spriteSheet, playButtonNormalSprite); err != nil {
		return err
	}
	spriteSize.Width *= spriteScale
	spriteSize.Height *= spriteScale

	// add the play button
	playButton = wld.Add(entity.New(
		ui.SpriteButton{
			Sheet:  spriteSheet,
			Normal: playButtonNormalSprite,
			Hover:  playButtonHoverSprite,
			Scale:  gameScale.Min * spriteScale,
			Event: events.PlayMusicEvent{ // on click send a event to play the music
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

	// add the stop button
	wld.Add(entity.New(
		ui.SpriteButton{
			Sheet:  spriteSheet,
			Normal: stopButtonNormalSprite,
			Hover:  stopButtonHoverSprite,
			Scale:  gameScale.Min * spriteScale,
			Event: events.StopMusicEvent{ // on click send a event to stop the music
				Name: musicFile,
			},
		},
		geometry.Point{
			X: ((designResolution.Width / 2) + (spriteSize.Width / 2) + buttonsGap) * gameScale.Point.X,
			Y: (designResolution.Height / 2) * gameScale.Point.Y,
		},
		effects.Layer{Depth: 0},
	))

	// the scale for our gopher sprite
	spriteScale = float32(2)

	// add the gopher with it animations
	gopher = wld.Add(entity.New(
		animation.Animation{
			Sequences: map[string]animation.Sequence{
				danceAnim: {
					Sheet:  spriteSheet,
					Base:   gopherDance,
					Scale:  gameScale.Min * spriteScale,
					Frames: gopherDanceFrames,
					Delay:  gopherFrameDelay,
				},
				idleAnim: {
					Sheet:  spriteSheet,
					Base:   gopherIdle,
					Scale:  gameScale.Min * spriteScale,
					Frames: gopherIdleFrames,
					Delay:  gopherFrameDelay,
				},
			},
			Current: idleAnim, // current animation is idle
			Speed:   1,
		},
		geometry.Point{
			X: (designResolution.Width / 2) * gameScale.Point.X,
			Y: (designResolution.Height/2)*gameScale.Point.Y - gopherVerticalGap,
		},
	))

	// add the system that will handle the UI
	wld.AddSystem(UISystem())
	return err
}

// System for our game ui
type uiSystem struct{}

// notify a new event
func (us uiSystem) Notify(_ *world.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	// if is the music state has change
	case events.MusicStateChangeEvent:
		// get the play button ui.SpriteButton
		sb := ui.Get.SpriteButton(playButton)
		// get the gopher current animation.Animation
		anim := animation.Get.Animation(gopher)
		switch e.New {
		// if the music is playing
		case audio.StatePlaying:
			// update the play button hover and normal sprites pause sprites
			sb.Normal = pauseButtonNormalSprite
			sb.Hover = pauseButtonHoverSprite
			sb.Event = events.PauseMusicEvent{ // on click send an event to pause the music
				Name: e.Name,
			}
			anim.Current = danceAnim
		// if the music has stopped
		case audio.StateStopped:
			// update the play button hover and normal sprites with play sprites
			sb.Normal = playButtonNormalSprite
			sb.Hover = playButtonHoverSprite
			sb.Event = events.PlayMusicEvent{ // on click send a event to play the music
				Name:  e.Name,
				Loops: audio.LoopForever,
			}
			anim.Current = idleAnim
		// if the music pause
		case audio.StatePaused:
			// update the play button hover an normal sprites with play sprites
			sb.Normal = playButtonNormalSprite
			sb.Hover = playButtonHoverSprite
			sb.Event = events.ResumeMusicEvent{
				Name: e.Name,
			}
			anim.Current = idleAnim
		}
		// update the play button entity
		playButton.Set(sb)
		// update the gopher entity
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
