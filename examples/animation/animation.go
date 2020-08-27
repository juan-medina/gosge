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
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/device"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"reflect"
)

var opt = options.Options{
	Title:      "AnimationFrame Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
}

const (
	robotSheet = "resources/robot.json"
	robotIdle  = "Idle_%02d.png"
	robotRun   = "Run_%02d.png"
	runAnim    = "run"
	idleAnim   = "idle"
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	// robot is our robot sprite
	robot *entity.Entity
)

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

func loadGame(eng engine.Engine) error {
	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	if err := eng.LoadSpriteSheet(robotSheet); err != nil {
		return err
	}

	wld := eng.World()

	spritePos := geometry.Point{
		X: designResolution.Width / 2 * gameScale.Point.X,
		Y: (designResolution.Height - 150) * gameScale.Point.Y,
	}

	robot = wld.Add(entity.New(
		spritePos,
		Animation{
			Frames: map[string]AnimationFrame{
				runAnim: {
					Sheet:  robotSheet,
					Base:   robotRun,
					Scale:  gameScale.Min * 0.5,
					Frames: 8,
					Delay:  0.065,
				},
				idleAnim: {
					Sheet:  robotSheet,
					Base:   robotIdle,
					Scale:  gameScale.Min * 0.5,
					Frames: 10,
					Delay:  0.065,
				},
			},
			Current: idleAnim,
		},
	))

	wld.AddSystem(NewAnimationSystem())

	return nil
}

// AnimationFrame represent a set of frames that will be render with a delay
type AnimationFrame struct {
	Sheet    string  // Sheet is the sprite sheet where the animation sprites are
	Base     string  // Base is the base name for each frame. ex : Idle_%d.png
	Rotation float32 // Rotation for this AnimationFrame
	Scale    float32 // Scale for this AnimationFrame
	Frames   int32   // Frames are the number of frame in this animation
	Delay    float32 // Delay number of seconds to wait in each frame
	Current  float32 // Current is the time in this frame
	Frame    int32   // Frame is the current frame number
}

// AnimationFrameType is the reflect.Type for an AnimationFrame
var AnimationFrameType = reflect.TypeOf(AnimationFrame{})

// Animation allow to easily switch animations
type Animation struct {
	Frames  map[string]AnimationFrame // AnimationFrame is the different AnimationFrame for this set
	Current string                    // Current is the animation that like to run
	Running string                    // Running is the animation that is running
}

// AnimationType is the reflect.Type for an Animation
var AnimationType = reflect.TypeOf(Animation{})

type animationSystem struct{}

func (as *animationSystem) Update(world *world.World, delta float32) error {
	for it := world.Iterator(AnimationType); it.HasNext(); {
		ent := it.Value()

		anim := ent.Get(AnimationType).(Animation)

		if anim.Running != anim.Current {
			anim.Running = anim.Current
			if _, ok := anim.Frames[anim.Current]; !ok {
				return fmt.Errorf("can not find animation: %q", anim.Current)
			}
			ent.Set(anim.Frames[anim.Current])
			ent.Set(anim)
		}
	}

	for it := world.Iterator(AnimationFrameType); it.HasNext(); {
		ent := it.Value()

		anf := ent.Get(AnimationFrameType).(AnimationFrame)

		if anf.Current += delta; anf.Current > anf.Delay {
			anf.Current = 0
			if anf.Frame++; anf.Frame >= anf.Frames {
				anf.Frame = 0
			}
		}

		spriteName := fmt.Sprintf(anf.Base, anf.Frame+1)

		spr := sprite.Sprite{
			Sheet:    anf.Sheet,
			Name:     spriteName,
			Scale:    anf.Scale,
			Rotation: anf.Rotation,
		}

		ent.Set(spr)
		ent.Set(anf)
	}

	return nil
}

func (as animationSystem) Notify(_ *world.World, e interface{}, _ float32) error {

	switch v := e.(type) {
	case events.MouseUpEvent:
		if v.MouseButton == device.MouseLeftButton {
			anim := robot.Get(AnimationType).(Animation)

			if anim.Current == idleAnim {
				anim.Current = runAnim
			} else {
				anim.Current = idleAnim
			}

			robot.Set(anim)
		}
	}
	return nil
}

// NewAnimationSystem creates a animation world.System
func NewAnimationSystem() world.System {
	return &animationSystem{}
}
