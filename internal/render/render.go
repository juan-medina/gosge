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
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/juan-medina/gosge/internal/components"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/options"
)

var saveOpts = options.Options{}

// Init the rendering device
func Init(opt options.Options) {
	saveOpts = opt
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(opt.Size.Width), int32(opt.Size.Height), opt.Title)
}

// End the rendering device
func End() {
	UnloadAllTextures()
	rl.CloseWindow()
}

// BeginFrame for rendering
func BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(color2RayColor(saveOpts.ClearColor))
}

// EndFrame for rendering
func EndFrame() {
	rl.EndDrawing()
}

// ShouldClose returns if th engine should close
func ShouldClose() bool {
	return rl.WindowShouldClose()
}

// GetScreenSize get the current screen size
func GetScreenSize() (width int, height int) {
	return rl.GetScreenWidth(), rl.GetScreenHeight()
}

// DrawText will draw a text.Text in the given geometry.Position with the correspondent color.Color
func DrawText(txt text.Text, pos geometry.Position, color color.Color) {
	font := rl.GetFontDefault()

	vec := rl.Vector2{
		X: pos.X,
		Y: pos.Y,
	}

	if txt.HAlignment != text.LeftHAlignment || txt.VAlignment != text.BottomVAlignment {
		av := rl.MeasureTextEx(font, txt.String, txt.Size, txt.Spacing)

		switch txt.HAlignment {
		case text.LeftHAlignment:
			av.X = 0
		case text.CenterHAlignment:
			av.X = -av.X / 2
		case text.RightHAlignment:
			av.X = -av.X
		}

		switch txt.VAlignment {
		case text.BottomVAlignment:
			av.Y = -av.Y
		case text.MiddleVAlignment:
			av.Y = -av.Y / 2
		case text.TopVAlignment:
			av.Y = 0
			break
		}
		vec.X += av.X
		vec.Y += av.Y
	}

	rl.DrawTextEx(font, txt.String, vec, txt.Size, txt.Spacing, color2RayColor(color))
}

func color2RayColor(color color.Color) rl.Color {
	return rl.NewColor(color.R, color.G, color.B, color.A)
}

// IsScreenScaleChange returns if the current screen scale has changed
func IsScreenScaleChange() bool {
	return rl.IsWindowResized()
}

// GetFrameTime returns the time from the delta time for current frame
func GetFrameTime() float32 {
	return rl.GetFrameTime()
}

var textureHold = make(map[string]rl.Texture2D, 0)

// LoadTexture giving it file name into VRAM
func LoadTexture(fileName string) error {
	if t := rl.LoadTexture(fileName); t.ID != 0 {
		textureHold[fileName] = t
		return nil
	}
	return fmt.Errorf("error loading texture: %q", fileName)
}

// UnloadAllTextures from VRAM
func UnloadAllTextures() {
	for k, v := range textureHold {
		delete(textureHold, k)
		rl.UnloadTexture(v)
	}
}

// DrawSprite draws a sprite.Sprite in the given geometry.Position with the tint color.Color
func DrawSprite(def components.SpriteDef, sprite sprite.Sprite, pos geometry.Position, tint color.Color) error {
	if val, ok := textureHold[def.Texture]; ok {
		scale := sprite.Scale
		px := def.Origin.Size.Width / 2
		py := def.Origin.Size.Height / 2
		rc := color2RayColor(tint)
		rotation := sprite.Rotation
		sourceRec := rl.Rectangle{
			X:      def.Origin.From.X,
			Y:      def.Origin.From.Y,
			Width:  def.Origin.Size.Width,
			Height: def.Origin.Size.Height,
		}
		destRec := rl.Rectangle{
			X:      pos.X - (px * scale),
			Y:      pos.Y - (py * scale),
			Width:  def.Origin.Size.Width * scale,
			Height: def.Origin.Size.Height * scale,
		}
		origin := rl.Vector2{X: 0, Y: 0}
		rl.DrawTexturePro(val, sourceRec, destRec, origin, rotation, rc)
	} else {
		return fmt.Errorf("error drawing sprite, texture not found: %q", sprite.Name)
	}
	return nil
}

// GetMousePosition returns the current position of the mouse
func GetMousePosition() (x, y float32) {
	pos := rl.GetMousePosition()
	return pos.X, pos.Y
}
