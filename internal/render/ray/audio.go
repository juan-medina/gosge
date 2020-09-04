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

package ray

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/juan-medina/gosge/internal/components"
	"os"
)

// LoadMusic giving it file name into memory
func (rr RenderImpl) LoadMusic(fileName string) (components.MusicDef, error) {
	var file *os.File
	var err error

	if file, err = os.Open(fileName); err != nil {
		return components.MusicDef{}, fmt.Errorf("we could not find the music file %q", fileName)
	}
	_ = file.Close()
	rlMusic := rl.LoadMusicStream(fileName)
	return components.MusicDef{
		Data: rlMusic,
	}, nil
}

// UnloadMusic giving it file from memory
func (rr RenderImpl) UnloadMusic(musicDef components.MusicDef) {
	rl.UnloadMusicStream(musicDef.Data.(rl.Music))
}

// PlayMusic plays the given components.MusicDef
func (rr RenderImpl) PlayMusic(musicDef components.MusicDef) {
	rlMusic := musicDef.Data.(rl.Music)
	rl.PlayMusicStream(rlMusic)
}

// PauseMusic pauses the given components.MusicDef
func (rr RenderImpl) PauseMusic(musicDef components.MusicDef) {
	rl.PauseMusicStream(musicDef.Data.(rl.Music))
}

// StopMusic stop the given components.MusicDef
func (rr RenderImpl) StopMusic(musicDef components.MusicDef) {
	rl.StopMusicStream(musicDef.Data.(rl.Music))
}

// ResumeMusic resumes the given components.MusicDef
func (rr RenderImpl) ResumeMusic(musicDef components.MusicDef) {
	rl.ResumeMusicStream(musicDef.Data.(rl.Music))
}

// UpdateMusic update the stream of the given components.MusicDef
func (rr RenderImpl) UpdateMusic(musicDef components.MusicDef) {
	rl.UpdateMusicStream(musicDef.Data.(rl.Music))
}

// LoadSound giving it file name into memory
func (rr *RenderImpl) LoadSound(fileName string) (components.SoundDef, error) {
	var file *os.File
	var err error

	if file, err = os.Open(fileName); err != nil {
		return components.SoundDef{}, fmt.Errorf("we could not find the sound file %q", fileName)
	}
	_ = file.Close()
	rlSound := rl.LoadSound(fileName)
	return components.SoundDef{
		Data: rlSound,
	}, nil
}

// UnloadSound giving it file from memory
func (rr *RenderImpl) UnloadSound(soundDef components.SoundDef) {
	rl.UnloadSound(soundDef.Data.(rl.Sound))
}

// PlaySound plays the given components.SoundDef
func (rr *RenderImpl) PlaySound(soundDef components.SoundDef) {
	rl.PlaySoundMulti(soundDef.Data.(rl.Sound))
}

// StopAllSounds currently playing
func (rr *RenderImpl) StopAllSounds() {
	rl.StopSoundMulti()
}
