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
	"github.com/juan-medina/gosge/components/geometry"
)

//Box is a rectangular outline that we could draw in a geometry.Point with a color.Solid
type Box struct {
	Size      geometry.Size // The box size
	Scale     float32       // The box scale
	Thickness int32         // Thickness of the line
}

// Type return this goecs.ComponentType
func (b Box) Type() goecs.ComponentType {
	return TYPE.Box
}

//SolidBox is a rectangular shape that we could draw in a geometry.Point with a color.Solid or color.Gradient
type SolidBox struct {
	Size  geometry.Size // The box size
	Scale float32       // The box scale
}

// Type return this goecs.ComponentType
func (b SolidBox) Type() goecs.ComponentType {
	return TYPE.SolidBox
}

// Contains return if a box at a geometry.Point contains a point
func (b SolidBox) Contains(at geometry.Point, point geometry.Point) bool {
	return geometry.Rect{
		From: at,
		Size: geometry.Size{
			Width:  b.Size.Width * b.Scale,
			Height: b.Size.Height * b.Scale,
		},
	}.IsPointInRect(point)
}

// GetReactAt return return a geometry.Rect at a given point
func (b Box) GetReactAt(at geometry.Point) geometry.Rect {
	return geometry.Rect{
		From: at,
		Size: geometry.Size{
			Width:  b.Size.Width * b.Scale,
			Height: b.Size.Height * b.Scale,
		},
	}
}

// Contains return if a box at a geometry.Point contains a point
func (b Box) Contains(at geometry.Point, point geometry.Point) bool {
	return b.GetReactAt(at).IsPointInRect(point)
}

// Line is a component that represent a line
type Line struct {
	To        geometry.Point // To where the line goes
	Thickness float32        // Thickness of the line
}

// Type return this goecs.ComponentType
func (l Line) Type() goecs.ComponentType {
	return TYPE.Line
}

type types struct {
	// Box is the goecs.ComponentType for shapes.Box
	Box goecs.ComponentType
	// SolidBox is the goecs.ComponentType for shapes.SolidBox
	SolidBox goecs.ComponentType
	// Line is the goecs.ComponentType for shapes.Line
	Line goecs.ComponentType
}

// TYPE hold the goecs.ComponentType for our shapes components
var TYPE = types{
	Box:      goecs.NewComponentType(),
	SolidBox: goecs.NewComponentType(),
	Line:     goecs.NewComponentType(),
}

type gets struct {
	// Box gets a shapes.Box from a goecs.Entity
	Box func(e *goecs.Entity) Box
	// SolidBox gets a shapes.SolidBox from a goecs.Entity
	SolidBox func(e *goecs.Entity) SolidBox
	// Line gets a shapes.Line from a goecs.Entity
	Line func(e *goecs.Entity) Line
}

// Get a geometry component
var Get = gets{
	// Box gets a shapes.Box from a goecs.Entity
	Box: func(e *goecs.Entity) Box {
		return e.Get(TYPE.Box).(Box)
	},
	// SolidBox gets a shapes.SolidBox from a goecs.Entity
	SolidBox: func(e *goecs.Entity) SolidBox {
		return e.Get(TYPE.SolidBox).(SolidBox)
	},
	// Line gets a shapes.Line from a goecs.Entity
	Line: func(e *goecs.Entity) Line {
		return e.Get(TYPE.Line).(Line)
	},
}
