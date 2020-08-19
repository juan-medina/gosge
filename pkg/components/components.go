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

package components

import (
	"reflect"
)

type HAlignment int

const (
	LeftHAlignment = HAlignment(iota)
	RightHAlignment
	CenterHAlignment
)

type VAlignment int

const (
	BottomVAlignment = VAlignment(iota)
	TopVAlignment
	MiddleVAlignment
)

type UiText struct {
	String     string
	Size       float64
	Spacing    float64
	VAlignment VAlignment
	HAlignment HAlignment
}

var UiTextType = reflect.TypeOf(UiText{})

type Pos struct {
	X float64
	Y float64
}

var PosType = reflect.TypeOf(Pos{})

type RGBAColor struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

var RGBAColorType = reflect.TypeOf(RGBAColor{})

func (rc RGBAColor) Blend(other RGBAColor, scale float64) RGBAColor {
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

	return NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
}

func NewColor(r, g, b, a uint8) RGBAColor {
	return RGBAColor{R: r, G: g, B: b, A: a}
}

type AlternateColor struct {
	From    RGBAColor
	To      RGBAColor
	Time    float64
	Current float64
}

var AlternateColorType = reflect.TypeOf(AlternateColor{})
