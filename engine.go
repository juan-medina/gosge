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

package gosge

import (
	"fmt"
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/managers"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
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

// InitFunc is a function that will get call for our game to load
type InitFunc func(eng *Engine) error

// Engine is our game engine
type Engine struct {
	opt       options.Options
	world     *goecs.World
	status    engineStatus
	init      InitFunc
	frameTime float32
	sm        managers.StorageManager
	dm        managers.DeviceManager
	stages    map[string]InitFunc
}

// SetBackgroundColor changes the current background color.Solid
func (e *Engine) SetBackgroundColor(color color.Solid) {
	e.opt.BackGround = color
	e.dm.SetBackgroundColor(color)
}

// GetScreenSize returns the current screen size
func (e *Engine) GetScreenSize() geometry.Size {
	return e.dm.GetScreenSize()
}

// GetSpriteSize returns the geometry.Size of a given sprite
func (e *Engine) GetSpriteSize(sheet string, name string) (geometry.Size, error) {
	return e.sm.GetSpriteSize(sheet, name)
}

// LoadSpriteSheet preloads a sprite.
func (e *Engine) LoadSpriteSheet(fileName string) error {
	return e.sm.LoadSpriteSheet(fileName)
}

// World returns the game goecs.World
func (e *Engine) World() *goecs.World {
	return e.world
}

// Listener listen to game events
func (e *Engine) Listener(_ *goecs.World, event interface{}, _ float32) error {
	switch v := event.(type) {
	case events.GameCloseEvent:
		e.status = statusEnding
	case events.ChangeGameStage:
		return e.changeStage(v.Stage)
	}
	return nil
}

func (e *Engine) initialize() error {
	log.Info().Interface("options", e.opt).Msg("Initializing Engine")
	e.dm.Init(e.opt)
	e.dm.BeginFrame()
	e.dm.EndFrame()

	ss := e.dm.GetScreenSize()
	log.Trace().Interface("size", ss).Msg("Screen ready")

	e.status = statusPrepare

	return nil
}

func (e *Engine) drawLoading() {
	time := float32(.125)
	clr := e.opt.BackGround
	e.dm.SetBackgroundColor(color.Black)
	for time > 0 {
		e.dm.BeginFrame()
		e.dm.EndFrame()
		time -= e.dm.GetFrameTime()
	}
	e.dm.SetBackgroundColor(clr)
}

func (e *Engine) register(mng managers.Manager, priority int32) {
	switch m := mng.(type) {
	case managers.WithSystemAndListener:
		e.world.AddSystemWithPriority(m.System, priority)
		e.world.AddListenerWithPriority(m.Listener, priority)
	case managers.WithSystem:
		e.world.AddSystemWithPriority(m.System, priority)
	case managers.WithListener:
		e.world.AddListenerWithPriority(m.Listener, priority)
	default:
		panic("can not register manager")
	}
}

func (e *Engine) prepare() error {
	e.drawLoading()
	err := e.init(e)

	// main managers will update before the game managers
	e.register(e, highPriority)
	e.register(managers.Events(e.dm), highPriority)

	// add the sound manager
	e.register(managers.Sounds(e.dm, e.sm), highPriority)

	// add the music manager
	e.register(managers.Music(e.dm, e.sm), highPriority)

	// animation system will run after game systems but before the rendering managers
	e.register(managers.Animation(), lowPriority)

	// tiled manager will run after game systems but before the rendering managers
	e.register(managers.TiledMaps(e.sm), lowPriority)

	// ui manager will run after game system but before the effect managers
	e.register(managers.UI(e.dm, e.sm), lowPriority)

	// effects manager will run after game system but before the rendering managers
	e.register(managers.Effects(), lowPriority)

	// rendering manager will run last
	e.register(managers.Rendering(e.dm, e.sm), lastPriority)
	e.status = statusRunning

	return err
}

func (e *Engine) running() error {
	// begin frame
	e.dm.BeginFrame()

	// update the systems
	err := e.world.Update(e.frameTime)

	// we end the frame regardless of if we have an error
	e.dm.EndFrame()

	// if we do not have an error check if we should close
	if err == nil {
		if e.dm.ShouldClose() {
			e.status = statusEnding
		}
	}

	return err
}

func (e *Engine) end() error {
	e.sm.Clear()
	e.dm.End()
	return nil
}

// Run a game within the engine
func (e *Engine) Run() error {
	var err error = nil
	for e.status != statusEnding && err == nil {
		e.frameTime = e.dm.GetFrameTime()
		switch e.status {
		case statusInitializing:
			err = e.initialize()
		case statusPrepare:
			err = e.prepare()
		case statusChangeStage:
			err = e.prepare()
		case statusRunning:
			err = e.running()
		}
	}
	if err == nil && e.status == statusEnding {
		_ = e.end()
	}
	return err
}

// AddGameStage adds a new game stage to our game with the given name, for changing
// to that stage send a events.ChangeStage
func (e *Engine) AddGameStage(name string, init InitFunc) {
	e.stages[name] = init
}

func (e *Engine) changeStage(name string) error {
	if _, ok := e.stages[name]; ok {
		e.dm.StopAllSounds()
		// clear all entities and systems
		e.world.Clear()
		// clear all storage
		e.sm.Clear()

		e.init = e.stages[name]
		e.status = statusChangeStage

		return nil
	}
	return fmt.Errorf("stage %q not found", name)
}

// MeasureText return the geometry.Size of a string with a defined size and spacing
func (e Engine) MeasureText(font string, str string, size float32) (result geometry.Size, err error) {
	var fnt components.FontDef
	if fnt, err = e.sm.GetFontDef(font); err == nil {
		result = e.dm.MeasureText(fnt, str, size)
	}

	return
}

// LoadFont preloads a font
func (e Engine) LoadFont(fileName string) error {
	return e.sm.LoadFont(fileName)
}

// LoadSprite preloads a single sprite.Sprite
func (e Engine) LoadSprite(filename string, pivot geometry.Point) error {
	return e.sm.LoadSingleSprite("", filename, pivot)
}

// SpriteAtContains indicates if a sprite.Sprite at a given geometry.Point contains a geometry.Point
func (e Engine) SpriteAtContains(spr sprite.Sprite, at geometry.Point, point geometry.Point) bool {
	return e.sm.SpriteAtContains(spr, at, point)
}

// LoadMusic preloads a music stream
func (e Engine) LoadMusic(filename string) error {
	return e.sm.LoadMusic(filename)
}

//LoadSound preload a sound wave
func (e Engine) LoadSound(filename string) error {
	return e.sm.LoadSound(filename)
}

// LoadTiledMap preload a tiled map
func (e Engine) LoadTiledMap(filename string) error {
	return e.sm.LoadTiledMap(filename)
}

// GeTiledMapSize returns the geometry.Size of a given tile map
func (e Engine) GeTiledMapSize(name string) (size geometry.Size, err error) {
	var def components.TiledMapDef
	if def, err = e.sm.GetTiledMapDef(name); err == nil {
		size = def.Size
	}
	return
}

// New return a Engine
func New(opt options.Options, init InitFunc) *Engine {
	dm := managers.Device()
	return &Engine{
		opt:    opt,
		world:  goecs.Default(),
		status: statusInitializing,
		init:   init,
		sm:     managers.Storage(dm),
		dm:     dm,
		stages: make(map[string]InitFunc),
	}
}
