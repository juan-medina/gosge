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
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
)

// Game options
var opt = options.Options{
	Title:      "GOSGE Animation Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed:   true,
	// Width:      2048,
	// Height:     1536,
}

// Game constants
const (
	fontName   = "resources/go_regular.fnt"               // font name for the bottom text
	fontSize   = 60                                       // bottom text size
	fontLayer  = -10                                      // the layer for the top will be on front
	robotSheet = "resources/robot.json"                   // sprite sheet with the robot animation
	robotIdle  = "Idle_%02d.png"                          // robot idle animations frames name
	robotRun   = "Run_%02d.png"                           // robot running animation frames name name
	robotPosY  = 220                                      // Y position for our robot
	robotLayer = 0                                        // layer for our robot, above the city bellow the text
	runAnim    = "run"                                    // name for the running animation
	idleAnim   = "idle"                                   // name for the idle animation
	numLayers  = 8                                        // number of layers for the city
	layerFile  = "resources/background/city/layer_%d.png" // each layer frame name
	moveSpeed  = 75                                       // our movement speed
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	// robot is our robot sprite
	robot goecs.EntityID
	// gameScale is our game scale
	gameScale geometry.Scale
)

func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng *gosge.Engine) error {
	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale = eng.GetScreenSize().CalculateScale(designResolution)

	// pre-load font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	// pre-load the city layers sprites
	for l := 1; l <= numLayers; l++ {
		// generate the layer file name
		layerFile := fmt.Sprintf(layerFile, l)
		// pre-load the sprite
		if err := eng.LoadSprite(layerFile, geometry.Point{}); err != nil {
			return err
		}
	}

	// Load the sprite sheet with the robot animations
	if err := eng.LoadSpriteSheet(robotSheet); err != nil {
		return err
	}

	// get the ECS world
	world := eng.World()

	// for each city layer
	for l := int32(1); l <= numLayers; l++ {
		// generate the layer file name
		layerFile := fmt.Sprintf(layerFile, l)
		// add a sprite that will be attached to the sprite that we scroll
		attached := world.AddEntity(
			sprite.Sprite{
				Name:  layerFile,
				Scale: gameScale.Max,
			},
			geometry.Point{X: -9000}, // we move it out of the screen
			effects.Layer{Depth: float32(l)},
		)
		// add each layer to the screen
		world.AddEntity(
			sprite.Sprite{
				Name:  layerFile,
				Scale: gameScale.Max,
			},
			geometry.Point{},
			effects.Layer{Depth: float32(l)}, // each layer is have it own depth
			attach{id: attached},             // attach the layer that is attached
		)
	}

	// calculate the position for our robot
	spritePos := geometry.Point{
		X: designResolution.Width / 2 * gameScale.Point.X,
		Y: (designResolution.Height - robotPosY) * gameScale.Max,
	}

	// add the robot with it animations
	robot = world.AddEntity(
		spritePos,
		animation.Animation{
			Sequences: map[string]animation.Sequence{
				runAnim: {
					Sheet:  robotSheet,
					Base:   robotRun,
					Scale:  gameScale.Max * 0.5,
					Frames: 8,
					Delay:  0.065,
				},
				idleAnim: {
					Sheet:  robotSheet,
					Base:   robotIdle,
					Scale:  gameScale.Max * 0.5,
					Frames: 10,
					Delay:  0.065,
				},
			},
			Current: idleAnim, // current animation is idle
			Speed:   1,
		},
		effects.Layer{Depth: robotLayer},
	)

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     "press <ESC> to close, cursors to move",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       fontSize * gameScale.Max,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height * gameScale.Max,
		},
		effects.AlternateColor{
			Time: .25,
			From: color.White,
			To:   color.White.Alpha(0),
		},
		effects.Layer{Depth: fontLayer},
	)

	// system that move our robot and the layers
	world.AddSystem(robotMoveSystem)

	// listen to keys
	world.AddListener(keysListener, events.TYPE.KeyUpEvent, events.TYPE.KeyDownEvent)

	return nil
}

// a component to have an attached entity
type attach struct {
	id goecs.EntityID
}

func (a attach) Type() goecs.ComponentType {
	return attachType
}

// the goecs.ComponentType of the attach struct
var attachType = goecs.NewComponentType()

// on update
func robotMoveSystem(world *goecs.World, delta float32) error {
	// get the robot id
	robotEnt := world.Get(robot)
	// get the animation from our robot
	anim := animation.Get.Animation(robotEnt)
	// if its running
	if anim.Current == runAnim {
		// get from the ECS all the entities that have a layer a position and an attachment
		for it := world.Iterator(effects.TYPE.Layer, geometry.TYPE.Point, attachType); it != nil; it = it.Next() {
			// get the entity
			ent := it.Value()

			/// get the layer and the position
			ly := effects.Get.Layer(ent)
			pos := geometry.Get.Point(ent)

			// the direction that we are moving in 1/-1 format
			direction := float32(1)
			if anim.FlipX {
				direction = -1
			}

			// what is the width of this layer is
			layerWidth := designResolution.Width * gameScale.Max

			// moving according the direction we are facing
			pos.X -= delta * moveSpeed * (8 - ly.Depth) * direction * gameScale.Max

			// check if we need to wrap since we are completely off-screen
			if pos.X < -layerWidth {
				pos.X = 0
			} else if pos.X > layerWidth {
				pos.X = 0
			}

			// this will be attachment pos
			attachPos := geometry.Point{}

			// if the layer is moving to the right
			if pos.X > 0 {
				attachPos.X = pos.X - layerWidth // attach to the left
			} else {
				attachPos.X = pos.X + layerWidth // attach to the right
			}
			// get the robot id
			id := ent.Get(attachType).(attach).id
			world.Get(id).Set(attachPos)
			// update layer pos
			ent.Set(pos)
		}
	}
	return nil
}

// notified on events
func keysListener(world *goecs.World, e goecs.Component, _ float32) error {
	switch v := e.(type) {
	// is we get a key down event
	case events.KeyDownEvent:
		// it is the left o right cursor change animations
		if v.Key == device.KeyRight || v.Key == device.KeyLeft {
			// get the robot id
			robotEnt := world.Get(robot)
			// get the animation
			anim := animation.Get.Animation(robotEnt)

			// if we are idle
			if anim.Current == idleAnim {
				anim.Current = runAnim // run
			}
			// we flip the animation if the key was left
			anim.FlipX = v.Key == device.KeyLeft
			// update the entity animation
			robotEnt.Set(anim)
		}
	case events.KeyUpEvent:
		// it is the left o right cursor change animations
		if v.Key == device.KeyRight || v.Key == device.KeyLeft {
			// get the robot id
			robotEnt := world.Get(robot)
			// get the animation
			anim := animation.Get.Animation(robotEnt)

			// if we are running
			if anim.Current == runAnim {
				anim.Current = idleAnim // idle
			}

			// we flip the animation if the key was left
			anim.FlipX = v.Key == device.KeyLeft
			// update the entity animation
			robotEnt.Set(anim)
		}
	}
	return nil
}
