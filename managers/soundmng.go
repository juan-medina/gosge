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
	"github.com/juan-medina/gosge/events"
)

type soundManager struct {
	dm DeviceManager
	sm *StorageManager
}

func (sm soundManager) Signals() []goecs.ComponentType {
	return []goecs.ComponentType{events.TYPE.PlaySoundEvent, events.TYPE.ChangeMasterVolumeEvent}
}

func (sm soundManager) Listener(_ *goecs.World, event goecs.Component, _ float32) error {
	switch e := event.(type) {
	case events.PlaySoundEvent:
		if def, err := sm.sm.GetSoundDef(e.Name); err == nil {
			sm.dm.PlaySound(def, e.Volume)
		} else {
			return err
		}
	case events.ChangeMasterVolumeEvent:
		sm.dm.SetMasterVolume(e.Volume)
	}
	return nil
}

// Sounds returns a managers.WithListener that handle audio waves
func Sounds(dm DeviceManager, sm *StorageManager) WithListener {
	return &soundManager{
		dm: dm,
		sm: sm,
	}
}
