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
	init   func(gWorld *world.World)
}

func newEngineState(opt options.Options, init func(gWorld *world.World)) *engineState {
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

func (es *engineState) initialize() {
	render.Init(es.opt)

	es.updateGameSettings()
	render.BeginFrame()
	es.init(es.gWorld)
	render.EndFrame()

	//es.gWorld.AddSystemToGroup(systems.xxxRenderingSystem(), renderingGroup)
	es.gWorld.AddSystemToGroup(systems.UiRenderingSystem(), uiGroup)

	es.status = statusRunning
}

func (es *engineState) render2D() {
	render.Begin2D()
	es.gWorld.UpdateGroup(renderingGroup)
	render.End2D()
}

func (es *engineState) renderUI() {
	es.gWorld.UpdateGroup(uiGroup)
}

func (es *engineState) running() {
	es.updateGameSettings()
	render.BeginFrame()

	es.gWorld.Update()
	es.render2D()
	es.renderUI()

	render.EndFrame()

	if render.ShouldClose() {
		es.status = statusEnding
	}
}

func (es *engineState) end() {
	render.End()
}

func (es *engineState) run() {
	for es.status != statusEnding {
		switch es.status {
		case statusInitializing:
			es.initialize()
		case statusRunning:
			es.running()
		case statusEnding:
			es.end()
		}
	}
}

func Run(opt options.Options, init func(gWorld *world.World)) {
	newEngineState(opt, init).run()
}
