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
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
)

// FlatButtonColors represent the colors of the button
type FlatButtonColors struct {
	Solid    color.Solid    // Solid is a color.Solid
	Gradient color.Gradient // Gradient is a color.Gradient
	Set      bool           // Set indicates is this has been set
}

// FlatButton is a UI element for displaying a button
type FlatButton struct {
	Shadow geometry.Size    // Shadow is the offset of the shadow on the ui.FlatButton
	Normal FlatButtonColors // Normal is the ui.FlatButtonColors on normal state
	Hover  FlatButtonColors // Hover is the  ui.FlatButtonColors on hover state
	Event  interface{}      // Event is the event that will be trigger when this button is click
}
