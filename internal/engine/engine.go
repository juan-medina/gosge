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
	"fmt"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/internal/store"
	"github.com/juan-medina/gosge/internal/systems"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/options"
)

type engineStatus int

const (
	statusInitializing = engineStatus(iota)
	statusRunning
	statusEnding
	statusUnloadCurrentStage
	statusLoadNextStage
	statusRunStage
)

const (
	lastPriority = int32(-1000)
	lowPriority  = int32(-500)
	highPriority = int32(500)
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
	frameTime float32
	ss        store.SpriteStorage
	rdr       render.Render
	stages    map[string]engine.GameStage
	cStage    engine.GameStage
	nStage    engine.GameStage
}

func (ei *engineImpl) GetScreenSize() geometry.Size {
	return ei.rdr.GetScreenSize()
}

func (ei *engineImpl) GetSpriteSize(sheet string, name string) (geometry.Size, error) {
	return ei.ss.GetSpriteSize(sheet, name)
}

func (ei *engineImpl) LoadSpriteSheet(fileName string) error {
	return ei.ss.LoadSpriteSheet(fileName)
}

func (ei *engineImpl) LoadTexture(fileName string) error {
	return ei.rdr.LoadTexture(fileName)
}

func (ei *engineImpl) World() *world.World {
	return ei.gWorld
}

func (ei *engineImpl) Update(_ *world.World, _ float32) error {
	return nil
}

func (ei *engineImpl) Notify(_ *world.World, event interface{}, _ float32) error {
	switch event.(type) {
	case events.GameCloseEvent:
		ei.status = statusEnding
	}
	return nil
}

// New return a engine internal implementation
func New(opt options.Options, init engine.InitFunc) Impl {
	rdr := render.New()
	return &engineImpl{
		opt:    opt,
		gWorld: world.New(),
		status: statusInitializing,
		init:   init,
		ss:     store.NewSpriteStorage(rdr),
		rdr:    rdr,
		stages: make(map[string]engine.GameStage),
		cStage: nil,
		nStage: nil,
	}
}

func (ei *engineImpl) initialize() error {
	ei.rdr.Init(ei.opt)
	ei.rdr.BeginFrame()
	err := ei.init(ei)
	ei.rdr.EndFrame()

	if err == nil {
		// main systems will update before the game systems
		ei.gWorld.AddSystemWithPriority(ei, highPriority)
		ei.gWorld.AddSystemWithPriority(systems.EventSystem(ei.rdr), highPriority)

		// effect system will run after game system but before the rendering systems
		ei.gWorld.AddSystemWithPriority(systems.AlternateColorSystem(), lowPriority)

		// rendering system will run last
		ei.gWorld.AddSystemWithPriority(systems.RenderingSystem(ei.rdr, ei.ss), lastPriority)

		if ei.nStage == nil {
			ei.status = statusRunning
		} else {
			ei.status = statusLoadNextStage
		}

	}
	return err
}

func (ei *engineImpl) running() error {
	// begin frame
	ei.rdr.BeginFrame()

	// update the systems
	err := ei.gWorld.Update(ei.frameTime)

	// we end the frame regardless of if we have an error
	ei.rdr.EndFrame()

	// if we do not have an error check if we should close
	if err == nil {
		if ei.rdr.ShouldClose() {
			ei.status = statusEnding
		}
	}

	return err
}

func (ei *engineImpl) end() error {
	ei.rdr.End()
	return nil
}

func (ei *engineImpl) Run() error {
	var err error = nil
	for ei.status != statusEnding && err == nil {
		ei.frameTime = ei.rdr.GetFrameTime()
		switch ei.status {
		case statusInitializing:
			err = ei.initialize()
		case statusRunning:
			err = ei.running()
		case statusEnding:
			err = ei.end()
		case statusUnloadCurrentStage:
			err = ei.unloadCurrentStage()
		case statusLoadNextStage:
			err = ei.loadNextStage()
		case statusRunStage:
			err = ei.runStage()
		}
	}
	if err != nil && ei.status == statusRunning {
		_ = ei.end()
	}
	return err
}

func (ei *engineImpl) AddGameStage(name string, stage engine.GameStage) {
	ei.stages[name] = stage
}

func (ei *engineImpl) ChangeStage(name string) error {
	if _, ok := ei.stages[name]; ok {
		ei.nStage = ei.stages[name]

		// if we have a current stage unload if not load the new
		if ei.cStage != nil {
			ei.status = statusUnloadCurrentStage
		} else {
			ei.status = statusLoadNextStage
		}

		return nil
	}
	return fmt.Errorf("stage %q not found", name)
}

func (ei *engineImpl) unloadCurrentStage() error {
	var err error

	// begin frame
	ei.rdr.BeginFrame()

	err = ei.cStage.Unload(ei)

	if err == nil {
		ei.cStage = nil
		ei.status = statusLoadNextStage
	} else {
		ei.status = statusRunning
	}

	// we end the frame regardless of if we have an error
	ei.rdr.EndFrame()

	return err
}

func (ei *engineImpl) loadNextStage() error {
	var err error

	// begin frame
	ei.rdr.BeginFrame()

	err = ei.nStage.Load(ei)
	if err == nil {
		ei.cStage = ei.nStage
		ei.status = statusRunStage
	} else {
		ei.status = statusRunning
	}

	// we end the frame regardless of if we have an error
	ei.rdr.EndFrame()

	return err
}

func (ei *engineImpl) runStage() error {
	var err error

	// begin frame
	ei.rdr.BeginFrame()

	if ei.cStage != nil {
		err = ei.cStage.Run(ei)
	}

	ei.status = statusRunning

	// we end the frame regardless of if we have an error
	ei.rdr.EndFrame()

	return err
}
