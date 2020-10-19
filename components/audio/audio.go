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

// Package audio are the components for audio and music
package audio

import (
	"github.com/juan-medina/goecs"
)

// Music represent an music stream
type Music struct {
	Name string // Name is the filename for our music stream
}

// Type return this goecs.ComponentType
func (m Music) Type() goecs.ComponentType {
	return TYPE.Music
}

// MusicPlayingState represent the playing state for a music stream
type MusicPlayingState int32

//goland:noinspection GoUnusedConst
const (
	StateStopped = MusicPlayingState(iota) // StateStopped is the audio.MusicPlayingState for a stopped music
	StatePlaying                           // StatePlaying is the audio.MusicPlayingState for a playing music
	StatePaused                            // StatePaused is the audio.MusicPlayingState for a paused music
)

// MusicState represent the state for a music
type MusicState struct {
	PlayingState MusicPlayingState // PlayingState is the current audio.MusicPlayingState for our music
	Name         string            // Name is the current music Name
}

// Type return this goecs.ComponentType
func (m MusicState) Type() goecs.ComponentType {
	return TYPE.MusicState
}

type types struct {
	// Music is the goecs.ComponentType for audio.Music
	Music goecs.ComponentType
	// MusicPlayingState is the goecs.ComponentType for audio.Music
	MusicState goecs.ComponentType
}

// TYPE hold the goecs.ComponentType for our audio components
var TYPE = types{
	Music:      goecs.NewComponentType(),
	MusicState: goecs.NewComponentType(),
}

type gets struct {
	// Music gets a audio.Music from a goecs.Entity
	Music func(e *goecs.Entity) Music
	// Music gets a audio.MusicState from a goecs.Entity
	MusicState func(e *goecs.Entity) MusicState
}

// Get a audio component
var Get = gets{
	// Music gets a audio.Music from a goecs.Entity
	Music: func(e *goecs.Entity) Music {
		return e.Get(TYPE.Music).(Music)
	},
	// MusicState gets a audio.MusicState from a goecs.Entity
	MusicState: func(e *goecs.Entity) MusicState {
		return e.Get(TYPE.MusicState).(MusicState)
	},
}
