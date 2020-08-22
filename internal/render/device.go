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
	rl "github.com/gen2brain/raylib-go/raylib"
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
