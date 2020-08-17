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
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/options"
	"image/color"
)

var saveOpts = options.Options{}

func Init(opt options.Options) {
	saveOpts = opt
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(opt.Width, opt.Height, opt.Title)
}

func End() {
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

func DrawText(text components.Text, pos components.Pos, txColor color.Color) {
	font := rl.GetFontDefault()

	vec := rl.Vector2{
		X: float32(pos.X),
		Y: float32(pos.Y),
	}

	if text.HAlignment != components.LeftHAlignment || text.VAlignment != components.BottomVAlignment {
		av := rl.MeasureTextEx(font, text.String, float32(text.Size), 10)

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
			av.Y = 0
		case components.MiddleVAlignment:
			av.Y = -av.Y / 2
		case components.TopVAlignment:
			av.Y = -av.Y
			break
		}
		vec.X += av.X
		vec.Y += av.Y
	}

	rl.DrawTextEx(font, text.String, vec, float32(text.Size), 10, color2RayColor(txColor))
}

func color2RayColor(color color.Color) rl.Color {
	r, g, b, a := color.RGBA()
	return rl.NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
}
