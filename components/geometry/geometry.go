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

// Package geometry handle components with geometry concepts
package geometry

import (
	"github.com/juan-medina/goecs"
	"math"
	"reflect"
)

// Point represent a x/y coordinate
type Point struct {
	X float32 // The x coordinate
	Y float32 // The y coordinate
}

// Clamp a geometry.Point to min max value
func (pos *Point) Clamp(min Point, max Point) {
	if pos.X < min.X {
		pos.X = min.X
	}

	if pos.X > max.X {
		pos.X = max.X
	}

	if pos.Y > max.Y {
		pos.Y = max.Y
	}

	if pos.Y < min.Y {
		pos.Y = min.Y
	}
}

// Size represent the size of an object
type Size struct {
	Width  float32 // Width is vertical length of the size
	Height float32 // Height is the horizontal length of the size
}

// Scale is the current scale between two Size
type Scale struct {
	Point Point   // Point is the scale per point
	Min   float32 // Min is the minimum scale between X and Y
	Max   float32 // Max is the maximum scale between X and Y
}

// CalculateScale calculates the Scale between to Size
func (s Size) CalculateScale(other Size) Scale {
	sx := s.Width / other.Width
	sy := s.Height / other.Height
	return Scale{
		Min: float32(math.Min(float64(sx), float64(sy))),
		Max: float32(math.Max(float64(sx), float64(sy))),
		Point: Point{
			X: sx,
			Y: sy,
		},
	}
}

// Scale return a geometry.Size scaled by a factor
func (s Size) Scale(factor float32) Size {
	return Size{
		Width:  s.Width * factor,
		Height: s.Height * factor,
	}
}

// Rect is a rectangular area
type Rect struct {
	From Point // From is the origin of the area
	Size Size  // Size is the size of the area
}

// IsPointInRect return if a given geometry.Point is inside the geometry.Rect
func (r Rect) IsPointInRect(point Point) bool {
	return r.From.X <= point.X &&
		r.From.Y <= point.Y &&
		r.From.X+r.Size.Width >= point.X &&
		r.From.Y+r.Size.Height >= point.Y
}

func (r Rect) collidesOneDirection(other Rect) bool {
	point1 := Point{X: other.From.X, Y: other.From.Y}
	point2 := Point{X: other.From.X, Y: other.From.Y + other.Size.Height}
	point3 := Point{X: other.From.X + other.Size.Width, Y: other.From.Y}
	point4 := Point{X: other.From.X + other.Size.Width, Y: other.From.Y + other.Size.Height}

	return r.IsPointInRect(point1) || r.IsPointInRect(point2) || r.IsPointInRect(point3) || r.IsPointInRect(point4)
}

// Collides return if a given geometry.Rect collides with this geometry.Rect
func (r Rect) Collides(other Rect) bool {
	return r.collidesOneDirection(other) || other.collidesOneDirection(r)
}

type types struct {
	// Point is the reflect.Type for geometry.Point
	Point reflect.Type
}

// TYPE hold the reflect.Type for our geometry components
var TYPE = types{
	Point: reflect.TypeOf(Point{}),
}

type gets struct {
	// Point gets a geometry.Point from a goecs.Entity
	Point func(e *goecs.Entity) Point
}

// Get a geometry component
var Get = gets{
	// Point gets a geometry.Point from a goecs.Entity
	Point: func(e *goecs.Entity) Point {
		return e.Get(TYPE.Point).(Point)
	},
}
