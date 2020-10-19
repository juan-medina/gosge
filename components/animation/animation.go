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
	"github.com/juan-medina/goecs"
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

// Type return this goecs.ComponentType
func (s Sequence) Type() goecs.ComponentType {
	return TYPE.Sequence
}

// State allow to easily switch animations
type State struct {
	Current string  // Current is the animation that is running
	Speed   float32 // Speed is the current animation speed
	Time    float32 // Time is the time in this frame
	Frame   int32   // Frame is the current frame number
}

// Type return this goecs.ComponentType
func (s State) Type() goecs.ComponentType {
	return TYPE.State
}

// Animation allow to easily switch animations
type Animation struct {
	Sequences map[string]Sequence // Sequences is the different AnimationSequence for this Animation
	Current   string              // Current is the animation that be like to run
	Speed     float32             // Speed is a multiplier of th speed of the current animation
	FlipX     bool                // FlipX indicates if the Animation is flipped in the X-Assis
	FlipY     bool                // FlipY indicates if the Animation is flipped in the Y-Assis
}

// Type return this goecs.ComponentType
func (a Animation) Type() goecs.ComponentType {
	return TYPE.Animation
}

type types struct {
	// Animation is the goecs.ComponentType for animation.Animation
	Animation goecs.ComponentType
	// State is the goecs.ComponentType for animation.State
	State goecs.ComponentType
	// Sequence is the goecs.ComponentType for animation.Sequence
	Sequence goecs.ComponentType
}

// TYPE hold the goecs.ComponentType for our animation components
var TYPE = types{
	Animation: goecs.NewComponentType(),
	State:     goecs.NewComponentType(),
	Sequence:  goecs.NewComponentType(),
}

type gets struct {
	// AlternateColor gets a AlternateColor from a goecs.Entity
	Animation func(e *goecs.Entity) Animation
	// Layer gets a Layer from a goecs.Entity
	State func(e *goecs.Entity) State
	// Layer gets a Layer from a goecs.Entity
	Sequence func(e *goecs.Entity) Sequence
}

// Get animation component
var Get = gets{
	// Animation gets a animation.Animation from a goecs.Entity
	Animation: func(e *goecs.Entity) Animation {
		return e.Get(TYPE.Animation).(Animation)
	},
	// State gets a animation.State from a goecs.Entity
	State: func(e *goecs.Entity) State {
		return e.Get(TYPE.State).(State)
	},
	// Sequence gets a animation.Sequence from a goecs.Entity
	Sequence: func(e *goecs.Entity) Sequence {
		return e.Get(TYPE.Sequence).(Sequence)
	},
}
