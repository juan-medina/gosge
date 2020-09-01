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
	"github.com/juan-medina/goecs/pkg/entity"
	"reflect"
)

//goland:noinspection GoUnusedConst
const (
	LoopForever = -1 // LoopForever indicates to loop the audio forever
)

// Music represent an music stream
type Music struct {
	Name  string // Name is the filename for our music stream
	Loops int32  // Loops is the number of loops, we could use audio.LoopForever
}

// MusicPlayingState represent the playing state for a music stream
type MusicPlayingState int32

//goland:noinspection GoUnusedConst
const (
	Stopped = MusicPlayingState(iota) // Stopped is the audio.MusicPlayingState for a stopped music
	Playing                           // Playing is the audio.MusicPlayingState for a playing music
	Paused                            // Paused is the audio.MusicPlayingState for a paused music
)

// MusicState represent the state for a music
type MusicState struct {
	State MusicPlayingState // State is the current audio.MusicState for our music
	Name  string            // Name is the current music Name
}

type types struct {
	// Music is the reflect.Type for audio.Music
	Music reflect.Type
	// MusicPlayingState is the reflect.Type for audio.Music
	MusicState reflect.Type
}

// TYPE hold the reflect.Type for our audio components
var TYPE = types{
	Music:      reflect.TypeOf(Music{}),
	MusicState: reflect.TypeOf(MusicState{}),
}

type gets struct {
	// Music gets a audio.Music from a entity.Entity
	Music func(e *entity.Entity) Music
	// Music gets a audio.MusicState from a entity.Entity
	MusicState func(e *entity.Entity) MusicState
}

// Get a audio component
var Get = gets{
	// Music gets a audio.Music from a entity.Entity
	Music: func(e *entity.Entity) Music {
		return e.Get(TYPE.Music).(Music)
	},
	// MusicState gets a audio.MusicState from a entity.Entity
	MusicState: func(e *entity.Entity) MusicState {
		return e.Get(TYPE.MusicState).(MusicState)
	},
}
