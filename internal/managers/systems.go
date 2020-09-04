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

// Package managers are the world.System for the engine
package managers

import "github.com/juan-medina/goecs/pkg/world"

// Manager can have world.System or world.Listener
type Manager interface{}

// WithSystem have a world.System
type WithSystem interface {
	// System is a world.System
	System(wld *world.World, delta float32) error
}

// WithListener have a world.Listener
type WithListener interface {
	// Listener is a world.Listener
	Listener(wld *world.World, signal interface{}, delta float32) error
}

// WithSystemAndListener have a world.System and a world.Listener
type WithSystemAndListener interface {
	// System is a world.System
	System(wld *world.World, delta float32) error
	// Listener is a world.Listener
	Listener(wld *world.World, signal interface{}, delta float32) error
}
