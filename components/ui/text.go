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
)

// HAlignment horizontal alignment for a Text
type HAlignment int

// Horizontal Text alignments
const (
	LeftHAlignment   = HAlignment(iota) // LeftHAlignment indicates left Text HAlignment
	RightHAlignment                     // RightHAlignment indicates right Text HAlignment
	CenterHAlignment                    // RightHAlignment indicates center Text HAlignment
)

// VAlignment vertical alignment for a Text
type VAlignment int

// Vertical Text alignments
const (
	BottomVAlignment = VAlignment(iota) // BottomVAlignment indicates bottom text.VAlignment
	TopVAlignment                       // TopVAlignment indicates top text.VAlignment
	MiddleVAlignment                    // MiddleVAlignment indicates middle text.VAlignment
)

// Text is a graphical text to drawn on the screen
type Text struct {
	String     string     // String is the Text string
	Size       float32    // Size is the Text size
	Font       string     // Font is the font to use
	VAlignment VAlignment // VAlignment is the text.VAlignment
	HAlignment HAlignment // HAlignment is the text.HAlignment
}

// Type return this goecs.ComponentType
func (t Text) Type() goecs.ComponentType {
	return TYPE.Text
}
