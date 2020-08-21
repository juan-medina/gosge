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

// EffectState is the state for an effect
type EffectState int

// EffectState states
const (
	NoState      = EffectState(iota) // NoState represents a not set EffectState
	StateStopped                     // StateStopped represent an EffectState that is stopped
	StateRunning                     // StateRunning represent an EffectState that is running
)

// AlternateColor effects will cycle between two colors From and To in given Time with an Optional Delay
type AlternateColor struct {
	From    color.Color // From is color.Color that we start from
	To      color.Color // To is the color.Color that we will end to
	Time    float32     // Time is how long will be get to go from From to To in seconds
	Delay   float32     // Delay is how long will stay in the final To until switching again to To
	Current float32     // Current time that the effects is running
	State   EffectState // State is the current EffectState
}

type types struct {
	// AlternateColor is the reflect.Type for effects.AlternateColor
	AlternateColor reflect.Type
}

// TYPE hold the reflect.Type for our effects components
var TYPE = types{
	AlternateColor: reflect.TypeOf(AlternateColor{}),
}

type gets struct {
	// AlternateColor is the reflect.Type for effects.AlternateColor
	AlternateColor func(e *entity.Entity) AlternateColor
}

// Get hold the reflect.Type for our effects components
var Get = gets{
	// AlternateColor gets a AlternateColor from a entity.Entity
	AlternateColor: func(e *entity.Entity) AlternateColor {
		return e.Get(TYPE.AlternateColor).(AlternateColor)
	},
}
