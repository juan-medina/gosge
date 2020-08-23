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

package ray

import (
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/juan-medina/gosge/pkg/options"
	"os"
)

// Init the rendering device
func (rr *RenderImpl) Init(opt options.Options) {
	rr.saveOpts = opt

	w := rl.GetMonitorWidth(opt.Monitor)
	h := rl.GetMonitorHeight(opt.Monitor)
	rl.InitWindow(int32(w), int32(h), opt.Title)

	if opt.Icon != "" {
		if file, err := os.Open(opt.Icon); err == nil {
			_ = file.Close()
			img := rl.LoadImage(opt.Icon)
			rl.SetWindowIcon(*img)
			rl.UnloadImage(img)
		}
	}

	rl.ToggleFullscreen()
}

// End the rendering device
func (rr RenderImpl) End() {
	rr.UnloadAllTextures()
	rl.CloseWindow()
}

// BeginFrame for rendering
func (rr RenderImpl) BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(rr.color2RayColor(rr.saveOpts.BackGround))
}

// EndFrame for rendering
func (rr RenderImpl) EndFrame() {
	rl.EndDrawing()
}

// ShouldClose returns if th engine should close
func (rr RenderImpl) ShouldClose() bool {
	return rl.WindowShouldClose()
}
