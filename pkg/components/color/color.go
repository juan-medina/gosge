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

package color

import "reflect"

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

var TYPE = reflect.TypeOf(Color{})

func New(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

func (rc Color) Alpha(alpha uint8) Color {
	return New(rc.R, rc.G, rc.B, alpha)
}

func (rc Color) Blend(other Color, scale float64) Color {
	r1 := float64(rc.R)
	g1 := float64(rc.G)
	b1 := float64(rc.B)
	a1 := float64(rc.A)

	r2 := float64(other.R)
	g2 := float64(other.G)
	b2 := float64(other.B)
	a2 := float64(other.A)

	r := r1 + ((r2 - r1) * scale)
	g := g1 + ((g2 - g1) * scale)
	b := b1 + ((b2 - b1) * scale)
	a := a1 + ((a2 - a1) * scale)

	return New(uint8(r), uint8(g), uint8(b), uint8(a))
}

//goland:noinspection GoUnusedGlobalVariable
var (
	Black      = New(0, 0, 0, 255)
	White      = New(255, 255, 255, 255)
	Magenta    = New(255, 0, 255, 255)
	LightGray  = New(200, 200, 200, 25)
	Gray       = New(130, 130, 130, 255)
	DarkGray   = New(80, 80, 80, 255)
	Yellow     = New(253, 249, 0, 255)
	Gold       = New(255, 203, 0, 255)
	Orange     = New(255, 161, 0, 255)
	Pink       = New(255, 109, 194, 255)
	Red        = New(230, 41, 55, 255)
	Maroon     = New(190, 33, 55, 255)
	Green      = New(0, 228, 48, 255)
	Lime       = New(0, 158, 47, 255)
	DarkGreen  = New(0, 117, 44, 255)
	SkyBlue    = New(102, 191, 255, 255)
	Blue       = New(0, 121, 241, 255)
	DarkBlue   = New(0, 82, 172, 255)
	Purple     = New(200, 122, 255, 255)
	Violet     = New(135, 60, 190, 255)
	DarkPurple = New(112, 31, 126, 255)
	Beige      = New(211, 176, 131, 255)
	Brown      = New(127, 106, 79, 255)
	DarkBrown  = New(76, 63, 47, 255)
	Gopher     = New(106, 215, 229, 255)
)
