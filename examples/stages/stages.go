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
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/engine"
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
	eng.AddGameStage("main", mainStage{})
	eng.AddGameStage("menu", menuStage{})

	scs := &stageChangeSystem{eng: eng, stages: []string{"main", "menu"}}
	eng.World().AddSystem(scs)

	return eng.ChangeStage("menu")
}

type mainStage struct{}

func (m mainStage) Load(_ engine.Engine) error {
	return nil
}

func (m mainStage) Run(eng engine.Engine) error {
	gWorld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// add the centered text
	gWorld.Add(entity.New(
		text.Text{
			String:     "Main Stage",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.MiddleVAlignment,
			Size:       300 * gameScale.Min,
			Spacing:    10,
		},
		geometry.Position{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		effects.AlternateColor{
			Time:  1,
			Delay: 0.5,
			From:  color.Red,
			To:    color.Yellow,
		},
	))

	return nil
}

func (m mainStage) Unload(eng engine.Engine) error {
	w := eng.World()
	for it := w.Iterator(geometry.TYPE.Position); it.HasNext(); {
		_ = w.Remove(it.Value())
	}
	return nil
}

type menuStage struct{}

func (m menuStage) Load(_ engine.Engine) error {
	return nil
}

func (m menuStage) Run(eng engine.Engine) error {
	gWorld := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// add the centered text
	gWorld.Add(entity.New(
		text.Text{
			String:     "Menu",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.MiddleVAlignment,
			Size:       300 * gameScale.Min,
			Spacing:    10,
		},
		geometry.Position{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		effects.AlternateColor{
			Time:  1,
			Delay: 0.5,
			From:  color.Red,
			To:    color.Yellow,
		},
	))

	return eng.ChangeStage("main")
}

func (m menuStage) Unload(eng engine.Engine) error {
	w := eng.World()
	for it := w.Iterator(geometry.TYPE.Position); it.HasNext(); {
		_ = w.Remove(it.Value())
	}
	return nil
}

type stageChangeSystem struct {
	eng     engine.Engine
	time    float32
	stages  []string
	current int
}

func (s *stageChangeSystem) Update(_ *world.World, delta float32) error {
	s.time += delta

	if s.time > 5 {
		s.time = 0
		s.current++
		if s.current >= len(s.stages) {
			s.current = 0
		}
		return s.eng.ChangeStage(s.stages[s.current])
	}
	return nil
}

func (s stageChangeSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}
