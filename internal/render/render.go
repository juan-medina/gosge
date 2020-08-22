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

// Package render device rendering implementation
package render

import (
	"github.com/juan-medina/gosge/internal/components"
	"github.com/juan-medina/gosge/internal/render/ray"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/options"
)

//Render is the interface for our rendering system
type Render interface {
	// Init the rendering device
	Init(opt options.Options)
	// End the rendering device
	End()
	// BeginFrame for rendering
	BeginFrame()
	// EndFrame for rendering
	EndFrame()
	// ShouldClose returns if th engine should close
	ShouldClose() bool

	// GetScreenSize get the current screen size
	GetScreenSize() geometry.Size
	// GetMousePosition returns the current position of the mouse
	GetMousePosition() geometry.Position
	// IsScreenScaleChange returns if the current screen scale has changed
	IsScreenScaleChange() bool

	// GetFrameTime returns the time from the delta time for current frame
	GetFrameTime() float32

	// LoadTexture giving it file name into VRAM
	LoadTexture(fileName string) error
	// UnloadAllTextures from VRAM
	UnloadAllTextures()
	// DrawText will draw a text.Text in the given geometry.Position with the correspondent color.Color
	DrawText(txt text.Text, pos geometry.Position, color color.Solid)
	// DrawSprite draws a sprite.Sprite in the given geometry.Position with the tint color.Color
	DrawSprite(def components.SpriteDef, sprite sprite.Sprite, pos geometry.Position, tint color.Solid) error
	// DrawSolidBox draws a solid box with an color.Solid and a scale
	DrawSolidBox(pos geometry.Position, box shapes.Box, solid color.Solid)
	// DrawGradientBox draws a solid box with an color.Solid and a scale
	DrawGradientBox(pos geometry.Position, box shapes.Box, gradient color.Gradient)
}

// New return the Render system
func New() Render {
	return ray.New()
}
