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

func (m musicSystem) Notify(wld *world.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	case events.PlayMusicEvent:
		if def, err := m.ds.GetMusicDef(e.Name); err == nil {
			m.rdr.PlayMusic(def, e.Loops)
			wld.Add(entity.New(
				audio.Music{
					Name:  e.Name,
					Loops: e.Loops,
				},
				audio.MusicState{
					Name:  e.Name,
					State: audio.Playing,
				},
			))
		} else {
			return err
		}
	}

	return nil
}

// MusicSystem returns a world.System that handle music stream
func MusicSystem(rdr render.Render, ds storage.Storage) world.System {
	return &musicSystem{
		rdr: rdr,
		ds:  ds,
	}
}
