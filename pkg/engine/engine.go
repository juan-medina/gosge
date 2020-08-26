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

// Package engine contains the definition of the Engine
package engine

import (
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
)

// Engine is the interface for the game engine
type Engine interface {
	// World returns the game world.World
	World() *world.World
	// LoadTexture preloads a texture
	LoadTexture(fileName string) error
	// LoadSpriteSheet preloads a sprite.Sprite sheet
	LoadSpriteSheet(fileName string) error
	// GetSpriteSize returns the geometry.Size of a given sprite
	GetSpriteSize(sheet string, name string) (geometry.Size, error)
	// GetScreenSize returns the current screen size
	GetScreenSize() geometry.Size
	// AddGameStage adds a new game stage to our game with the given name, for changing
	//to that stage send events.ChangeStage
	AddGameStage(name string, init InitFunc)
	// SetBackgroundColor changes the current background color.Solid
	SetBackgroundColor(color color.Solid)
	// MeasureText return the geometry.Size of a string with a defined size and spacing
	MeasureText(str string, size, spacing float32) geometry.Size
}

// InitFunc is a function that will get call for our game to load
type InitFunc func(eng Engine) error
