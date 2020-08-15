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
	"image/color"
	"reflect"
)

type HAlignment int

const (
	LeftHAlignment   = HAlignment(iota)
	RightHAlignment  = HAlignment(iota)
	CenterHAlignment = HAlignment(iota)
)

type VAlignment int

const (
	BottomVAlignment = VAlignment(iota)
	TopVAlignment    = VAlignment(iota)
	MiddleVAlignment = VAlignment(iota)
)

type Text struct {
	String     string
	Size       int
	VAlignment VAlignment
	HAlignment HAlignment
}

var TextType = reflect.TypeOf(Text{})

type Pos struct {
	X float64
	Y float64
}

var PosType = reflect.TypeOf(Pos{})

type GameSettings struct {
	Width  int
	Height int
	Title  string
}

var GameSettingsType = reflect.TypeOf(GameSettings{})

var ColorType = reflect.TypeOf(color.RGBA{})
