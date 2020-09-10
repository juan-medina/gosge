package managers

import (
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/events"
)

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

type eventManager struct {
	mme events.MouseMoveEvent
	dm  DeviceManager
	ks  []device.KeyStatus
}

func (em eventManager) Listener(world *goecs.World, signal interface{}, delta float32) error {
	switch e := signal.(type) {
	case events.DelaySignal:
		e.Time -= delta
		if e.Time <= 0 {
			return world.Signal(e.Signal)
		}
		return world.Signal(e)
	}
	return nil
}

var (
	mouseButtonsTocCheck = []device.MouseButton{
		device.MouseLeftButton, device.MouseMiddleButton, device.MouseRightButton,
	}
)

func (em eventManager) sendGameClose(world *goecs.World) error {
	return world.Signal(events.GameCloseEvent{})
}

func (em eventManager) sendMouseMove(world *goecs.World) error {
	return world.Signal(em.mme)
}

func (em eventManager) sendMouseRelease(world *goecs.World, button device.MouseButton) error {
	return world.Signal(events.MouseUpEvent{Point: em.mme.Point, MouseButton: button})
}

func (em eventManager) sendKeyEvent(world *goecs.World, key device.Key, status device.KeyStatus) error {
	return world.Signal(events.KeyEvent{Key: key, Status: status})
}

func (em *eventManager) System(world *goecs.World, _ float32) error {
	if em.dm.ShouldClose() {
		if err := em.sendGameClose(world); err != nil {
			return err
		}
	}

	mp := em.dm.GetMousePoint()
	if em.mme.Point != mp {
		em.mme.Point = mp
		if err := em.sendMouseMove(world); err != nil {
			return err
		}
	}

	for _, v := range mouseButtonsTocCheck {
		if em.dm.IsMouseRelease(v) {
			if err := em.sendMouseRelease(world, v); err != nil {
				return err
			}
		}
	}

	for key := device.FirstKey + 1; key < device.TotalKeys; key++ {
		status := em.dm.GetKeyStatus(key)
		if !status.Equals(em.ks[key]) {
			em.ks[key] = status
			if err := em.sendKeyEvent(world, key, status); err != nil {
				return err
			}
		}
	}

	return nil
}

// Events returns a managers.WithSystem that will handle signals
func Events(dm DeviceManager) WithSystemAndListener {
	return &eventManager{
		dm: dm,
		mme: events.MouseMoveEvent{
			Point: geometry.Point{
				X: -1,
				Y: -1,
			},
		},
		ks: make([]device.KeyStatus, device.TotalKeys),
	}
}
