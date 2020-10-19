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

package managers

import (
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components/audio"
	"github.com/juan-medina/gosge/events"
)

type musicManager struct {
	dm DeviceManager
	sm *StorageManager
}

func (mm musicManager) Signals() []goecs.ComponentType {
	return []goecs.ComponentType{
		events.TYPE.PlayMusicEvent,
		events.TYPE.StopMusicEvent,
		events.TYPE.PauseMusicEvent,
		events.TYPE.ResumeMusicEvent,
		events.TYPE.ChangeMusicVolumeEvent,
	}
}

func (mm musicManager) System(world *goecs.World, _ float32) error {
	for it := world.Iterator(audio.TYPE.Music, audio.TYPE.MusicState); it != nil; it = it.Next() {
		ent := it.Value()
		//music := audio.Get.Music(ent)
		state := audio.Get.MusicState(ent)

		if def, err := mm.sm.GetMusicDef(state.Name); err == nil {
			mm.dm.UpdateMusic(def)
		} else {
			return err
		}
	}
	return nil
}

func (mm musicManager) findMusicEnt(world *goecs.World, name string) *goecs.Entity {
	for it := world.Iterator(audio.TYPE.MusicState); it != nil; it = it.Next() {
		ent := it.Value()
		state := audio.Get.MusicState(ent)
		if state.Name == name {
			return ent
		}
	}
	return nil
}

func (mm musicManager) Listener(world *goecs.World, event goecs.Component, _ float32) error {
	switch e := event.(type) {
	case events.PlayMusicEvent:
		return mm.playMusicEvent(world, e)
	case events.StopMusicEvent:
		return mm.stopMusicEvent(world, e)
	case events.PauseMusicEvent:
		return mm.pauseMusicEvent(world, e)
	case events.ResumeMusicEvent:
		return mm.resumeMusicEvent(world, e)
	case events.ChangeMusicVolumeEvent:
		return mm.changeVolume(world, e)
	}
	return nil
}

func (mm musicManager) playMusicEvent(world *goecs.World, pme events.PlayMusicEvent) error {
	if def, err := mm.sm.GetMusicDef(pme.Name); err == nil {
		var ent *goecs.Entity
		if ent = mm.findMusicEnt(world, pme.Name); ent == nil {
			entID := world.AddEntity(
				audio.Music{
					Name: pme.Name,
				},
				audio.MusicState{
					Name:         pme.Name,
					PlayingState: audio.StateStopped,
				},
			)
			var err error
			if ent, err = world.Get(entID); err != nil {
				return err
			}
		}
		state := audio.Get.MusicState(ent)
		if state.PlayingState == audio.StateStopped || state.PlayingState == audio.StatePaused {
			old := state.PlayingState
			state.PlayingState = audio.StatePlaying
			ent.Set(state)
			mm.dm.PlayMusic(def, pme.Volume)
			mm.sendMusicStateChangeEvent(world, pme.Name, old, state.PlayingState)
		}
	} else {
		return err
	}
	return nil
}

func (mm musicManager) stopMusicEvent(world *goecs.World, sme events.StopMusicEvent) error {
	if def, err := mm.sm.GetMusicDef(sme.Name); err == nil {
		if ent := mm.findMusicEnt(world, sme.Name); ent != nil {
			state := audio.Get.MusicState(ent)
			if state.PlayingState == audio.StatePlaying || state.PlayingState == audio.StatePaused {
				old := state.PlayingState
				state.PlayingState = audio.StateStopped
				ent.Set(state)
				mm.dm.StopMusic(def)
				mm.sendMusicStateChangeEvent(world, sme.Name, old, state.PlayingState)
			}
		}
	} else {
		return err
	}
	return nil
}

func (mm musicManager) pauseMusicEvent(world *goecs.World, pme events.PauseMusicEvent) error {
	if def, err := mm.sm.GetMusicDef(pme.Name); err == nil {
		if ent := mm.findMusicEnt(world, pme.Name); ent != nil {
			state := audio.Get.MusicState(ent)
			if state.PlayingState == audio.StatePlaying {
				old := state.PlayingState
				state.PlayingState = audio.StatePaused
				ent.Set(state)
				mm.dm.PauseMusic(def)
				mm.sendMusicStateChangeEvent(world, pme.Name, old, state.PlayingState)
			}
		}
	} else {
		return err
	}
	return nil
}

func (mm musicManager) resumeMusicEvent(world *goecs.World, rme events.ResumeMusicEvent) error {
	if def, err := mm.sm.GetMusicDef(rme.Name); err == nil {
		if ent := mm.findMusicEnt(world, rme.Name); ent != nil {
			state := audio.Get.MusicState(ent)
			if state.PlayingState == audio.StatePaused {
				old := state.PlayingState
				state.PlayingState = audio.StatePlaying
				ent.Set(state)
				mm.dm.ResumeMusic(def)
				mm.sendMusicStateChangeEvent(world, rme.Name, old, state.PlayingState)
			}
		}
	} else {
		return err
	}
	return nil
}

func (mm musicManager) sendMusicStateChangeEvent(world *goecs.World, name string, old audio.MusicPlayingState, new audio.MusicPlayingState) {
	world.Signal(events.MusicStateChangeEvent{
		Name: name,
		Old:  old,
		New:  new,
	})
}

func (mm musicManager) changeVolume(world *goecs.World, cmv events.ChangeMusicVolumeEvent) error {
	if def, err := mm.sm.GetMusicDef(cmv.Name); err == nil {
		if ent := mm.findMusicEnt(world, cmv.Name); ent != nil {
			state := audio.Get.MusicState(ent)
			if state.PlayingState == audio.StatePlaying || state.PlayingState == audio.StatePaused {
				mm.dm.ChangeMusicVolume(def, cmv.Volume)
			}
		}
	}
	return nil
}

// Music returns a managers.WithSystemAndListener that handle music stream
func Music(dm DeviceManager, sm *StorageManager) WithSystemAndListener {
	return &musicManager{
		dm: dm,
		sm: sm,
	}
}
