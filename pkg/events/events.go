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
	"github.com/juan-medina/gosge/pkg/components/audio"
	"github.com/juan-medina/gosge/pkg/components/device"
	"github.com/juan-medina/gosge/pkg/components/geometry"
)

// GameCloseEvent is an event that indicates that game need to close
type GameCloseEvent struct{}

// MouseMoveEvent is an event that indicates that the mouse is moving
type MouseMoveEvent struct {
	geometry.Point
}

// MouseUpEvent is an event that indicates that the mouse is release
type MouseUpEvent struct {
	// Point is the geometry.Point where the mouse is  when released
	Point geometry.Point
	// MouseButton is the device.MouseButton been released
	MouseButton device.MouseButton
}

// ChangeGameStage is an event that indicates that change game stage, all entities,
//systems, sprites sheets and textures will be removed. If the Stage does not exist
//the game.Run method will return an error. Stages must be created with engine.AddGameStage
type ChangeGameStage struct {
	// Stage is the name of the stage to change to, it must be created with engine.AddGameStage
	Stage string
}

// KeyEvent this event triggers when the status of a key change
type KeyEvent struct {
	Key    device.Key
	Status device.KeyStatus
}

// PlayMusicEvent is an event to play a music stream
type PlayMusicEvent struct {
	Name string
}

// StopMusicEvent is an event to play a music stream
type StopMusicEvent struct {
	Name string
}

// PauseMusicEvent is an event to pause a music stream
type PauseMusicEvent struct {
	Name string
}

// ResumeMusicEvent is an event to resume a music stream
type ResumeMusicEvent struct {
	Name string
}

// MusicStateChangeEvent is a event trigger when a music state change
type MusicStateChangeEvent struct {
	Name string                  // Name of the music
	Old  audio.MusicPlayingState // Old is the previous audio.MusicPlayingState
	New  audio.MusicPlayingState // New is the current audio.MusicPlayingState
}

// PlaySoundEvent is an event to play a sound wave
type PlaySoundEvent struct {
	Name string
}
