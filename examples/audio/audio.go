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
	"github.com/juan-medina/gosge/components/animation"
	"github.com/juan-medina/gosge/components/audio"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
	"reflect"
)

var opt = options.Options{
	Title:      "GOSGE Audio Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed: true,
	// Width:    2048,
	// Height:   1536,
}

const (
	fontName                = "resources/go_regular.fnt"   // the font for the game
	fontSmall               = 60                           // the font size
	musicFile               = "resources/audio/loop.ogg"   // music file
	gopherSound             = "resources/audio/gopher.wav" // gopher sound
	clickSound              = "resources/audio/click.wav"  // click sound
	spriteSheet             = "resources/audio.json"       // sprite sheet with our graphics
	playButtonNormalSprite  = "play_button_normal.png"     // the normal sprite for the play button
	playButtonHoverSprite   = "play_button_hover.png"      // the hover sprite for the play button
	stopButtonNormalSprite  = "stop_button_normal.png"     // the normal sprite for the stop button
	stopButtonHoverSprite   = "stop_button_hover.png"      // the hover sprite for the stop button
	pauseButtonNormalSprite = "pause_button_normal.png"    // the normal sprite for the pause button
	pauseButtonHoverSprite  = "pause_button_hover.png"     // the hover sprite for the pause button
	gopherVerticalGap       = 300                          // vertical gap from the center for our gopher
	danceAnim               = "dance"                      // the dance animation name
	idleAnim                = "idle"                       // the idle animation name
	gopherIdle              = "gopher_coffee_%02d.png"     // sprite frame base name for the idle animation
	gopherIdleFrames        = 10                           // number of frames for the gopher idle animation
	gopherDance             = "gopher_dance_%02d.png"      // sprite frame base name for the dance animation
	gopherDanceFrames       = 24                           // number of frames for the gopher dance animation
	gopherFrameDelay        = 0.100                        // delay between frames for the gopher animations
	danceAnimSpeedIncrease  = 0.20                         // how much more faster will dance the gopher on click
	barWith                 = 300                          // our bars width
	barHeight               = 40                           // our bars height
	defaultMusicVolume      = float32(1)                   // the current music volume
	defaultSoundVolume      = float32(1)                   // the current sound volume
	defaultMasterVolume     = float32(1)                   // the current master volume
)

var (
	playButton  *goecs.Entity // the play button entity
	stopButton  *goecs.Entity // the stop button
	gopher      *goecs.Entity // the gopher sprite entity
	geng        *gosge.Engine // the game engine
	masterLabel *goecs.Entity // the master volume label
	masterBar   *goecs.Entity // the master progress bar
	musicLabel  *goecs.Entity // the music volume label
	musicBar    *goecs.Entity // the music progress bar
	soundLabel  *goecs.Entity // the sound volume label
	soundBar    *goecs.Entity // the sound progress bar
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

// BarType is the bar that we click
type BarType int

// BarType for each bar
const (
	MasterBar = BarType(iota) // our master bar
	MusicBar                  // our music bar
	SoundBar                  // our sound bar
)

// BarClickEvent is trigger when the bar is clicked
type BarClickEvent struct {
	Type BarType
}

// BarClickEventType is the reflect.Type for BarClickEvent
var BarClickEventType = reflect.TypeOf(BarClickEvent{})

func loadGame(eng *gosge.Engine) error {
	geng = eng
	var err error

	// set initial values
	settings := eng.GetSettings()
	masterVolume := settings.GetFloat32("masterVolume", defaultMasterVolume)
	musicVolume := settings.GetFloat32("musicVolume", defaultMusicVolume)
	soundVolume := settings.GetFloat32("soundVolume", defaultSoundVolume)

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

	// Preload sounds
	if err = eng.LoadSound(gopherSound); err != nil {
		return err
	}
	if err = eng.LoadSound(clickSound); err != nil {
		return err
	}

	// Get the ECS world
	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     "click on the buttons or the gopher, press <ESC> to close",
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

	// calculate the size of the buttons
	spriteScale := float32(0.25)
	var spriteSize geometry.Size
	if spriteSize, err = eng.GetSpriteSize(spriteSheet, playButtonNormalSprite); err != nil {
		return err
	}
	spriteSize.Width *= spriteScale
	spriteSize.Height *= spriteScale

	buttonsGap := spriteSize.Width / 4

	// add the play button
	playButton = world.AddEntity(
		ui.SpriteButton{
			Sheet:   spriteSheet,
			Normal:  playButtonNormalSprite,
			Hover:   playButtonHoverSprite,
			Clicked: playButtonNormalSprite,
			Scale:   gameScale.Max * spriteScale,
			Sound:   clickSound,
			Volume:  soundVolume,
			Event: events.PlayMusicEvent{ // on click send a event to play the music
				Name:   musicFile,
				Volume: musicVolume,
			},
		},
		geometry.Point{
			X: ((designResolution.Width / 2) - (spriteSize.Width / 2) - buttonsGap) * gameScale.Point.X,
			Y: (designResolution.Height / 2) * gameScale.Point.Y,
		},
		effects.Layer{Depth: 0},
	)

	// add the stop button
	stopButton = world.AddEntity(
		ui.SpriteButton{
			Sheet:   spriteSheet,
			Normal:  stopButtonNormalSprite,
			Hover:   stopButtonHoverSprite,
			Clicked: stopButtonNormalSprite,
			Scale:   gameScale.Max * spriteScale,
			Sound:   clickSound,
			Volume:  soundVolume,
			Event: events.StopMusicEvent{ // on click send a event to stop the music
				Name: musicFile,
			},
		},
		geometry.Point{
			X: ((designResolution.Width / 2) + (spriteSize.Width / 2) + buttonsGap) * gameScale.Point.X,
			Y: (designResolution.Height / 2) * gameScale.Point.Y,
		},
		effects.Layer{Depth: 0},
	)

	// the scale for our gopher sprite
	spriteScale = float32(2)

	// add the gopher with it animations
	gopher = world.AddEntity(
		animation.Animation{
			Sequences: map[string]animation.Sequence{
				danceAnim: {
					Sheet:  spriteSheet,
					Base:   gopherDance,
					Scale:  gameScale.Max * spriteScale,
					Frames: gopherDanceFrames,
					Delay:  gopherFrameDelay,
				},
				idleAnim: {
					Sheet:  spriteSheet,
					Base:   gopherIdle,
					Scale:  gameScale.Max * spriteScale,
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
	)

	textPos := geometry.Point{
		X: (designResolution.Width / 2) * gameScale.Point.X,
		Y: (designResolution.Height / 2) * gameScale.Point.Y,
	}

	var textSize geometry.Size
	if textSize, err = eng.MeasureText(fontName, "Master Volume: 100%", fontSmall); err != nil {
		return err
	}

	// master bar
	textPos.Y += spriteSize.Height * gameScale.Max

	barPos := textPos
	barPos.Y += textSize.Height * gameScale.Max * 0.5
	barPos.X -= barWith * gameScale.Max * 0.5

	masterLabel, masterBar = createBar(world, "Master Volume : %d%%", textPos, barPos, gameScale, MasterBar)

	// music bar
	textPos.Y = barPos.Y + (barHeight * 2 * gameScale.Max)

	barPos = textPos
	barPos.Y += textSize.Height * gameScale.Max * 0.5
	barPos.X -= barWith * gameScale.Max * 0.5

	musicLabel, musicBar = createBar(world, "Music Volume : %d%%", textPos, barPos, gameScale, MusicBar)

	// sound bar
	textPos.Y = barPos.Y + (barHeight * 2 * gameScale.Max)

	barPos = textPos
	barPos.Y += textSize.Height * gameScale.Max * 0.5
	barPos.X -= barWith * gameScale.Max * 0.5

	soundLabel, soundBar = createBar(world, "Sound Volume : %d%%", textPos, barPos, gameScale, SoundBar)

	// add the listener for mouse clicks
	world.AddListener(mouseListener, events.TYPE.MouseUpEvent)

	// add the listener to update the ui when music status change
	world.AddListener(musicStateListener, events.TYPE.MusicStateChangeEvent)

	// listen on bar clicks
	world.AddListener(barClickListener, BarClickEventType)

	textPos.Y += barHeight * gameScale.Max

	// set the master volume
	world.Signal(events.ChangeMasterVolumeEvent{Volume: masterVolume})

	return nil
}

func createBar(world *goecs.World, text string, textPos, barPos geometry.Point, scale geometry.Scale,
	barType BarType) (label, bar *goecs.Entity) {

	settings := geng.GetSettings()
	soundVolume := settings.GetFloat32("soundVolume", defaultSoundVolume)
	masterVolume := settings.GetFloat32("masterVolume", defaultMusicVolume)
	musicVolume := settings.GetFloat32("musicVolume", defaultMusicVolume)

	pro := ui.ProgressBar{
		Min:     0,
		Max:     100,
		Current: 100,
		Shadow: geometry.Size{
			Width:  5 * scale.Max,
			Height: 5 * scale.Max,
		},
		Sound:  clickSound,
		Volume: soundVolume,
		Event: BarClickEvent{
			Type: barType,
		},
	}

	switch barType {
	case MasterBar:
		pro.Current = 100 * masterVolume
		break
	case MusicBar:
		pro.Current = 100 * musicVolume
		break
	case SoundBar:
		pro.Current = 100 * soundVolume
		break
	}

	labelText := fmt.Sprintf(text, int(pro.Current))

	// add the label
	label = world.AddEntity(
		ui.Text{
			String:     labelText,
			Size:       fontSmall,
			Font:       fontName,
			VAlignment: ui.MiddleVAlignment,
			HAlignment: ui.CenterHAlignment,
		},
		textPos,
		color.White,
	)

	// add the bar
	bar = world.AddEntity(
		pro,
		barPos,
		shapes.Box{
			Size: geometry.Size{
				Width:  barWith,
				Height: barHeight,
			},
			Scale:     scale.Max,
			Thickness: int32(2 * scale.Max),
		},
		ui.ProgressBarColor{
			Gradient: color.Gradient{
				From:      color.Blue,
				To:        color.SkyBlue,
				Direction: color.GradientHorizontal,
			},
			Border: color.White,
			Empty:  color.Solid{R: 192, G: 192, B: 192, A: 255},
		},
	)

	return
}

func barClickListener(world *goecs.World, signal interface{}, _ float32) error {
	switch e := signal.(type) {
	case BarClickEvent:
		settings := geng.GetSettings()
		masterVolume := settings.GetFloat32("masterVolume", defaultMusicVolume)
		musicVolume := settings.GetFloat32("musicVolume", defaultMusicVolume)
		soundVolume := settings.GetFloat32("soundVolume", defaultSoundVolume)
		switch e.Type {
		case MasterBar:
			bar := ui.Get.ProgressBar(masterBar)
			text := ui.Get.Text(masterLabel)
			text.String = fmt.Sprintf("Master Volume : %d%%", int(bar.Current))
			masterLabel.Set(text)
			masterVolume = float32(int(bar.Current)) / 100
			settings.SetFloat32("masterVolume", masterVolume)
			updateUIVolume()
			world.Signal(events.ChangeMasterVolumeEvent{Volume: masterVolume})
		case MusicBar:
			bar := ui.Get.ProgressBar(musicBar)
			text := ui.Get.Text(musicLabel)
			text.String = fmt.Sprintf("Music Volume : %d%%", int(bar.Current))
			musicVolume = float32(int(bar.Current)) / 100
			settings.SetFloat32("musicVolume", musicVolume)
			musicLabel.Set(text)
			updateUIVolume()
			world.Signal(events.ChangeMusicVolumeEvent{
				Name:   musicFile,
				Volume: musicVolume,
			})
		case SoundBar:
			bar := ui.Get.ProgressBar(soundBar)
			text := ui.Get.Text(soundLabel)
			text.String = fmt.Sprintf("Sound Volume : %d%%", int(bar.Current))
			soundVolume = float32(int(bar.Current)) / 100
			settings.SetFloat32("soundVolume", soundVolume)
			soundLabel.Set(text)
			updateUIVolume()
		}
	}
	return nil
}

func updateUIVolume() {
	settings := geng.GetSettings()
	musicVolume := settings.GetFloat32("musicVolume", defaultMusicVolume)
	soundVolume := settings.GetFloat32("soundVolume", defaultSoundVolume)

	btn := ui.Get.SpriteButton(playButton)
	btn.Volume = soundVolume
	playButton.Set(btn)

	btn = ui.Get.SpriteButton(stopButton)
	btn.Volume = soundVolume
	stopButton.Set(btn)

	bar := ui.Get.ProgressBar(masterBar)
	bar.Volume = soundVolume
	masterBar.Set(bar)

	bar = ui.Get.ProgressBar(musicBar)
	bar.Volume = soundVolume
	musicBar.Set(bar)

	bar = ui.Get.ProgressBar(soundBar)
	bar.Volume = soundVolume
	soundBar.Set(bar)

	sb := ui.Get.SpriteButton(playButton)
	switch sb.Event.(type) {
	case events.PlayMusicEvent:
		sb.Event = events.PlayMusicEvent{
			Name:   musicFile,
			Volume: musicVolume,
		}
		playButton.Set(sb)
	}
}

func mouseListener(world *goecs.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	// check if we get a mouse up
	case events.MouseUpEvent:
		settings := geng.GetSettings()
		soundVolume := settings.GetFloat32("soundVolume", defaultSoundVolume)
		// get our gopher sprite and position
		pos := geometry.Get.Point(gopher)
		spr := sprite.Get(gopher)
		// if we click on the gopher
		if geng.SpriteAtContains(spr, pos, e.Point) {
			// make the gopher move faster
			anim := animation.Get.Animation(gopher)
			// if we are dancing, dance faster
			if anim.Current == danceAnim {
				anim.Speed += danceAnimSpeedIncrease
			}
			gopher.Set(anim)
			// play the gopher sound
			world.Signal(events.PlaySoundEvent{Name: gopherSound, Volume: soundVolume})
		}
	}
	return nil
}

func musicStateListener(_ *goecs.World, event interface{}, _ float32) error {
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
			sb.Clicked = pauseButtonNormalSprite
			sb.Event = events.PauseMusicEvent{ // on click send an event to pause the music
				Name: e.Name,
			}
			anim.Current = danceAnim
		// if the music has stopped
		case audio.StateStopped:
			settings := geng.GetSettings()
			musicVolume := settings.GetFloat32("musicVolume", defaultMusicVolume)
			// update the play button hover and normal sprites with play sprites
			sb.Normal = playButtonNormalSprite
			sb.Hover = playButtonHoverSprite
			sb.Clicked = playButtonNormalSprite
			sb.Event = events.PlayMusicEvent{ // on click send a event to play the music
				Name:   e.Name,
				Volume: musicVolume,
			}
			anim.Current = idleAnim
			anim.Speed = 1.0
		// if the music pause
		case audio.StatePaused:
			// update the play button hover an normal sprites with play sprites
			sb.Normal = playButtonNormalSprite
			sb.Hover = playButtonHoverSprite
			sb.Clicked = playButtonNormalSprite
			sb.Event = events.ResumeMusicEvent{
				Name: e.Name,
			}
			anim.Current = idleAnim
		}
		// reset anim speed
		anim.Speed = 1
		// update the play button entity
		playButton.Set(sb)
		// update the gopher entity
		gopher.Set(anim)
	}
	return nil
}
