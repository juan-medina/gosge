package managers

import (
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/pkg/components/device"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/events"
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
	rdr render.Render
	ks  []device.KeyStatus
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
	if em.rdr.ShouldClose() {
		if err := em.sendGameClose(world); err != nil {
			return err
		}
	}

	mp := em.rdr.GetMousePoint()
	if em.mme.Point != mp {
		em.mme.Point = mp
		if err := em.sendMouseMove(world); err != nil {
			return err
		}
	}

	for _, v := range mouseButtonsTocCheck {
		if em.rdr.IsMouseRelease(v) {
			if err := em.sendMouseRelease(world, v); err != nil {
				return err
			}
		}
	}

	for key := device.FirstKey + 1; key < device.TotalKeys; key++ {
		status := em.rdr.GetKeyStatus(key)
		if !status.Equals(em.ks[key]) || (status.Down) {
			em.ks[key] = status
			if err := em.sendKeyEvent(world, key, status); err != nil {
				return err
			}
		}
	}

	return nil
}

// Events returns a managers.WithSystem that will handle signals
func Events(rdr render.Render) WithSystem {
	return &eventManager{
		rdr: rdr,
		mme: events.MouseMoveEvent{
			Point: geometry.Point{
				X: -1,
				Y: -1,
			},
		},
		ks: make([]device.KeyStatus, device.TotalKeys),
	}
}
