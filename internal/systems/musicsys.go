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

package systems

import (
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/internal/storage"
	"github.com/juan-medina/gosge/pkg/components/audio"
	"github.com/juan-medina/gosge/pkg/events"
)

type musicSystem struct {
	rdr render.Render
	ds  storage.Storage
}

func (m musicSystem) Update(wld *world.World, _ float32) error {
	for it := wld.Iterator(audio.TYPE.Music, audio.TYPE.MusicState); it.HasNext(); {
		ent := it.Value()
		//music := audio.Get.Music(ent)
		state := audio.Get.MusicState(ent)

		if def, err := m.ds.GetMusicDef(state.Name); err == nil {
			m.rdr.UpdateMusic(def)
		} else {
			return err
		}
	}
	return nil
}

func (m musicSystem) findMusicEnt(wld *world.World, name string) *entity.Entity {
	for it := wld.Iterator(audio.TYPE.MusicState); it.HasNext(); {
		ent := it.Value()
		state := audio.Get.MusicState(ent)
		if state.Name == name {
			return ent
		}
	}
	return nil
}

func (m musicSystem) Notify(wld *world.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	case events.PlayMusicEvent:
		return m.playMusicEvent(wld, e)
	case events.StopMusicEvent:
		return m.stopMusicEvent(wld, e)
	case events.PauseMusicEvent:
		return m.pauseMusicEvent(wld, e)
	case events.ResumeMusicEvent:
		return m.resumeMusicEvent(wld, e)
	}
	return nil
}

func (m musicSystem) playMusicEvent(wld *world.World, pme events.PlayMusicEvent) error {
	if def, err := m.ds.GetMusicDef(pme.Name); err == nil {
		var ent *entity.Entity
		if ent = m.findMusicEnt(wld, pme.Name); ent == nil {
			ent = wld.Add(entity.New(
				audio.Music{
					Name:  pme.Name,
					Loops: pme.Loops,
				},
				audio.MusicState{
					Name:         pme.Name,
					PlayingState: audio.StateStopped,
				},
			))
		}
		state := audio.Get.MusicState(ent)
		if state.PlayingState == audio.StateStopped || state.PlayingState == audio.StatePaused {
			old := state.PlayingState
			state.PlayingState = audio.StatePlaying
			ent.Set(state)
			m.rdr.PlayMusic(def, pme.Loops)
			return m.sendMusicStateChangeEvent(wld, pme.Name, old, state.PlayingState)
		}
	} else {
		return err
	}
	return nil
}

func (m musicSystem) stopMusicEvent(wld *world.World, sme events.StopMusicEvent) error {
	if def, err := m.ds.GetMusicDef(sme.Name); err == nil {
		if ent := m.findMusicEnt(wld, sme.Name); ent != nil {
			state := audio.Get.MusicState(ent)
			if state.PlayingState == audio.StatePlaying || state.PlayingState == audio.StatePaused {
				old := state.PlayingState
				state.PlayingState = audio.StateStopped
				ent.Set(state)
				m.rdr.StopMusic(def)
				return m.sendMusicStateChangeEvent(wld, sme.Name, old, state.PlayingState)
			}
		}
	} else {
		return err
	}
	return nil
}

func (m musicSystem) pauseMusicEvent(wld *world.World, pme events.PauseMusicEvent) error {
	if def, err := m.ds.GetMusicDef(pme.Name); err == nil {
		if ent := m.findMusicEnt(wld, pme.Name); ent != nil {
			state := audio.Get.MusicState(ent)
			if state.PlayingState == audio.StatePlaying {
				old := state.PlayingState
				state.PlayingState = audio.StatePaused
				ent.Set(state)
				m.rdr.PauseMusic(def)
				return m.sendMusicStateChangeEvent(wld, pme.Name, old, state.PlayingState)
			}
		}
	} else {
		return err
	}
	return nil
}

func (m musicSystem) resumeMusicEvent(wld *world.World, rme events.ResumeMusicEvent) error {
	if def, err := m.ds.GetMusicDef(rme.Name); err == nil {
		if ent := m.findMusicEnt(wld, rme.Name); ent != nil {
			state := audio.Get.MusicState(ent)
			if state.PlayingState == audio.StatePaused {
				old := state.PlayingState
				state.PlayingState = audio.StatePlaying
				ent.Set(state)
				m.rdr.ResumeMusic(def)
				return m.sendMusicStateChangeEvent(wld, rme.Name, old, state.PlayingState)
			}
		}
	} else {
		return err
	}
	return nil
}

func (m musicSystem) sendMusicStateChangeEvent(wld *world.World, name string, old audio.MusicPlayingState, new audio.MusicPlayingState) error {
	return wld.Notify(events.MusicStateChangeEvent{
		Name: name,
		Old:  old,
		New:  new,
	})
}

// MusicSystem returns a world.System that handle music stream
func MusicSystem(rdr render.Render, ds storage.Storage) world.System {
	return &musicSystem{
		rdr: rdr,
		ds:  ds,
	}
}
