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

package managers

import (
	"github.com/juan-medina/gosge/components"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/managers/ray"
	"github.com/juan-medina/gosge/options"
)

//DeviceManager is the interface for our device manager
type DeviceManager interface {
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
	// SetExitKey will set the exit key
	SetExitKey(key device.Key)
	// DisableExitKey will disable the exit key
	DisableExitKey()

	// GetScreenSize get the current screen size
	GetScreenSize() geometry.Size

	// GetMousePoint returns the current Point of the mouse
	GetMousePoint() geometry.Point

	// IsMouseRelease check if the given MouseButton has been release
	IsMouseRelease(button device.MouseButton) bool

	// IsMousePressed check if the given MouseButton has been pressed
	IsMousePressed(button device.MouseButton) bool

	// IsKeyPressed returns if given device.Key is pressed
	IsKeyPressed(key device.Key) bool

	// IsKeyReleased returns if given device.Key is released
	IsKeyReleased(key device.Key) bool

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
	PlayMusic(musicDef components.MusicDef, volume float32)
	// PauseMusic pauses the given components.MusicDef
	PauseMusic(musicDef components.MusicDef)
	// StopMusic stop the given components.MusicDef
	StopMusic(musicDef components.MusicDef)
	// ResumeMusic resumes the given components.MusicDef
	ResumeMusic(musicDef components.MusicDef)
	// UpdateMusic update the stream of the given components.MusicDef
	UpdateMusic(musicDef components.MusicDef)
	// ChangeMusicVolume change the given components.MusicDef volume
	ChangeMusicVolume(musicDef components.MusicDef, volume float32)
	// LoadSound giving it file name into memory
	LoadSound(fileName string) (components.SoundDef, error)
	// UnloadSound giving it file from memory
	UnloadSound(soundDef components.SoundDef)
	// PlaySound plays the given components.SoundDef
	PlaySound(soundDef components.SoundDef, volume float32)
	// SetMasterVolume change the master volume
	SetMasterVolume(volume float32)
	// StopAllSounds currently playing
	StopAllSounds()

	// SetBackgroundColor changes the current background color.Solid
	SetBackgroundColor(color color.Solid)

	// DrawText will draw a text.Text in the given geometry.Point with the correspondent color.Color
	DrawText(ftd components.FontDef, txt ui.Text, pos geometry.Point, color color.Solid)
	// DrawSprite draws a sprite.Sprite in the given geometry.Point with the tint color.Color
	DrawSprite(def components.SpriteDef, sprite sprite.Sprite, pos geometry.Point, tint color.Solid) error
	// DrawBox draws a box outline with an color.Solid and a scale
	DrawBox(pos geometry.Point, box shapes.Box, solid color.Solid)
	// DrawSolidBox draws a solid box with an color.Solid and a scale
	DrawSolidBox(pos geometry.Point, box shapes.SolidBox, solid color.Solid)
	// DrawGradientBox draws a solid box with an color.Solid and a scale
	DrawGradientBox(pos geometry.Point, box shapes.SolidBox, gradient color.Gradient)
	// MeasureText return the geometry.Size of a string with a defined size
	MeasureText(fnt components.FontDef, str string, size float32) geometry.Size
	// DrawLine between from and to with a given thickness and color.Solid
	DrawLine(from, to geometry.Point, thickness float32, color color.Solid)
	// BeginScissor start a scissor draw (define screen area for following drawing)
	BeginScissor(from geometry.Point, size geometry.Size)
	// EndScissor end the current scissor
	EndScissor()
}

// Device return the DeviceManager
func Device() DeviceManager {
	return ray.New()
}
