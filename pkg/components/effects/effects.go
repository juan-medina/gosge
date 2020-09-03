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

// Package effects include different components for adding effects
package effects

import (
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/gosge/pkg/components/color"
	"reflect"
)

// AlternateColorState is the state for an effect
type AlternateColorState struct {
	CurrentTime float32     // CurrentTime time that this effects has been running
	Running     bool        // Running indicates if this effect is running
	From        color.Solid // From is color.Solid that we start from in the current state
	To          color.Solid // To is the color.Solid that we will end to the current state
}

// AlternateColor effects will cycle between two colors From and To in given Time with an Optional Delay
type AlternateColor struct {
	From  color.Solid // From is color.Solid that we start from
	To    color.Solid // To is the color.Solid that we will end to
	Time  float32     // Time is how long will be get to go from From to To in seconds
	Delay float32     // Delay is how long will stay in the final To until switching again to To
}

// Layer effect is use to render thins in a logical layer
type Layer struct {
	Depth float32 // Depth is the screen depth for this effects.Layer
}

type types struct {
	// AlternateColorState is the reflect.Type for effects.AlternateColorState
	AlternateColorState reflect.Type
	// AlternateColor is the reflect.Type for effects.AlternateColor
	AlternateColor reflect.Type
	// Layer is the reflect.Type for effects.Layer
	Layer reflect.Type
}

// TYPE hold the reflect.Type for our effects components
var TYPE = types{
	AlternateColorState: reflect.TypeOf(AlternateColorState{}),
	AlternateColor:      reflect.TypeOf(AlternateColor{}),
	Layer:               reflect.TypeOf(Layer{}),
}

type gets struct {
	// AlternateColorState gets a AlternateColorState from a entity.Entity
	AlternateColorState func(e *entity.Entity) AlternateColorState
	// AlternateColor gets a AlternateColor from a entity.Entity
	AlternateColor func(e *entity.Entity) AlternateColor
	// Layer gets a Layer from a entity.Entity
	Layer func(e *entity.Entity) Layer
}

// Get effect component
var Get = gets{
	// AlternateColorState gets a AlternateColorState from a entity.Entity
	AlternateColorState: func(e *entity.Entity) AlternateColorState {
		return e.Get(TYPE.AlternateColorState).(AlternateColorState)
	},
	// AlternateColor gets a AlternateColor from a entity.Entity
	AlternateColor: func(e *entity.Entity) AlternateColor {
		return e.Get(TYPE.AlternateColor).(AlternateColor)
	},
	// Layer gets a Layer from a entity.Entity
	Layer: func(e *entity.Entity) Layer {
		return e.Get(TYPE.Layer).(Layer)
	},
}
