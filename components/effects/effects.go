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
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components/color"
)

// AlternateColorState is the state for an effect
type AlternateColorState struct {
	CurrentTime float32     // CurrentTime time that this effects has been running
	Running     bool        // Running indicates if this effect is running
	From        color.Solid // From is color.Solid that we start from in the current state
	To          color.Solid // To is the color.Solid that we will end to the current state
}

// Type return this goecs.ComponentType
func (a AlternateColorState) Type() goecs.ComponentType {
	return TYPE.AlternateColorState
}

// AlternateColor effects will cycle between two colors From and To in given Time with an Optional Delay
type AlternateColor struct {
	From  color.Solid // From is color.Solid that we start from
	To    color.Solid // To is the color.Solid that we will end to
	Time  float32     // Time is how long will be get to go from From to To in seconds
	Delay float32     // Delay is how long will stay in the final To until switching again to To
}

// Type return this goecs.ComponentType
func (a AlternateColor) Type() goecs.ComponentType {
	return TYPE.AlternateColor
}

// Layer effect is use to render thins in a logical layer
type Layer struct {
	Depth float32 // Depth is the screen depth for this effects.Layer
}

// Type return this goecs.ComponentType
func (l Layer) Type() goecs.ComponentType {
	return TYPE.Layer
}

// Hide indicates that this entity should be not draw
type Hide struct{}

// Type return this goecs.ComponentType
func (h Hide) Type() goecs.ComponentType {
	return TYPE.Hide
}

type types struct {
	// AlternateColorState is the goecs.ComponentType for effects.AlternateColorState
	AlternateColorState goecs.ComponentType
	// AlternateColor is the goecs.ComponentType for effects.AlternateColor
	AlternateColor goecs.ComponentType
	// Layer is the goecs.ComponentType for effects.Layer
	Layer goecs.ComponentType
	// Hide is the goecs.ComponentType for effects.Hide
	Hide goecs.ComponentType
}

// TYPE hold the goecs.ComponentType for our effects components
var TYPE = types{
	AlternateColorState: goecs.NewComponentType(),
	AlternateColor:      goecs.NewComponentType(),
	Layer:               goecs.NewComponentType(),
	Hide:                goecs.NewComponentType(),
}

type gets struct {
	// AlternateColorState gets a AlternateColorState from a goecs.Entity
	AlternateColorState func(e *goecs.Entity) AlternateColorState
	// AlternateColor gets a AlternateColor from a goecs.Entity
	AlternateColor func(e *goecs.Entity) AlternateColor
	// Layer gets a Layer from a goecs.Entity
	Layer func(e *goecs.Entity) Layer
	// Hide gets a Hide from a goecs.Entity
	Hide func(e *goecs.Entity) Hide
}

// Get effect component
var Get = gets{
	// AlternateColorState gets a AlternateColorState from a goecs.Entity
	AlternateColorState: func(e *goecs.Entity) AlternateColorState {
		return e.Get(TYPE.AlternateColorState).(AlternateColorState)
	},
	// AlternateColor gets a AlternateColor from a goecs.Entity
	AlternateColor: func(e *goecs.Entity) AlternateColor {
		return e.Get(TYPE.AlternateColor).(AlternateColor)
	},
	// Layer gets a Layer from a goecs.Entity
	Layer: func(e *goecs.Entity) Layer {
		return e.Get(TYPE.Layer).(Layer)
	},
	// HIde gets a Hide from a goecs.Entity
	Hide: func(e *goecs.Entity) Hide {
		return e.Get(TYPE.Hide).(Hide)
	},
}
