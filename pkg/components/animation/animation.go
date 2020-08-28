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

// Package animation are the components for running and sprite.Sprite animations
package animation

import (
	"github.com/juan-medina/goecs/pkg/entity"
	"reflect"
)

// Sequence represent a set of frames that will be render with a delay
type Sequence struct {
	Sheet    string  // Sheet is the sprite sheet where the animation sprites are
	Base     string  // Base is the base name for each frame. ex : Idle_%d.png
	Rotation float32 // Rotation for this AnimationSequence
	Scale    float32 // Scale for this AnimationSequence
	Frames   int32   // Frames are the number of frame in this animation
	Delay    float32 // Delay number of seconds to wait in each frame
}

// State allow to easily switch animations
type State struct {
	Current string  // Current is the animation that is running
	Speed   float32 // Speed is the current animation speed
	Time    float32 // Time is the time in this frame
	Frame   int32   // Frame is the current frame number
}

// Animation allow to easily switch animations
type Animation struct {
	Sequences map[string]Sequence // Sequences is the different AnimationSequence for this Animation
	Current   string              // Current is the animation that be like to run
	Speed     float32             // Speed is a multiplier of th speed of the current animation
	FlipX     bool                // FlipX indicates if the Animation is flipped in the X-Assis
	FlipY     bool                // FlipY indicates if the Animation is flipped in the Y-Assis
}

type types struct {
	// Animation is the reflect.Type for animation.Animation
	Animation reflect.Type
	// State is the reflect.Type for animation.State
	State reflect.Type
	// Sequence is the reflect.Type for animation.Sequence
	Sequence reflect.Type
}

// TYPE hold the reflect.Type for our animation components
var TYPE = types{
	Animation: reflect.TypeOf(Animation{}),
	State:     reflect.TypeOf(State{}),
	Sequence:  reflect.TypeOf(Sequence{}),
}

type gets struct {
	// AlternateColor gets a AlternateColor from a entity.Entity
	Animation func(e *entity.Entity) Animation
	// Layer gets a Layer from a entity.Entity
	State func(e *entity.Entity) State
	// Layer gets a Layer from a entity.Entity
	Sequence func(e *entity.Entity) Sequence
}

// Get animation component
var Get = gets{
	// Animation gets a animation.Animation from a entity.Entity
	Animation: func(e *entity.Entity) Animation {
		return e.Get(TYPE.Animation).(Animation)
	},
	// State gets a animation.State from a entity.Entity
	State: func(e *entity.Entity) State {
		return e.Get(TYPE.State).(State)
	},
	// Sequence gets a animation.Sequence from a entity.Entity
	Sequence: func(e *entity.Entity) Sequence {
		return e.Get(TYPE.Sequence).(Sequence)
	},
}
