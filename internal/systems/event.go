package systems

import (
	"github.com/juan-medina/goecs/pkg/world"
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

type eventSystem struct {
	mme events.MouseMoveEvent
	rdr render.Render
	ks  []device.KeyStatus
}

var (
	mouseButtonsTocCheck = []device.MouseButton{
		device.MouseLeftButton, device.MouseMiddleButton, device.MouseRightButton,
	}
)

func (es eventSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}

func (es eventSystem) sendGameClose(wld *world.World) error {
	return wld.Notify(events.GameCloseEvent{})
}

func (es eventSystem) sendMouseMove(wld *world.World) error {
	return wld.Notify(es.mme)
}

func (es eventSystem) sendMouseRelease(wld *world.World, button device.MouseButton) error {
	return wld.Notify(events.MouseUpEvent{Point: es.mme.Point, MouseButton: button})
}

func (es eventSystem) sendKeyEvent(wld *world.World, key device.Key, status device.KeyStatus) error {
	return wld.Notify(events.KeyEvent{Key: key, Status: status})
}

func (es *eventSystem) Update(world *world.World, _ float32) error {
	if es.rdr.ShouldClose() {
		if err := es.sendGameClose(world); err != nil {
			return err
		}
	}

	mp := es.rdr.GetMousePoint()
	if es.mme.Point != mp {
		es.mme.Point = mp
		if err := es.sendMouseMove(world); err != nil {
			return err
		}
	}

	for _, v := range mouseButtonsTocCheck {
		if es.rdr.IsMouseRelease(v) {
			if err := es.sendMouseRelease(world, v); err != nil {
				return err
			}
		}
	}

	for key := device.FirstKey + 1; key < device.TotalKeys; key++ {
		status := es.rdr.GetKeyStatus(key)
		if !status.Equals(es.ks[key]) || (status.Down || status.Up) {
			es.ks[key] = status
			if err := es.sendKeyEvent(world, key, status); err != nil {
				return err
			}
		}
	}

	return nil
}

// EventSystem returns a world.System that will handle events
func EventSystem(rdr render.Render) world.System {
	return &eventSystem{
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
