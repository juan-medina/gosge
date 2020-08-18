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

package engine

import (
	"github.com/juan-medina/goecs/pkg/entitiy"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/options"
	"github.com/juan-medina/gosge/pkg/render"
	"github.com/juan-medina/gosge/pkg/systems"
	"math"
)

type engineStatus int

const (
	statusInitializing = engineStatus(iota)
	statusRunning
	statusEnding
)

const (
	uiGroup        = "UI_GROUP"
	renderingGroup = "RENDERING_GROUP"
)

type engineState struct {
	opt    options.Options
	gWorld *world.World
	status engineStatus
	init   func(gWorld *world.World) error
}

func newEngineState(opt options.Options, init func(gWorld *world.World) error) *engineState {
	return &engineState{
		opt:    opt,
		gWorld: world.New(),
		status: statusInitializing,
		init:   init,
	}
}

func (es *engineState) updateGameSettings() {
	gsEnt := es.gWorld.Entity(components.GameSettingsType)
	if gsEnt == nil {
		gsEnt = es.gWorld.Add(entitiy.New(components.GameSettings{}))
	}
	w, h := render.GetScreenSize()
	gs := components.GameSettings{
		Current: struct {
			Width  int
			Height int
		}{
			Width:  w,
			Height: h,
		},
		Original: struct {
			Width  int
			Height int
		}{
			Width:  es.opt.Width,
			Height: es.opt.Height,
		},
	}
	sx := float64(gs.Current.Width) / float64(gs.Original.Width)
	sy := float64(gs.Current.Height) / float64(gs.Original.Height)
	gs.Scale = math.Min(sx, sy)

	gsEnt.Add(gs)
}

func (es *engineState) initialize() error {
	render.Init(es.opt)

	es.updateGameSettings()
	render.BeginFrame()
	err := es.init(es.gWorld)
	render.EndFrame()

	if err == nil {
		//es.gWorld.AddSystemToGroup(systems.xxxRenderingSystem(), renderingGroup)
		es.gWorld.AddSystemToGroup(systems.UiRenderingSystem(), uiGroup)

		es.status = statusRunning
	}
	return err
}

func (es *engineState) render2D() error {
	render.Begin2D()
	err := es.gWorld.UpdateGroup(renderingGroup, 0)
	render.End2D()
	return err
}

func (es *engineState) renderUI() error {
	return es.gWorld.UpdateGroup(uiGroup, 0)
}

func (es *engineState) running() error {
	es.updateGameSettings()
	render.BeginFrame()

	var err error
	if err = es.gWorld.Update(0); err == nil {
		if err = es.render2D(); err == nil {
			err = es.renderUI()
		}
	}

	render.EndFrame()

	if err == nil {
		if render.ShouldClose() {
			es.status = statusEnding
		}
	}
	return err
}

func (es *engineState) end() error {
	render.End()
	return nil
}

func (es *engineState) run() error {
	var err error = nil
	for es.status != statusEnding && err == nil {
		switch es.status {
		case statusInitializing:
			err = es.initialize()
		case statusRunning:
			err = es.running()
		case statusEnding:
			err = es.end()
		}
	}
	if err != nil && es.status == statusRunning {
		_ = es.end()
	}
	return err
}

func Run(opt options.Options, init func(gWorld *world.World) error) error {
	return newEngineState(opt, init).run()
}
