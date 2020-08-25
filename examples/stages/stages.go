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
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"reflect"
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

func mainStage(eng engine.Engine) error {
	// pre load sprites
	if err := eng.LoadSpriteSheet("resources/stages.json"); err != nil {
		return err
	}

	gWorld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	eng.SetBackgroundColor(color.SkyBlue)

	// add the centered text
	gWorld.Add(entity.New(
		text.Text{
			String:     "Main Stage",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.TopVAlignment,
			Size:       300 * gameScale.Min,
			Spacing:    10,
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

	gWorld.Add(entity.New(
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

	wld := eng.World()
	wld.AddSystem(&stageChangeSystem{stage: "menu", time: 5})

	return nil
}

func menuStage(eng engine.Engine) error {
	// pre load sprites
	if err := eng.LoadSpriteSheet("resources/stages.json"); err != nil {
		return err
	}

	gWorld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	eng.SetBackgroundColor(color.Gopher)

	// add the centered text
	gWorld.Add(entity.New(
		text.Text{
			String:     "Menu",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.TopVAlignment,
			Size:       300 * gameScale.Min,
			Spacing:    10,
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

	gWorld.Add(entity.New(
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

	createButton(gWorld, designResolution.Width/2, designResolution.Height-200, 400, 100, gameScale,
		color.SkyBlue, color.Yellow, "Play!")

	createButton(gWorld, designResolution.Width/2, designResolution.Height-90, 400, 100, gameScale,
		color.Beige, color.Yellow, "Exit")

	gWorld.AddSystem(newButtonOverSystem())

	//eng.World().AddSystem(&stageChangeSystem{stage: "main", time: 5})

	return nil
}

type buttonColor struct {
	normal color.Solid
	hover  color.Solid
}

type buttonText struct {
	normal color.Solid
	hover  color.Solid
	ref    *entity.Entity
}
type button struct {
	color buttonColor
	txt   buttonText
}

var (
	tagType = reflect.TypeOf(button{})
)

func createButton(wld *world.World, x, y, w, h float32, scale geometry.Scale, bc, tc color.Solid, str string) {
	boxPos := geometry.Point{
		X: (x - (w / 2)) * scale.Point.X,
		Y: (y - (h / 2)) * scale.Point.Y,
	}

	textPos := geometry.Point{
		X: boxPos.X + ((w / 2) * scale.Point.X),
		Y: boxPos.Y + ((h / 2) * scale.Point.Y),
	}

	txt := wld.Add(entity.New(
		text.Text{
			String:     str,
			Size:       (h - 20) * scale.Min,
			Spacing:    10 * scale.Min,
			VAlignment: text.MiddleVAlignment,
			HAlignment: text.CenterHAlignment,
		},
		color.Yellow,
		textPos,
		effects.Layer{Depth: -11},
	))

	bnc := bc.Blend(color.Black, 0.2)
	tnc := tc.Blend(color.Black, 0.2)

	wld.Add(entity.New(
		boxPos,
		shapes.Box{
			Size: geometry.Size{
				Width:  w,
				Height: h,
			},
			Scale: scale.Min,
		},
		color.DarkBlue,
		effects.Layer{Depth: -10},
		button{
			color: buttonColor{normal: bnc, hover: bc},
			txt:   buttonText{normal: tnc, hover: tc, ref: txt},
		},
	))

	boxPos.X += 10
	boxPos.Y += 10

	wld.Add(entity.New(
		boxPos,
		shapes.Box{
			Size: geometry.Size{
				Width:  w,
				Height: h,
			},
			Scale: scale.Min,
		},
		color.DarkGray,
		effects.Layer{Depth: 0},
	))
}

// stageChangeSystem is a world.System that will automatically change to game stage in a giving time
type stageChangeSystem struct {
	time  float32 // time to change to the given stage
	stage string  // stage name, must be create with engine.AddGameStage
}

func (s *stageChangeSystem) Update(world *world.World, delta float32) error {
	// we discount time until zero, them we will notify to change the stage
	if s.time -= delta; s.time <= 0 {
		return world.Notify(events.ChangeGameStage{Stage: s.stage})
	}
	return nil
}

func (s stageChangeSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}

func newButtonOverSystem() world.System {
	return &buttonOverSystem{}
}

type buttonOverSystem struct{}

func (bos buttonOverSystem) Update(_ *world.World, _ float32) error {
	return nil
}

func (bos *buttonOverSystem) Notify(wld *world.World, e interface{}, _ float32) error {
	switch v := e.(type) {
	case events.MouseMoveEvent:
		for it := wld.Iterator(tagType, color.TYPE.Solid); it.HasNext(); {
			ent := it.Value()
			btn := ent.Get(tagType).(button)

			pos := geometry.Get.Point(ent)
			box := shapes.Get.Box(ent)

			clr := btn.color.normal
			tcl := btn.txt.normal

			if box.Contains(pos, v.Point) {
				clr = btn.color.hover
				tcl = btn.txt.hover
			}

			ent.Set(clr)
			btn.txt.ref.Set(tcl)
		}
	}
	return nil
}
