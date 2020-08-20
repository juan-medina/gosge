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

// Package color handles Color components
package color

import "reflect"

// Color represents a RGBA color
type Color struct {
	R uint8 // R is the green Color component
	G uint8 // G is the green Color component
	B uint8 // B is the blue Color component
	A uint8 // A is the alpha Color component
}

// TYPE is the reflect.type for Color
var TYPE = reflect.TypeOf(Color{})

// Alpha returns a new Color modifying the A component
func (rc Color) Alpha(alpha uint8) Color {
	return Color{R: rc.R, G: rc.G, B: rc.B, A: alpha}
}

// Blend the Color with a given Color using an scale
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

	return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

//goland:noinspection GoUnusedGlobalVariable
var (
	Black      = Color{A: 255}                         // Black Color
	White      = Color{R: 255, G: 255, B: 255, A: 255} // White Color
	Magenta    = Color{R: 255, B: 255, A: 255}         // Magenta Color
	LightGray  = Color{R: 200, G: 200, B: 200, A: 25}  // LightGray is a Light Gray Color
	Gray       = Color{R: 130, G: 130, B: 130, A: 255} // Gray Color
	DarkGray   = Color{R: 80, G: 80, B: 80, A: 255}    // DarkGray Dark Gray
	Yellow     = Color{R: 253, G: 249, A: 255}         // Yellow Color
	Gold       = Color{R: 255, G: 203, A: 255}         // Gold Color
	Orange     = Color{R: 255, G: 161, A: 255}         // Orange Color
	Pink       = Color{R: 255, G: 109, B: 194, A: 255} // Pink Color
	Red        = Color{R: 230, G: 41, B: 55, A: 255}   // Red Color
	Maroon     = Color{R: 190, G: 33, B: 55, A: 255}   // Marron Color
	Green      = Color{G: 228, B: 48, A: 255}          // Green Color
	Lime       = Color{G: 158, B: 47, A: 255}          // Lime Color
	DarkGreen  = Color{G: 117, B: 44, A: 255}          // DarkGreen is a Dark Green Color
	SkyBlue    = Color{R: 102, G: 191, B: 255, A: 255} // SkyBlue Color
	Blue       = Color{G: 121, B: 241, A: 255}         // Blue Color
	DarkBlue   = Color{G: 82, B: 172, A: 255}          // DarkBlue is a Dark Blue Color
	Purple     = Color{R: 200, G: 122, B: 255, A: 255} // Purple Color
	Violet     = Color{R: 135, G: 60, B: 190, A: 255}  // Violet Color
	DarkPurple = Color{R: 112, G: 31, B: 126, A: 255}  // DarkPurple is a Dark Purple Color
	Beige      = Color{R: 211, G: 176, B: 131, A: 255} // Beige Color
	Brown      = Color{R: 127, G: 106, B: 79, A: 255}  // Brown Color
	DarkBrown  = Color{R: 76, G: 63, B: 47, A: 255}    // DarkBrown Color
	Gopher     = Color{R: 106, G: 215, B: 229, A: 255} // Gopher Color
)