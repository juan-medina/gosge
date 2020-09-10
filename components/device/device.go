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

// Package device include device related types
package device

// MouseButton indicates a mouse button
type MouseButton int

//goland:noinspection GoUnusedConst
const (
	MouseLeftButton   = MouseButton(0) // MouseLeftButton is the left mouse button
	MouseRightButton  = MouseButton(1) // MouseRightButton is the right mouse button
	MouseMiddleButton = MouseButton(2) // MouseMiddleButton is the middle mouse button
)

// Key represents a device key
type Key int

//goland:noinspection GoUnusedConst
const (
	FirstKey     = Key(iota) // FirstKey is the initial value for keys
	KeyLeft                  // KeyLeft if the left cursor key
	KeyRight                 // KeyRight if the right cursor key
	KeyUp                    // KeyUp if the up cursor key
	KeyDown                  // KeyDown if the down cursor key
	KeySpace                 // KeySpace if the space cursor key
	KeyAltLeft               // KeySpace if the left alt key
	KeyCtrlLeft              // KeyCtrlLeft if the left control key
	KeyAltRight              // KeyAltRight if the right alt key
	KeyCtrlRight             // KeyCtrlRight if the right control key
	TotalKeys                // TotalKeys is the total number of keys
)
