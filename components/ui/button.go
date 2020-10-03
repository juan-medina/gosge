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
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/geometry"
)

// ButtonColor represent the colors of the button
type ButtonColor struct {
	Solid    color.Solid    // Solid is a color.Solid
	Gradient color.Gradient // Gradient is a color.Gradient
	Border   color.Solid    // Border is the border color.Solid
	Text     color.Solid    // Text is the text color.Solid
}

// ButtonHoverColors is the hover colors for a FlatButton
type ButtonHoverColors struct {
	Normal  ButtonColor // Normal is the ui.ButtonColor on normal state
	Hover   ButtonColor // Hover is the  ui.ButtonColor on hover state
	Clicked ButtonColor // Clicked is the  ui.ButtonColor on clicked state
}

// FlatButton is a UI element for displaying a button
type FlatButton struct {
	Shadow geometry.Size // Shadow is the offset of the shadow on the ui.FlatButton
	State  ControlState  // State is this ui.FlatButton ui.ControlState
	Sound  string        // Sound is the click sound
	Volume float32       // Volume is the volume for click Sound
	Event  interface{}   // Event is the event that will be trigger when this button is click
}

// SpriteButton is a UI element for displaying a image base button
type SpriteButton struct {
	Sheet   string       // Sheet is the Sprite sheet
	Normal  string       // Normal is the sprite on normal state
	Hover   string       // Hover is the sprite on hover state
	Clicked string       // Clicked is the sprite on clicked state
	State   ControlState // State is the ui.SpriteButton ui.ControlState
	Scale   float32      // Scale is the Sprite scale
	Sound   string       // Sound is the click sound
	Volume  float32      // Volume is the volume for click Sound
	Event   interface{}  // Event is the event that will be trigger when this button is click
}
