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

func (rc RGBAColor) Alpha(alpha uint8) RGBAColor {
	return NewColor(rc.R, rc.G, rc.B, alpha)
}

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

//goland:noinspection GoUnusedGlobalVariable
var (
	BlackColor      = NewColor(0, 0, 0, 255)
	WhiteColor      = NewColor(255, 255, 255, 255)
	MagentaColor    = NewColor(255, 0, 255, 255)
	LightGrayColor  = NewColor(200, 200, 200, 25)
	GrayColor       = NewColor(130, 130, 130, 255)
	DarkGrayColor   = NewColor(80, 80, 80, 255)
	YellowColor     = NewColor(253, 249, 0, 255)
	GoldColor       = NewColor(255, 203, 0, 255)
	OrangeColor     = NewColor(255, 161, 0, 255)
	PinkColor       = NewColor(255, 109, 194, 255)
	RedColor        = NewColor(230, 41, 55, 255)
	MaroonColor     = NewColor(190, 33, 55, 255)
	GreenColor      = NewColor(0, 228, 48, 255)
	LimeColor       = NewColor(0, 158, 47, 255)
	DarkGreenColor  = NewColor(0, 117, 44, 255)
	SkyBlueColor    = NewColor(102, 191, 255, 255)
	BlueColor       = NewColor(0, 121, 241, 255)
	DarkBlueColor   = NewColor(0, 82, 172, 255)
	PurpleColor     = NewColor(200, 122, 255, 255)
	VioletColor     = NewColor(135, 60, 190, 255)
	DarkPurpleColor = NewColor(112, 31, 126, 255)
	BeigeColor      = NewColor(211, 176, 131, 255)
	BrownColor      = NewColor(127, 106, 79, 255)
	DarkBrownColor  = NewColor(76, 63, 47, 255)
)

type AlternateColor struct {
	From    RGBAColor
	To      RGBAColor
	Time    float64
	Current float64
}

var AlternateColorType = reflect.TypeOf(AlternateColor{})

type Sprite struct {
	FileName string
	Rotation float64
	Scale    float64
}

var SpriteType = reflect.TypeOf(Sprite{})
