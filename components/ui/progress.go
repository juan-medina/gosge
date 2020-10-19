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

package ui

import (
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/geometry"
)

// ProgressBar is a component representing a progress bar
type ProgressBar struct {
	Min     float32       // Min is the minimum value for the ui.ProgressBar
	Max     float32       // Max is the maximum value for the ui.ProgressBar
	Current float32       // Current is the current value for the ui.ProgressBar
	Shadow  geometry.Size // Shadow is the offset of the shadow on the ui.ProgressBar
	Sound   string        // Sound is the click sound
	Volume  float32       // Volume is the volume for click Sound
	Event   interface{}   // Event is the event that will be trigger when this bar is click
}

// Type return this goecs.ComponentType
func (p ProgressBar) Type() goecs.ComponentType {
	return TYPE.ProgressBar
}

// ProgressBarColor represent the colors of the ui.ProgressBar
type ProgressBarColor struct {
	Solid    color.Solid    // Solid is a color.Solid
	Gradient color.Gradient // Gradient is a color.Gradient
	Empty    color.Solid    // Empty is the empty bar color
	Border   color.Solid    // Border is a color.Solid
}

// Type return this goecs.ComponentType
func (p ProgressBarColor) Type() goecs.ComponentType {
	return TYPE.ProgressBarColor
}

// ProgressBarHoverColor is the hover colors for a ui.ProgressBar
type ProgressBarHoverColor struct {
	Normal   ProgressBarColor // Normal is the ui.ProgressBarColor on normal state
	Hover    ProgressBarColor // Hover is the  ui.ProgressBarColor on hover state
	Clicked  ProgressBarColor // Clicked is the  ui.ProgressBarColor on clicked state
	Disabled ProgressBarColor // Disabled is the  ui.ProgressBarColor on disabled state
}

// Type return this goecs.ComponentType
func (p ProgressBarHoverColor) Type() goecs.ComponentType {
	return TYPE.ProgressBarHoverColor
}
