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

package tiled

import (
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"reflect"
)

// Map is a tiled.Map
type Map struct {
	Name string // Name of our tiled.Map
}

// MapState is the state for tiled.Map
type MapState struct {
	Pos geometry.Point
}

type types struct {
	// Map is the reflect.Type for tiled.Map
	Map reflect.Type
	// MapState is the reflect.Type for tiled.MapState
	MapState reflect.Type
}

// TYPE hold the reflect.Type for our tiled components
var TYPE = types{
	Map:      reflect.TypeOf(Map{}),
	MapState: reflect.TypeOf(MapState{}),
}

type gets struct {
	// Map gets a tiled.Map from a entity.Entity
	Map func(e *entity.Entity) Map
	// MapState gets a tiled.MapState from a entity.Entity
	MapState func(e *entity.Entity) MapState
}

// Get a geometry component
var Get = gets{
	// Map gets a tiled.Map from a entity.Entity
	Map: func(e *entity.Entity) Map {
		return e.Get(TYPE.Map).(Map)
	},
	// MapState gets a tiled.MapState from a entity.Entity
	MapState: func(e *entity.Entity) MapState {
		return e.Get(TYPE.MapState).(MapState)
	},
}
