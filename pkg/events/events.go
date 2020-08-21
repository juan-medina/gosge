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

// Package events contain the events for our engine
package events

import "github.com/juan-medina/gosge/pkg/components/geometry"

// GameCloseEvent is an event that indicates that game need to close
type GameCloseEvent struct{}

// ScreenSizeChangeEvent is an event that indicates that the screen size has change
type ScreenSizeChangeEvent struct {
	// Current is the current screen size
	Current geometry.Size
	// Original is the original, from options.Options, screen size
	Original geometry.Size
	// Scale is the current scale between Current and Original
	Scale struct {
		Point geometry.Point //Point is the scale per point
		Min   float32        // Min is the minimum scale between X and Y
		Max   float32        // Max is the maximum scale between X and Y
	}
}

// MouseMoveEvent is an event that indicates that the mouse is moving
type MouseMoveEvent struct {
	geometry.Position
}
