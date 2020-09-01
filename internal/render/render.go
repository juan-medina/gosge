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
	"github.com/juan-medina/gosge/pkg/components/device"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/ui"
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

	// GetMousePoint returns the current Point of the mouse
	GetMousePoint() geometry.Point

	// IsMouseRelease check if the given MouseButton has been release
	IsMouseRelease(button device.MouseButton) bool

	// GetKeyStatus returns the device.KeyStatus for a given device.Key
	GetKeyStatus(key device.Key) device.KeyStatus

	// GetFrameTime returns the time from the delta time for current frame
	GetFrameTime() float32

	// LoadTexture giving it file name into VRAM
	LoadTexture(fileName string) (components.TextureDef, error)
	// UnloadTexture from VRAM
	UnloadTexture(textureDef components.TextureDef)

	// LoadFont giving it file name into VRAM
	LoadFont(fileName string) (components.FontDef, error)
	// UnloadFont from VRAM
	UnloadFont(textureDef components.FontDef)

	// LoadMusic giving it file name into memory
	LoadMusic(fileName string) (components.MusicDef, error)
	// UnloadMusic giving it file from memory
	UnloadMusic(musicDef components.MusicDef)
	// PlayMusic plays the given components.MusicDef
	PlayMusic(musicDef components.MusicDef, loops int32)
	// PauseMusic pauses the given components.MusicDef
	PauseMusic(musicDef components.MusicDef)
	// StopMusic stop the given components.MusicDef
	StopMusic(musicDef components.MusicDef)
	// ResumeMusic resumes the given components.MusicDef
	ResumeMusic(musicDef components.MusicDef)
	// UpdateMusic update the stream of the given components.MusicDef
	UpdateMusic(musicDef components.MusicDef)

	// SetBackgroundColor changes the current background color.Solid
	SetBackgroundColor(color color.Solid)

	// DrawText will draw a text.Text in the given geometry.Point with the correspondent color.Color
	DrawText(ftd components.FontDef, txt ui.Text, pos geometry.Point, color color.Solid)
	// DrawSprite draws a sprite.Sprite in the given geometry.Point with the tint color.Color
	DrawSprite(def components.SpriteDef, sprite sprite.Sprite, pos geometry.Point, tint color.Solid) error
	// DrawSolidBox draws a solid box with an color.Solid and a scale
	DrawSolidBox(pos geometry.Point, box shapes.Box, solid color.Solid)
	// DrawGradientBox draws a solid box with an color.Solid and a scale
	DrawGradientBox(pos geometry.Point, box shapes.Box, gradient color.Gradient)
	// MeasureText return the geometry.Size of a string with a defined size
	MeasureText(fnt components.FontDef, str string, size float32) geometry.Size
}

// New return the Render system
func New() Render {
	return ray.New()
}
