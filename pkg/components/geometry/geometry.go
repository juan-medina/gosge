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
	"github.com/juan-medina/goecs/pkg/entity"
	"reflect"
)

// Point represent a x/y coordinate
type Point struct {
	X float32 // The x coordinate
	Y float32 // The y coordinate
}

// Size represent the size of an object
type Size struct {
	Width  float32 // Width is vertical length of the size
	Height float32 // Height is the horizontal length of the size
}

// Rect is a rectangular area
type Rect struct {
	From Point // From is the origin of the area
	Size Size  // Size is the size of the area
}

// Position represent an X and Y screen position
type Position Point

type types struct {
	// Position is the reflect.Type for geometry.Position
	Position reflect.Type
}

// TYPE hold the reflect.Type for our geometry components
var TYPE = types{
	Position: reflect.TypeOf(Position{}),
}

type gets struct {
	// Position is the reflect.Type for geometry.Position
	Position func(e *entity.Entity) Position
}

// Get hold the reflect.Type for our geometry components
var Get = gets{
	// Position gets a geometry.Position from a entity.Entity
	Position: func(e *entity.Entity) Position {
		return e.Get(TYPE.Position).(Position)
	},
}
