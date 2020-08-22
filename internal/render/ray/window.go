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
	"github.com/juan-medina/gosge/pkg/components/geometry"
)

// GetScreenSize get the current screen size
func (rr RenderImpl) GetScreenSize() geometry.Size {
	return geometry.Size{Width: float32(rl.GetScreenWidth()), Height: float32(rl.GetScreenHeight())}
}

// GetMousePosition returns the current position of the mouse
func (rr RenderImpl) GetMousePosition() geometry.Position {
	pos := rl.GetMousePosition()
	return geometry.Position{X: pos.X, Y: pos.Y}
}

// IsScreenScaleChange returns if the current screen scale has changed
func (rr RenderImpl) IsScreenScaleChange() bool {
	return rl.IsWindowResized()
}
