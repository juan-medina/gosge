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

import (
	"github.com/juan-medina/goecs"
)

// Solid represents a RGBA color
type Solid struct {
	R uint8 // R is the green Color component
	G uint8 // G is the green Color component
	B uint8 // B is the blue Color component
	A uint8 // A is the alpha Color component
}

// Type return this goecs.ComponentType
func (rc Solid) Type() goecs.ComponentType {
	return TYPE.Solid
}

// Alpha returns a new Color modifying the A component
func (rc Solid) Alpha(alpha uint8) Solid {
	return Solid{R: rc.R, G: rc.G, B: rc.B, A: alpha}
}

// Blend the Color with a given Color using an scale
func (rc Solid) Blend(other Solid, scale float32) Solid {
	r1 := float32(rc.R)
	g1 := float32(rc.G)
	b1 := float32(rc.B)
	a1 := float32(rc.A)

	r2 := float32(other.R)
	g2 := float32(other.G)
	b2 := float32(other.B)
	a2 := float32(other.A)

	r := r1 + ((r2 - r1) * scale)
	g := g1 + ((g2 - g1) * scale)
	b := b1 + ((b2 - b1) * scale)
	a := a1 + ((a2 - a1) * scale)

	return Solid{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

// Inverse the color.Solid
func (rc Solid) Inverse() Solid {
	r1 := int32(rc.R)
	g1 := int32(rc.G)
	b1 := int32(rc.B)
	a1 := int32(rc.A)

	r := 255 - r1
	g := 255 - g1
	b := 255 - b1
	a := a1

	return Solid{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

// GrayScale converts a color.Solid to gray scale
func (rc Solid) GrayScale() Solid {
	r := float64(rc.R) * 0.2989
	g := float64(rc.G) * 0.5870
	b := float64(rc.B) * 0.1140
	gray := uint8(r + g + b)

	return Solid{R: gray, G: gray, B: gray, A: rc.A}
}

// Equals returns if two color.Solid are equal
func (rc Solid) Equals(other Solid) bool {
	return rc.R == other.R && rc.G == other.G && rc.B == other.B && rc.A == other.A
}

// GradientDirection is the direction of color.Gradient
type GradientDirection int

// Directions for a color.Gradient
//goland:noinspection ALL
const (
	GradientHorizontal = GradientDirection(iota) // GradientHorizontal is a horizontal color.Gradient
	GradientVertical                             // GradientHorizontal is a vertical color.Gradient
)

// Gradient represents a gradient color
type Gradient struct {
	From      Solid             // From what color.Solid the gradient stars
	To        Solid             // To what color.Solid the gradient ends
	Direction GradientDirection // Direction is GradientDirection for this color.Gradient
}

// Type return this goecs.ComponentType
func (g Gradient) Type() goecs.ComponentType {
	return TYPE.Gradient
}

//goland:noinspection GoUnusedGlobalVariable
var (
	Black      = Solid{A: 255}                         // Black Color
	White      = Solid{R: 255, G: 255, B: 255, A: 255} // White Color
	Magenta    = Solid{R: 255, B: 255, A: 255}         // Magenta Color
	LightGray  = Solid{R: 200, G: 200, B: 200, A: 255} // LightGray is a Light Gray Color
	Gray       = Solid{R: 130, G: 130, B: 130, A: 255} // Gray Color
	DarkGray   = Solid{R: 80, G: 80, B: 80, A: 255}    // DarkGray Dark Gray
	Yellow     = Solid{R: 253, G: 249, A: 255}         // Yellow Color
	Gold       = Solid{R: 255, G: 203, A: 255}         // Gold Color
	Orange     = Solid{R: 255, G: 161, A: 255}         // Orange Color
	Pink       = Solid{R: 255, G: 109, B: 194, A: 255} // Pink Color
	Red        = Solid{R: 230, G: 41, B: 55, A: 255}   // Red Color
	Maroon     = Solid{R: 190, G: 33, B: 55, A: 255}   // Marron Color
	Green      = Solid{G: 228, B: 48, A: 255}          // Green Color
	Lime       = Solid{G: 158, B: 47, A: 255}          // Lime Color
	DarkGreen  = Solid{G: 117, B: 44, A: 255}          // DarkGreen is a Dark Green Color
	SkyBlue    = Solid{R: 102, G: 191, B: 255, A: 255} // SkyBlue Color
	Blue       = Solid{G: 121, B: 241, A: 255}         // Blue Color
	DarkBlue   = Solid{G: 82, B: 172, A: 255}          // DarkBlue is a Dark Blue Color
	Purple     = Solid{R: 200, G: 122, B: 255, A: 255} // Purple Color
	Violet     = Solid{R: 135, G: 60, B: 190, A: 255}  // Violet Color
	DarkPurple = Solid{R: 112, G: 31, B: 126, A: 255}  // DarkPurple is a Dark Purple Color
	Beige      = Solid{R: 211, G: 176, B: 131, A: 255} // Beige Color
	Brown      = Solid{R: 127, G: 106, B: 79, A: 255}  // Brown Color
	DarkBrown  = Solid{R: 76, G: 63, B: 47, A: 255}    // DarkBrown Color
	Gopher     = Solid{R: 106, G: 215, B: 229, A: 255} // Gopher Color
)

type types struct {
	// Solid is the goecs.ComponentType for color.Solid
	Solid goecs.ComponentType
	// Gradient is the goecs.ComponentType for color.Gradient
	Gradient goecs.ComponentType
}

// TYPE hold the goecs.ComponentType for our color components
var TYPE = types{
	Solid:    goecs.NewComponentType(),
	Gradient: goecs.NewComponentType(),
}

type gets struct {
	// Solid get a color.Solid from a goecs.Entity
	Solid func(e *goecs.Entity) Solid
	// Gradient get a color.Gradient from a goecs.Entity
	Gradient func(e *goecs.Entity) Gradient
}

// Get a color component
var Get = gets{
	// Solid get a color.Solid from a goecs.Entity
	Solid: func(e *goecs.Entity) Solid {
		return e.Get(TYPE.Solid).(Solid)
	},
	// Gradient get a color.Gradient from a goecs.Entity
	Gradient: func(e *goecs.Entity) Gradient {
		return e.Get(TYPE.Gradient).(Gradient)
	},
}
