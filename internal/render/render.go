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

package render

import (
	"errors"
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/options"
)

var saveOpts = options.Options{}

func Init(opt options.Options) {
	saveOpts = opt
	//rl.SetTraceLog(rl.LogNone)
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(opt.Width), int32(opt.Height), opt.Title)
}

func End() {
	UnloadAllTextures()
	rl.CloseWindow()
}

func BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(color2RayColor(saveOpts.ClearColor))
}

func EndFrame() {
	rl.EndDrawing()
}

func ShouldClose() bool {
	return rl.WindowShouldClose()
}

func GetScreenSize() (width int, height int) {
	return rl.GetScreenWidth(), rl.GetScreenHeight()
}

func DrawText(text components.UiText, pos components.Pos, txColor components.RGBAColor) {
	font := rl.GetFontDefault()

	vec := rl.Vector2{
		X: float32(pos.X),
		Y: float32(pos.Y),
	}

	if text.HAlignment != components.LeftHAlignment || text.VAlignment != components.BottomVAlignment {
		av := rl.MeasureTextEx(font, text.String, float32(text.Size), float32(text.Spacing))

		switch text.HAlignment {
		case components.LeftHAlignment:
			av.X = 0
		case components.CenterHAlignment:
			av.X = -av.X / 2
		case components.RightHAlignment:
			av.X = -av.X
		}

		switch text.VAlignment {
		case components.BottomVAlignment:
			av.Y = -av.Y
		case components.MiddleVAlignment:
			av.Y = -av.Y / 2
		case components.TopVAlignment:
			av.Y = 0
			break
		}
		vec.X += av.X
		vec.Y += av.Y
	}

	rl.DrawTextEx(font, text.String, vec, float32(text.Size), float32(text.Spacing), color2RayColor(txColor))
}

func color2RayColor(color components.RGBAColor) rl.Color {
	return rl.NewColor(color.R, color.G, color.B, color.A)
}

func IsScreenScaleChange() bool {
	return rl.IsWindowResized()
}

func GetFrameTime() float64 {
	return float64(rl.GetFrameTime())
}

var textureHold = make(map[string]rl.Texture2D, 0)

func LoadTexture(fileName string) error {
	if t := rl.LoadTexture(fileName); t.ID != 0 {
		textureHold[fileName] = t
		return nil
	}
	return errors.New(fmt.Sprintf("error loading texture: %q", fileName))
}

func UnloadAllTextures() {
	for k, v := range textureHold {
		delete(textureHold, k)
		rl.UnloadTexture(v)
	}
}

func DrawSprite(sprite components.Sprite, pos components.Pos, tint components.RGBAColor) error {
	if val, ok := textureHold[sprite.FileName]; ok {
		scale := float32(sprite.Scale)
		px := float32(val.Width) / 2
		py := float32(val.Height) / 2
		vec := rl.Vector2{
			X: float32(pos.X) - (px * scale),
			Y: float32(pos.Y) - (py * scale),
		}
		color := color2RayColor(tint)
		rotation := float32(sprite.Rotation)
		rl.DrawTextureEx(val, vec, rotation, scale, color)
	} else {
		return errors.New(fmt.Sprintf("error drawing sprite, texture not found: %q", sprite.FileName))
	}

	return nil
}
