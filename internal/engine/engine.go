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
	"github.com/juan-medina/gosge/internal/components"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/internal/storage"
	"github.com/juan-medina/gosge/internal/systems"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/options"
)

type engineStatus int

const (
	statusInitializing = engineStatus(iota)
	statusPrepare
	statusChangeStage
	statusRunning
	statusEnding
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
	wld       *world.World
	status    engineStatus
	init      engine.InitFunc
	frameTime float32
	ds        storage.Storage
	rdr       render.Render
	stages    map[string]engine.InitFunc
}

func (ei *engineImpl) SetBackgroundColor(color color.Solid) {
	ei.rdr.SetBackgroundColor(color)
}

func (ei *engineImpl) GetScreenSize() geometry.Size {
	return ei.rdr.GetScreenSize()
}

func (ei *engineImpl) GetSpriteSize(sheet string, name string) (geometry.Size, error) {
	return ei.ds.GetSpriteSize(sheet, name)
}

func (ei *engineImpl) LoadSpriteSheet(fileName string) error {
	return ei.ds.LoadSpriteSheet(fileName)
}

func (ei *engineImpl) World() *world.World {
	return ei.wld
}

func (ei *engineImpl) Update(_ *world.World, _ float32) error {
	return nil
}

func (ei *engineImpl) Notify(_ *world.World, event interface{}, _ float32) error {
	switch v := event.(type) {
	case events.GameCloseEvent:
		ei.status = statusEnding
	case events.ChangeGameStage:
		return ei.changeStage(v.Stage)
	}
	return nil
}

// New return a engine internal implementation
func New(opt options.Options, init engine.InitFunc) Impl {
	rdr := render.New()
	return &engineImpl{
		opt:    opt,
		wld:    world.New(),
		status: statusInitializing,
		init:   init,
		ds:     storage.New(rdr),
		rdr:    rdr,
		stages: make(map[string]engine.InitFunc),
	}
}

func (ei *engineImpl) initialize() error {
	ei.rdr.Init(ei.opt)
	ei.rdr.BeginFrame()
	ei.rdr.EndFrame()

	ei.status = statusPrepare

	return nil
}

func (ei *engineImpl) prepare() error {
	ei.rdr.BeginFrame()
	ei.rdr.EndFrame()

	ei.rdr.BeginFrame()
	err := ei.init(ei)
	ei.rdr.EndFrame()

	// main systems will update before the game systems
	ei.wld.AddSystemWithPriority(ei, highPriority)
	ei.wld.AddSystemWithPriority(systems.EventSystem(ei.rdr), highPriority)

	// uis system will run after game system but before the effect systems
	ei.wld.AddSystemWithPriority(systems.UISystem(), lowPriority)

	// effect system will run after game system but before the rendering systems
	ei.wld.AddSystemWithPriority(systems.AlternateColorSystem(), lowPriority)

	// rendering system will run last
	ei.wld.AddSystemWithPriority(systems.RenderingSystem(ei.rdr, ei.ds), lastPriority)
	ei.status = statusRunning

	return err
}

func (ei *engineImpl) running() error {
	// begin frame
	ei.rdr.BeginFrame()

	// update the systems
	err := ei.wld.Update(ei.frameTime)

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
	ei.ds.Clear()
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
		case statusPrepare:
			err = ei.prepare()
		case statusChangeStage:
			err = ei.prepare()
		case statusRunning:
			err = ei.running()
		}
	}
	if err == nil && ei.status == statusEnding {
		_ = ei.end()
	}
	return err
}

func (ei *engineImpl) AddGameStage(name string, init engine.InitFunc) {
	ei.stages[name] = init
}

func (ei *engineImpl) changeStage(name string) error {
	if _, ok := ei.stages[name]; ok {
		// clear all entities and systems
		ei.wld.Clear()
		// clear all storage
		ei.ds.Clear()

		ei.init = ei.stages[name]
		ei.status = statusChangeStage

		return nil
	}
	return fmt.Errorf("stage %q not found", name)
}

func (ei engineImpl) MeasureText(font string, str string, size float32) (result geometry.Size, err error) {
	var fnt components.FontDef
	if fnt, err = ei.ds.GetFontDef(font); err == nil {
		result = ei.rdr.MeasureText(fnt, str, size)
	}

	return
}

func (ei engineImpl) LoadFont(fileName string) error {
	return ei.ds.LoadFont(fileName)
}
