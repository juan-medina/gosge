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

//Package shapes contains various drawable shapes
package shapes

import (
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"reflect"
)

//Box is a rectangular shape that we could draw in a geometry.Point with a color.Solid or color.Gradient
type Box struct {
	Size  geometry.Size // The box size
	Scale float32       // The box scale
}

// Contains return if a box at a geometry.Point contains a point
func (b Box) Contains(at geometry.Point, point geometry.Point) bool {
	return geometry.Rect{
		From: at,
		Size: geometry.Size{
			Width:  b.Size.Width * b.Scale,
			Height: b.Size.Height * b.Scale,
		},
	}.IsPointInRect(point)
}

type types struct {
	// Box is the reflect.Type for shapes.Box
	Box reflect.Type
}

// TYPE hold the reflect.Type for our shapes components
var TYPE = types{
	Box: reflect.TypeOf(Box{}),
}

type gets struct {
	// Box gets a shapes.Box from a goecs.Entity
	Box func(e *goecs.Entity) Box
}

// Get a geometry component
var Get = gets{
	// Box gets a shapes.Box from a goecs.Entity
	Box: func(e *goecs.Entity) Box {
		return e.Get(TYPE.Box).(Box)
	},
}
