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
	"github.com/juan-medina/gosge/pkg/components/color"
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

type AlternateColor struct {
	From    color.Color
	To      color.Color
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
