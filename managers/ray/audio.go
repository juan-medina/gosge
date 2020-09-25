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
	"github.com/juan-medina/gosge/components"
	"os"
)

// LoadMusic giving it file name into memory
func (dmi DeviceManagerImpl) LoadMusic(fileName string) (components.MusicDef, error) {
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
func (dmi DeviceManagerImpl) UnloadMusic(musicDef components.MusicDef) {
	rl.UnloadMusicStream(musicDef.Data.(rl.Music))
}

// PlayMusic plays the given components.MusicDef
func (dmi DeviceManagerImpl) PlayMusic(musicDef components.MusicDef, volume float32) {
	rlMusic := musicDef.Data.(rl.Music)
	rl.PlayMusicStream(rlMusic)
	rl.SetMusicVolume(rlMusic, volume)
}

// ChangeMusicVolume change the given components.MusicDef volume
func (dmi DeviceManagerImpl) ChangeMusicVolume(musicDef components.MusicDef, volume float32) {
	rlMusic := musicDef.Data.(rl.Music)
	rl.SetMusicVolume(rlMusic, volume)
}

// PauseMusic pauses the given components.MusicDef
func (dmi DeviceManagerImpl) PauseMusic(musicDef components.MusicDef) {
	rl.PauseMusicStream(musicDef.Data.(rl.Music))
}

// StopMusic stop the given components.MusicDef
func (dmi DeviceManagerImpl) StopMusic(musicDef components.MusicDef) {
	rl.StopMusicStream(musicDef.Data.(rl.Music))
}

// ResumeMusic resumes the given components.MusicDef
func (dmi DeviceManagerImpl) ResumeMusic(musicDef components.MusicDef) {
	rl.ResumeMusicStream(musicDef.Data.(rl.Music))
}

// UpdateMusic update the stream of the given components.MusicDef
func (dmi DeviceManagerImpl) UpdateMusic(musicDef components.MusicDef) {
	rl.UpdateMusicStream(musicDef.Data.(rl.Music))
}

// LoadSound giving it file name into memory
func (dmi *DeviceManagerImpl) LoadSound(fileName string) (components.SoundDef, error) {
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
func (dmi *DeviceManagerImpl) UnloadSound(soundDef components.SoundDef) {
	rl.UnloadSound(soundDef.Data.(rl.Sound))
}

// PlaySound plays the given components.SoundDef
func (dmi *DeviceManagerImpl) PlaySound(soundDef components.SoundDef, volume float32) {
	rlSound := soundDef.Data.(rl.Sound)
	rl.SetSoundVolume(rlSound, volume)
	rl.PlaySoundMulti(rlSound)
}

// StopAllSounds currently playing
func (dmi *DeviceManagerImpl) StopAllSounds() {
	rl.StopSoundMulti()
}

// SetMasterVolume change the master volume
func (dmi *DeviceManagerImpl) SetMasterVolume(volume float32) {
	rl.SetMasterVolume(volume)
}
