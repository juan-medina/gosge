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

// Package events contain the events for our engine
package events

import (
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components/audio"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/geometry"
)

// GameCloseEvent is an event that indicates that game need to close
type GameCloseEvent struct{}

// Type is this goecs.ComponentType
func (g GameCloseEvent) Type() goecs.ComponentType {
	return TYPE.GameCloseEvent
}

// MouseMoveEvent is an event that indicates that the mouse is moving
type MouseMoveEvent struct {
	geometry.Point
}

// Type is this goecs.ComponentType
func (m MouseMoveEvent) Type() goecs.ComponentType {
	return TYPE.MouseMoveEvent
}

// MouseUpEvent is an event that indicates that the mouse is release
type MouseUpEvent struct {
	// Point is the geometry.Point where the mouse is when released
	Point geometry.Point
	// MouseButton is the device.MouseButton been released
	MouseButton device.MouseButton
}

// Type is this goecs.ComponentType
func (m MouseUpEvent) Type() goecs.ComponentType {
	return TYPE.MouseUpEvent
}

// MouseDownEvent is an event that indicates that the mouse is pressed
type MouseDownEvent struct {
	// Point is the geometry.Point where the mouse is when pressed
	Point geometry.Point
	// MouseButton is the device.MouseButton been pressed
	MouseButton device.MouseButton
}

// Type is this goecs.ComponentType
func (m MouseDownEvent) Type() goecs.ComponentType {
	return TYPE.MouseDownEvent
}

// ChangeGameStage is an event that indicates that change game stage, all entities,
//systems, sprites sheets and textures will be removed. If the Stage does not exist
//the game.Run method will return an error. Stages must be created with engine.AddGameStage
type ChangeGameStage struct {
	// Stage is the name of the stage to change to, it must be created with engine.AddGameStage
	Stage string
}

// Type is this goecs.ComponentType
func (c ChangeGameStage) Type() goecs.ComponentType {
	return TYPE.ChangeGameStage
}

// KeyUpEvent this event triggers when a key is up
type KeyUpEvent struct {
	Key device.Key
}

// Type is this goecs.ComponentType
func (k KeyUpEvent) Type() goecs.ComponentType {
	return TYPE.KeyUpEvent
}

// KeyDownEvent this event triggers when a key is down
type KeyDownEvent struct {
	Key device.Key
}

// Type is this goecs.ComponentType
func (k KeyDownEvent) Type() goecs.ComponentType {
	return TYPE.KeyDownEvent
}

// PlayMusicEvent is an event to play a music stream
type PlayMusicEvent struct {
	Name   string
	Volume float32
}

// Type is this goecs.ComponentType
func (p PlayMusicEvent) Type() goecs.ComponentType {
	return TYPE.PlayMusicEvent
}

// StopMusicEvent is an event to play a music stream
type StopMusicEvent struct {
	Name string
}

// Type is this goecs.ComponentType
func (s StopMusicEvent) Type() goecs.ComponentType {
	return TYPE.StopMusicEvent
}

// PauseMusicEvent is an event to pause a music stream
type PauseMusicEvent struct {
	Name string
}

// Type is this goecs.ComponentType
func (p PauseMusicEvent) Type() goecs.ComponentType {
	return TYPE.PauseMusicEvent
}

// ResumeMusicEvent is an event to resume a music stream
type ResumeMusicEvent struct {
	Name string
}

// Type is this goecs.ComponentType
func (r ResumeMusicEvent) Type() goecs.ComponentType {
	return TYPE.ResumeMusicEvent
}

// ChangeMusicVolumeEvent is an event to change the volume of a music stream
type ChangeMusicVolumeEvent struct {
	Name   string
	Volume float32
}

// Type is this goecs.ComponentType
func (c ChangeMusicVolumeEvent) Type() goecs.ComponentType {
	return TYPE.ChangeMusicVolumeEvent
}

// MusicStateChangeEvent is a event trigger when a music state change
type MusicStateChangeEvent struct {
	Name string                  // Name of the music
	Old  audio.MusicPlayingState // Old is the previous audio.MusicPlayingState
	New  audio.MusicPlayingState // New is the current audio.MusicPlayingState
}

// Type is this goecs.ComponentType
func (m MusicStateChangeEvent) Type() goecs.ComponentType {
	return TYPE.MusicStateChangeEvent
}

// PlaySoundEvent is an event to play a sound wave
type PlaySoundEvent struct {
	Name   string
	Volume float32
}

// Type is this goecs.ComponentType
func (p PlaySoundEvent) Type() goecs.ComponentType {
	return TYPE.PlaySoundEvent
}

// ChangeMasterVolumeEvent to a given Volume
type ChangeMasterVolumeEvent struct {
	Volume float32 // Volume to be set
}

// Type is this goecs.ComponentType
func (c ChangeMasterVolumeEvent) Type() goecs.ComponentType {
	return TYPE.ChangeMasterVolumeEvent
}

// DelaySignal is a signal that will happen after a time
type DelaySignal struct {
	Signal interface{} // Signal that will be emitted after the Time
	Time   float32     // Time to this Signal to be emitted, in seconds
}

// Type is this goecs.ComponentType
func (d DelaySignal) Type() goecs.ComponentType {
	return TYPE.DelaySignal
}

// FocusOnControlEvent is a signal to change the focussed control
type FocusOnControlEvent struct {
	Control goecs.EntityID
}

// Type is this goecs.ComponentType
func (f FocusOnControlEvent) Type() goecs.ComponentType {
	return TYPE.FocusOnControlEvent
}

// ClearFocusEvent is a signal to clear all focussed control
type ClearFocusEvent struct{}

// Type is this goecs.ComponentType
func (c ClearFocusEvent) Type() goecs.ComponentType {
	return TYPE.ClearFocusEvent
}

// GamePadButtonUpEvent this event triggers when a gamepad button is up
type GamePadButtonUpEvent struct {
	Gamepad int32
	Button  device.GamepadButton
}

// Type is this goecs.ComponentType
func (g GamePadButtonUpEvent) Type() goecs.ComponentType {
	return TYPE.GamePadButtonUpEvent
}

// GamePadButtonDownEvent this event triggers when a gamepad button is down
type GamePadButtonDownEvent struct {
	Gamepad int32
	Button  device.GamepadButton
}

// Type is this goecs.ComponentType
func (g GamePadButtonDownEvent) Type() goecs.ComponentType {
	return TYPE.GamePadButtonDownEvent
}

// GamePadStickMoveEvent this event triggers when a gamepad stick moves
type GamePadStickMoveEvent struct {
	Gamepad  int32
	Stick    device.GamepadStick
	Movement geometry.Point
}

// Type is this goecs.ComponentType
func (g GamePadStickMoveEvent) Type() goecs.ComponentType {
	return TYPE.GamePadStickMoveEvent
}

type types struct {
	// GameCloseEvent is the goecs.ComponentType for events.GameCloseEvent
	GameCloseEvent goecs.ComponentType
	// ChangeGameStage is the goecs.ComponentType for events.ChangeGameStage
	ChangeGameStage goecs.ComponentType
	// DelaySignal is the goecs.ComponentType for events.DelaySignal
	DelaySignal goecs.ComponentType
	// PlaySoundEvent is the goecs.ComponentType for events.PlaySoundEvent
	PlaySoundEvent goecs.ComponentType
	// ChangeMasterVolumeEvent is the goecs.ComponentType for events.ChangeMasterVolumeEvent
	ChangeMasterVolumeEvent goecs.ComponentType
	// PlayMusicEvent is the goecs.ComponentType for events.PlayMusicEvent
	PlayMusicEvent goecs.ComponentType
	// StopMusicEvent is the goecs.ComponentType for events.StopMusicEvent
	StopMusicEvent goecs.ComponentType
	// PauseMusicEvent is the goecs.ComponentType for events.PauseMusicEvent
	PauseMusicEvent goecs.ComponentType
	// ResumeMusicEvent is the goecs.ComponentType for events.ResumeMusicEvent
	ResumeMusicEvent goecs.ComponentType
	// ChangeMusicVolumeEvent is the goecs.ComponentType for events.ChangeMusicVolumeEvent
	ChangeMusicVolumeEvent goecs.ComponentType
	// MouseMoveEvent is the goecs.ComponentType for events.MouseMoveEvent
	MouseMoveEvent goecs.ComponentType
	// MouseDownEvent is the goecs.ComponentType for events.MouseDownEvent
	MouseDownEvent goecs.ComponentType
	// MouseUpEvent is the goecs.ComponentType for events.MouseUpEvent
	MouseUpEvent goecs.ComponentType
	// KeyDownEvent is the goecs.ComponentType for events.KeyDownEvent
	KeyDownEvent goecs.ComponentType
	// KeyUpEvent is the goecs.ComponentType for events.KeyUpEvent
	KeyUpEvent goecs.ComponentType
	// MusicStateChangeEvent is the goecs.ComponentType for events.MusicStateChangeEvent
	MusicStateChangeEvent goecs.ComponentType
	// FocusOnControlEvent is the goecs.ComponentType for events.FocusOnControlEvent
	FocusOnControlEvent goecs.ComponentType
	// ClearFocusEvent is the goecs.ComponentType for events.ClearFocusEvent
	ClearFocusEvent goecs.ComponentType
	// GamePadButtonUpEvent is the goecs.ComponentType for events.GamePadButtonUpEvent
	GamePadButtonUpEvent goecs.ComponentType
	// GamePadButtonDownEvent is the goecs.ComponentType for events.GamePadButtonDownEvent
	GamePadButtonDownEvent goecs.ComponentType
	// GamePadStickMoveEvent is the goecs.ComponentType for events.GamePadStickMoveEvent
	GamePadStickMoveEvent goecs.ComponentType
}

// TYPE hold the goecs.ComponentType for our events
var TYPE = types{
	GameCloseEvent:          goecs.NewComponentType(),
	ChangeGameStage:         goecs.NewComponentType(),
	DelaySignal:             goecs.NewComponentType(),
	PlaySoundEvent:          goecs.NewComponentType(),
	ChangeMasterVolumeEvent: goecs.NewComponentType(),
	PlayMusicEvent:          goecs.NewComponentType(),
	StopMusicEvent:          goecs.NewComponentType(),
	PauseMusicEvent:         goecs.NewComponentType(),
	ResumeMusicEvent:        goecs.NewComponentType(),
	ChangeMusicVolumeEvent:  goecs.NewComponentType(),
	MouseMoveEvent:          goecs.NewComponentType(),
	MouseDownEvent:          goecs.NewComponentType(),
	MouseUpEvent:            goecs.NewComponentType(),
	KeyDownEvent:            goecs.NewComponentType(),
	KeyUpEvent:              goecs.NewComponentType(),
	MusicStateChangeEvent:   goecs.NewComponentType(),
	FocusOnControlEvent:     goecs.NewComponentType(),
	ClearFocusEvent:         goecs.NewComponentType(),
	GamePadButtonUpEvent:    goecs.NewComponentType(),
	GamePadButtonDownEvent:  goecs.NewComponentType(),
	GamePadStickMoveEvent:   goecs.NewComponentType(),
}
