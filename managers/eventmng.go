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
	ssm [][]geometry.Point
}

func (em eventManager) Signals() []goecs.ComponentType {
	return []goecs.ComponentType{events.TYPE.DelaySignal}
}

func (em eventManager) Listener(world *goecs.World, signal goecs.Component, delta float32) error {
	switch e := signal.(type) {
	case events.DelaySignal:
		e.Time -= delta
		if e.Time <= 0 {
			world.Signal(e.Signal)
			return nil
		}
		world.Signal(e)
	}
	return nil
}

var (
	mouseButtonsTocCheck = []device.MouseButton{
		device.MouseLeftButton, device.MouseMiddleButton, device.MouseRightButton,
	}
)

func (em eventManager) sendGameClose(world *goecs.World) {
	world.Signal(events.GameCloseEvent{})
}

func (em eventManager) sendMouseMove(world *goecs.World) {
	world.Signal(em.mme)
}

func (em eventManager) sendMouseRelease(world *goecs.World, button device.MouseButton) {
	world.Signal(events.MouseUpEvent{Point: em.mme.Point, MouseButton: button})
}

func (em eventManager) sendMousePressed(world *goecs.World, button device.MouseButton) {
	world.Signal(events.MouseDownEvent{Point: em.mme.Point, MouseButton: button})
}

func (em eventManager) sendKeyDownEvent(world *goecs.World, key device.Key) {
	world.Signal(events.KeyDownEvent{Key: key})
}

func (em eventManager) sendKeyUpEvent(world *goecs.World, key device.Key) {
	world.Signal(events.KeyUpEvent{Key: key})
}

func (em eventManager) sendGamePadButtonUpEvent(world *goecs.World, gamepad int32, button device.GamepadButton) {
	world.Signal(events.GamePadButtonUpEvent{Gamepad: gamepad, Button: button})
}

func (em eventManager) sendGamePadButtonDownEvent(world *goecs.World, gamepad int32, button device.GamepadButton) {
	world.Signal(events.GamePadButtonDownEvent{Gamepad: gamepad, Button: button})
}

func (em eventManager) sendGamePadStickMoveEvent(world *goecs.World, gamepad int32, stick device.GamepadStick,
	movement geometry.Point) {
	world.Signal(events.GamePadStickMoveEvent{Gamepad: gamepad, Stick: stick, Movement: movement})
}

func (em *eventManager) System(world *goecs.World, _ float32) error {
	if em.dm.ShouldClose() {
		em.sendGameClose(world)
	} else {
		em.handleMouse(world)
		em.handleKeys(world)
		em.handleGamepad(world)
	}
	return nil
}

func (em eventManager) handleMouse(world *goecs.World) {
	mp := em.dm.GetMousePoint()
	if em.mme.Point != mp {
		em.mme.Point = mp
		em.sendMouseMove(world)
	}

	for _, button := range mouseButtonsTocCheck {
		if em.dm.IsMouseRelease(button) {
			em.sendMouseRelease(world, button)
		}
		if em.dm.IsMousePressed(button) {
			em.sendMousePressed(world, button)
		}
	}
}

func (em eventManager) handleKeys(world *goecs.World) {
	for key := device.FirstKey + 1; key < device.TotalKeys; key++ {
		if em.dm.IsKeyReleased(key) {
			em.sendKeyUpEvent(world, key)
		}
		if em.dm.IsKeyPressed(key) {
			em.sendKeyDownEvent(world, key)
		}
	}
}

func (em *eventManager) handleGamepad(world *goecs.World) {
	for pad := int32(0); pad < device.MaxGamePads; pad++ {
		for button := device.GamepadFirstButton + 1; button < device.TotalButtons; button++ {
			if em.dm.IsGamepadButtonReleased(pad, button) {
				em.sendGamePadButtonUpEvent(world, pad, button)
			}
			if em.dm.IsGamepadButtonPressed(pad, button) {
				em.sendGamePadButtonDownEvent(world, pad, button)
			}
		}
		for stick := device.GamepadFirstStick + 1; stick < device.TotalSticks; stick++ {
			mov := em.dm.GetGamepadStickMovement(pad, stick)
			saved := em.ssm[pad][stick]
			if mov.X != saved.X || mov.Y != saved.Y {
				em.ssm[pad][stick] = mov
				em.sendGamePadStickMoveEvent(world, pad, stick, mov)
			}
		}
	}
}

// Events returns a managers.WithSystem that will handle signals
func Events(dm DeviceManager) WithSystemAndListener {
	// initialize saved stick movement
	var ssm = make([][]geometry.Point, device.MaxGamePads)
	for pad := int32(0); pad < device.MaxGamePads; pad++ {
		ssm[pad] = make([]geometry.Point, device.TotalSticks)
		for stick := device.GamepadFirstStick; stick < device.TotalSticks; stick++ {
			ssm[pad][int32(stick)] = geometry.Point{}
		}
	}

	return &eventManager{
		dm: dm,
		mme: events.MouseMoveEvent{
			Point: geometry.Point{
				X: -1,
				Y: -1,
			},
		},
		ssm: ssm,
	}
}
