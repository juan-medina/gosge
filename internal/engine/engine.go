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

// Package engine internal engine implementation
package engine

import (
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/internal/systems"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/options"
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

// Impl a engine internal implementation
type Impl interface {
	// Run the engine
	Run() error
}

type engineImpl struct {
	opt       options.Options
	gWorld    *world.World
	status    engineStatus
	init      engine.InitFunc
	frameTime float64
}

func (ei *engineImpl) LoadTexture(fileName string) error {
	return render.LoadTexture(fileName)
}

func (ei *engineImpl) World() *world.World {
	return ei.gWorld
}

func (ei *engineImpl) Update(_ *world.World, _ float64) error {
	return nil
}

func (ei *engineImpl) Notify(_ *world.World, event interface{}, _ float64) error {
	switch event.(type) {
	case events.GameCloseEvent:
		ei.status = statusEnding
	}
	return nil
}

// New return a engine internal implementation
func New(opt options.Options, init engine.InitFunc) Impl {
	return &engineImpl{
		opt:    opt,
		gWorld: world.New(),
		status: statusInitializing,
		init:   init,
	}
}

func (ei *engineImpl) initialize() error {
	render.Init(ei.opt)
	render.BeginFrame()
	err := ei.init(ei)
	render.EndFrame()

	if err == nil {
		ei.gWorld.AddSystem(ei)
		ei.gWorld.AddSystem(systems.EventSystem())
		ei.gWorld.AddSystem(systems.AlternateColorSystem())
		ei.gWorld.AddSystemToGroup(systems.SpriteRenderingSystem(), renderingGroup)
		ei.gWorld.AddSystemToGroup(systems.UIRenderingSystem(), uiGroup)

		ei.status = statusRunning
	}
	return err
}

func (ei *engineImpl) render() error {
	err := ei.gWorld.UpdateGroup(renderingGroup, ei.frameTime)
	return err
}

func (ei *engineImpl) renderUI() error {
	return ei.gWorld.UpdateGroup(uiGroup, ei.frameTime)
}

func (ei *engineImpl) running() error {
	render.BeginFrame()

	var err error
	if err = ei.gWorld.Update(ei.frameTime); err == nil {
		if err = ei.render(); err == nil {
			err = ei.renderUI()
		}
	}

	render.EndFrame()

	if err == nil {
		if render.ShouldClose() {
			ei.status = statusEnding
		}
	}
	return err
}

func (ei *engineImpl) end() error {
	render.End()
	return nil
}

func (ei *engineImpl) Run() error {
	var err error = nil
	for ei.status != statusEnding && err == nil {
		ei.frameTime = render.GetFrameTime()
		switch ei.status {
		case statusInitializing:
			err = ei.initialize()
		case statusRunning:
			err = ei.running()
		case statusEnding:
			err = ei.end()
		}
	}
	if err != nil && ei.status == statusRunning {
		_ = ei.end()
	}
	return err
}
