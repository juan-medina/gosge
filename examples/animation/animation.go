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
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/device"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
)

var opt = options.Options{
	Title:      "Animation Game",
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
		animation.Animation{
			Sequences: map[string]animation.Sequence{
				runAnim: {
					Sheet:  robotSheet,
					Base:   robotRun,
					Scale:  gameScale.Min,
					Frames: 8,
					Delay:  0.065,
				},
				idleAnim: {
					Sheet:  robotSheet,
					Base:   robotIdle,
					Scale:  gameScale.Min,
					Frames: 10,
					Delay:  0.065,
				},
			},
			Current: idleAnim,
			Speed:   1,
		},
	))

	wld.AddSystem(RobotMoveSystem())

	return nil
}

type robotMoveSystem struct{}

func (rms *robotMoveSystem) Update(_ *world.World, _ float32) error {
	return nil
}

func (rms robotMoveSystem) Notify(_ *world.World, e interface{}, _ float32) error {
	switch v := e.(type) {
	case events.KeyEvent:
		if v.Key == device.KeyRight || v.Key == device.KeyLeft {
			anim := animation.Get.Animation(robot)
			if v.Status.Down {
				if anim.Current == idleAnim {
					anim.Current = runAnim
				}
			} else if v.Status.Released {
				if anim.Current == runAnim {
					anim.Current = idleAnim
				}
			}
			anim.FlipX = v.Key == device.KeyLeft
			robot.Set(anim)
		}
	}
	return nil
}

// RobotMoveSystem creates a system to move our robot
func RobotMoveSystem() world.System {
	return &robotMoveSystem{}
}
