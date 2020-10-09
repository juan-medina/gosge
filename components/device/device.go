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
	KeyEscape                // KeyEscape is the escape key
	KeyReturn                // KeyReturn is the return key
	KeyF1                    // KeyF1 is the F1 key
	KeyF2                    // KeyF2 is the F2 key
	KeyF3                    // KeyF3 is the F3 key
	KeyF4                    // KeyF4 is the F4 key
	KeyF5                    // KeyF5 is the F5 key
	KeyF6                    // KeyF6 is the F6 key
	KeyF7                    // KeyF7 is the F7 key
	KeyF8                    // KeyF8 is the F8 key
	KeyF9                    // KeyF9 is the F9 key
	KeyF10                   // KeyF10 is the F10 key
	KeyF11                   // KeyF11 is the F11 key
	KeyF12                   // KeyF12 is the F12 key
	TotalKeys                // TotalKeys is the total number of keys
)

// GamepadButton represents a gamepad button
type GamepadButton int32

//goland:noinspection GoUnusedConst
const (
	GamepadFirstButton   = GamepadButton(iota) // GamepadFirstButton initial value for game pad buttons
	GamepadUp                                  // GamepadUp DPAD UP
	GamepadRight                               // GamepadRight DPAD RIGHT
	GamepadDown                                // GamepadDown DPAD DOWN
	GamepadLeft                                // GamepadLeft DPAD LEFT
	GamepadButton1                             // GamepadButton1 XBOX Y, PS3 Triangle
	GamepadButton2                             // GamepadButton2 XBOX X, PS3 Square
	GamepadButton3                             // GamepadButton3 XBOX A, PS3 Cross
	GamepadButton4                             // GamepadButton4 XBOX B, PS3 Circle
	GamepadLeftTrigger1                        // GamepadLeftTrigger1 L1
	GamepadLeftTrigger2                        // GamepadLeftTrigger2 L2
	GamepadRightTrigger1                       // GamepadRightTrigger1 R1
	GamepadRightTrigger2                       // GamepadRightTrigger2 R2
	GamepadSelect                              // GamepadSelect XBOX Option, PS3 Select
	GamepadSpecial                             // GamepadSpecial XBOX X Button, PS3 PS Button
	GamepadStart                               // GamepadStart XBOX Start, PS3 Start
	GamepadLeftThumb                           // GamepadLeftThumb LEFT THUMB
	GamepadRightThumb                          // GamepadRightThumb RIGHT THUMB
	TotalButtons                               // TotalButtons total gamepad buttons
)

// GamepadStick represent a gamepad stick
type GamepadStick int32

//goland:noinspection GoUnusedConst
const (
	GamepadFirstStick = GamepadStick(iota) // GamepadFirstStick initial value for game pad stick
	GamepadLeftStick                       // GamepadLeftStick is the left stick
	GamepadRightStick                      // GamepadRightStick is the right stick
	TotalSticks                            // TotalSticks total gamepad sticks
)

// Device constants
const (
	MaxGamePads = 4 // Max Number of Gamepads
)
