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
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"math"
	"reflect"
)

const (
	normalColorDarkenFactor = 0.25 // normalColorDarkenFactor is how much dark the color will be on normal
	hoverColorDarkenFactor  = 0.10 // hoverColorDarkenFactor is how much dark the color will be on hover
	keyWaitDelay            = 0.20 // keyWaitDelay is the key wait delay
)

type uiManager struct {
	dm       DeviceManager
	cm       *CollisionManager
	clicked  *goecs.Entity
	focus    *goecs.Entity
	lastKey  device.Key
	keyDelay float32
}

func (uim *uiManager) Signals() []reflect.Type {
	return []reflect.Type{
		events.TYPE.MouseMoveEvent,
		events.TYPE.MouseDownEvent,
		events.TYPE.MouseUpEvent,
		events.TYPE.FocusOnControlEvent,
		events.TYPE.ClearFocusEvent,
		events.TYPE.KeyUpEvent,
		events.TYPE.KeyDownEvent,
		events.TYPE.GamePadButtonUpEvent,
		events.TYPE.GamePadButtonDownEvent,
		events.TYPE.GamePadStickMoveEvent,
	}
}

func (uim *uiManager) System(world *goecs.World, delta float32) error {
	uim.flatButtons(world)
	uim.progressBars(world)
	uim.spriteButtons(world)
	uim.handleKeys(world, delta)
	return nil
}

func (uim *uiManager) Listener(world *goecs.World, event interface{}, _ float32) error {
	switch v := event.(type) {
	case events.MouseMoveEvent:
		uim.flatButtonsMouseMove(world, v)
		uim.progressBarsMouseMove(world, v)
		uim.spriteButtonsMouseMove(world, v)
	case events.MouseDownEvent:
		uim.flatButtonsMouseDown(world, v)
		uim.progressBarsMouseDown(world, v)
		uim.spriteButtonsMouseDown(world, v)
	case events.MouseUpEvent:
		uim.flatButtonsMouseUp(world, v)
		uim.progressBarsMouseUp(world, v)
		uim.spriteButtonsMouseUp(world, v)
	case events.FocusOnControlEvent:
		uim.focusControl(world, v.Control)
	case events.ClearFocusEvent:
		uim.clearFocus(world)
	case events.KeyDownEvent:
		uim.keyDown(world, v.Key)
	case events.GamePadButtonDownEvent:
		uim.gamepadDown(world, v.Gamepad, v.Button)
	case events.KeyUpEvent:
		uim.keyUp(v.Key)
	case events.GamePadButtonUpEvent:
		uim.gamepadUp(v.Gamepad, v.Button)
	case events.GamePadStickMoveEvent:
		uim.gamepadStickMove(world, v.Gamepad, v.Stick, v.Movement)
	}
	return nil
}

func (uim *uiManager) flatButtons(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.FlatButton); it != nil; it = it.Next() {
		ent := it.Value()
		// calculate state if is not done yet
		if ent.NotContains(ui.TYPE.ControlState) {
			ent.Set(ui.ControlState{})
		}
		if ent.NotContains(ui.TYPE.ButtonHoverColors) {
			clr := ui.Get.ButtonColor(ent)
			bcl := ui.ButtonHoverColors{}

			bcl.Hover.Border = clr.Border.Blend(color.Black, hoverColorDarkenFactor)
			bcl.Normal.Border = clr.Border.Blend(color.Black, normalColorDarkenFactor)
			bcl.Clicked.Border = clr.Border
			bcl.Disabled.Border = clr.Border.GrayScale().Blend(color.Black, normalColorDarkenFactor)

			bcl.Hover.Text = clr.Text.Blend(color.Black, hoverColorDarkenFactor)
			bcl.Normal.Text = clr.Text.Blend(color.Black, normalColorDarkenFactor)
			bcl.Clicked.Text = clr.Text
			bcl.Disabled.Text = clr.Text.GrayScale().Blend(color.Black, normalColorDarkenFactor)

			if clr.Gradient.From.Equals(clr.Gradient.To) {
				bcl.Hover.Solid = clr.Solid.Blend(color.Black, hoverColorDarkenFactor)
				bcl.Normal.Solid = clr.Solid.Blend(color.Black, normalColorDarkenFactor)
				bcl.Clicked.Solid = clr.Solid
				bcl.Disabled.Solid = clr.Solid.GrayScale().Blend(color.Black, normalColorDarkenFactor)
			} else {
				bcl.Hover.Gradient = color.Gradient{
					From:      clr.Gradient.From.Blend(color.Black, hoverColorDarkenFactor),
					To:        clr.Gradient.To.Blend(color.Black, hoverColorDarkenFactor),
					Direction: clr.Gradient.Direction,
				}
				bcl.Normal.Gradient = color.Gradient{
					From:      clr.Gradient.From.Blend(color.Black, normalColorDarkenFactor),
					To:        clr.Gradient.To.Blend(color.Black, normalColorDarkenFactor),
					Direction: clr.Gradient.Direction,
				}
				bcl.Clicked.Gradient = clr.Gradient
				bcl.Disabled.Gradient = color.Gradient{
					From:      clr.Gradient.From.GrayScale().Blend(color.Black, normalColorDarkenFactor),
					To:        clr.Gradient.To.GrayScale().Blend(color.Black, normalColorDarkenFactor),
					Direction: clr.Gradient.Direction,
				}
			}
			ent.Set(bcl)
		} else {
			uim.flatButtonsColors(ent)
		}
	}
}

func (uim uiManager) flatButtonsColors(ent *goecs.Entity) {
	bcl := ui.Get.ButtonHoverColors(ent)
	state := ui.Get.ControlState(ent)
	if state.Disabled {
		ent.Set(bcl.Disabled)
	} else if state.Clicked {
		ent.Set(bcl.Clicked)
	} else if state.Hover {
		ent.Set(bcl.Hover)
	} else {
		ent.Set(bcl.Normal)
	}
	if state.Focused {
		if ent.Contains(color.TYPE.Solid) {
			ac := color.Get.Solid(ent)
			clr := ui.Get.ButtonColor(ent)
			clr.Border = ac
			clr.Text = ac
			ent.Set(clr)
		}
	}
}

func (uim *uiManager) flatButtonsMouseMove(world *goecs.World, mme events.MouseMoveEvent) {
	if uim.clicked != nil {
		return
	}
	for it := world.Iterator(ui.TYPE.FlatButton, ui.TYPE.ButtonHoverColors, geometry.TYPE.Point,
		shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		btn := ui.Get.FlatButton(ent)
		state := ui.Get.ControlState(ent)
		if state.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)
		state.Hover = box.Contains(pos, mme.Point)
		uim.flatButtonsColors(ent)
		ent.Set(btn)
		ent.Set(state)
	}
}

func (uim *uiManager) flatButtonsMouseDown(world *goecs.World, mde events.MouseDownEvent) {
	for it := world.Iterator(ui.TYPE.FlatButton, geometry.TYPE.Point, shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		btn := ui.Get.FlatButton(ent)
		state := ui.Get.ControlState(ent)
		if state.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)
		if box.Contains(pos, mde.Point) {
			state.Clicked = true
			uim.clicked = ent
		}
		uim.flatButtonsColors(ent)
		ent.Set(btn)
		ent.Set(state)
	}
}

func (uim *uiManager) flatButtonsMouseUp(world *goecs.World, _ events.MouseUpEvent) {
	for it := world.Iterator(ui.TYPE.FlatButton, geometry.TYPE.Point, shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		btn := ui.Get.FlatButton(ent)
		state := ui.Get.ControlState(ent)
		if state.Clicked {
			trigger := true
			if btn.CheckBox {
				if btn.Group == "" {
					state.Checked = !state.Checked
				} else {
					if !state.Checked {
						uim.uncheckGroup(world, btn.Group)
						state.Checked = true
					} else {
						trigger = false
					}
				}
			}
			state.Clicked = false
			if trigger {
				if btn.Sound != "" {
					world.Signal(events.PlaySoundEvent{Name: btn.Sound, Volume: btn.Volume})
				}
				world.Signal(btn.Event)
			}

			uim.clicked = nil
		}
		uim.flatButtonsColors(ent)
		ent.Set(btn)
		ent.Set(state)
	}
}

func (uim *uiManager) uncheckGroup(world *goecs.World, group string) {
	for it := world.Iterator(ui.TYPE.FlatButton, geometry.TYPE.Point, shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		btn := ui.Get.FlatButton(ent)
		state := ui.Get.ControlState(ent)
		if btn.Group == group {
			state.Checked = false
			ent.Set(state)
		}
	}
}

func (uim uiManager) progressBars(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.ProgressBar); it != nil; it = it.Next() {
		ent := it.Value()
		// calculate state if is not done yet
		if ent.NotContains(ui.TYPE.ControlState) {
			ent.Set(ui.ControlState{})
		}
		if ent.NotContains(ui.TYPE.ProgressBarHoverColor) {
			clr := ui.Get.ProgressBarColor(ent)
			phc := ui.ProgressBarHoverColor{}

			phc.Hover.Empty = clr.Empty.Blend(color.Black, hoverColorDarkenFactor)
			phc.Hover.Border = clr.Border.Blend(color.Black, hoverColorDarkenFactor)

			phc.Normal.Empty = clr.Empty.Blend(color.Black, normalColorDarkenFactor)
			phc.Normal.Border = clr.Border.Blend(color.Black, normalColorDarkenFactor)

			phc.Clicked.Empty = clr.Empty
			phc.Clicked.Border = clr.Border

			phc.Disabled.Empty = clr.Empty.GrayScale().Blend(color.Black, normalColorDarkenFactor)
			phc.Disabled.Border = clr.Border.GrayScale().Blend(color.Black, normalColorDarkenFactor)

			if clr.Gradient.From.Equals(clr.Gradient.To) {
				phc.Hover.Solid = clr.Solid.Blend(color.Black, hoverColorDarkenFactor)
				phc.Normal.Solid = clr.Solid.Blend(color.Black, normalColorDarkenFactor)
				phc.Clicked.Solid = clr.Solid
				phc.Disabled.Solid = clr.Solid.GrayScale().Blend(color.Black, normalColorDarkenFactor)
			} else {
				phc.Hover.Gradient = color.Gradient{
					From:      clr.Gradient.From.Blend(color.Black, hoverColorDarkenFactor),
					To:        clr.Gradient.To.Blend(color.Black, hoverColorDarkenFactor),
					Direction: clr.Gradient.Direction,
				}
				phc.Normal.Gradient = color.Gradient{
					From:      clr.Gradient.From.Blend(color.Black, normalColorDarkenFactor),
					To:        clr.Gradient.To.Blend(color.Black, normalColorDarkenFactor),
					Direction: clr.Gradient.Direction,
				}
				phc.Clicked.Gradient = clr.Gradient
				phc.Disabled.Gradient = color.Gradient{
					From:      clr.Gradient.From.GrayScale().Blend(color.Black, normalColorDarkenFactor),
					To:        clr.Gradient.To.GrayScale().Blend(color.Black, normalColorDarkenFactor),
					Direction: clr.Gradient.Direction,
				}
			}
			ent.Set(phc)
		} else {
			uim.progressBarsColors(ent)
		}
	}
}

func (uim uiManager) progressBarsColors(ent *goecs.Entity) {
	bcl := ui.Get.ProgressBarHoverColor(ent)
	state := ui.Get.ControlState(ent)
	if state.Disabled {
		ent.Set(bcl.Disabled)
	} else if state.Clicked {
		ent.Set(bcl.Clicked)
	} else if state.Hover {
		ent.Set(bcl.Hover)
	} else {
		ent.Set(bcl.Normal)
	}
	if state.Focused {
		if ent.Contains(color.TYPE.Solid) {
			ac := color.Get.Solid(ent)
			clr := ui.Get.ProgressBarColor(ent)
			clr.Border = ac
			ent.Set(clr)
		}
	}
}

func (uim *uiManager) progressBarsMouseMove(world *goecs.World, mme events.MouseMoveEvent) {
	if uim.clicked != nil {
		if uim.clicked.Contains(ui.TYPE.ProgressBar) {
			state := ui.Get.ControlState(uim.clicked)
			if !state.Disabled {
				uim.calculateBarCurrent(world, uim.clicked, mme.Point)
			}
		}
		return
	}
	for it := world.Iterator(ui.TYPE.ProgressBar, ui.TYPE.ProgressBarHoverColor, geometry.TYPE.Point,
		shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		bar := ui.Get.ProgressBar(ent)
		state := ui.Get.ControlState(ent)
		if bar.Event != nil && !state.Disabled {
			pos := geometry.Get.Point(ent)
			box := shapes.Get.Box(ent)
			state.Hover = box.Contains(pos, mme.Point)
			uim.progressBarsColors(ent)
			ent.Set(bar)
			ent.Set(state)
		}
	}
	return
}

func (uim *uiManager) progressBarsMouseDown(world *goecs.World, mde events.MouseDownEvent) {
	for it := world.Iterator(ui.TYPE.ProgressBar, ui.TYPE.ProgressBarHoverColor, geometry.TYPE.Point,
		shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		bar := ui.Get.ProgressBar(ent)
		state := ui.Get.ControlState(ent)
		if state.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)
		if box.Contains(pos, mde.Point) {
			state.Clicked = true
			uim.clicked = ent
			uim.progressBarsColors(ent)
			ent.Set(bar)
			ent.Set(state)
		}
	}
}

func (uim uiManager) calculateBarCurrent(world *goecs.World, ent *goecs.Entity, mouse geometry.Point) {
	bar := ui.Get.ProgressBar(ent)
	previous := bar.Current
	pos := geometry.Get.Point(ent)
	box := shapes.Get.Box(ent)

	total := box.Size.Width * box.Scale
	max := pos.X + total
	cur := max - mouse.X
	per := 1 - (cur / total)

	diff := bar.Max - bar.Min
	bar.Current = diff * per
	if bar.Current < bar.Min {
		bar.Current = bar.Min
	}
	if bar.Current > bar.Max {
		bar.Current = bar.Max
	}

	if previous != bar.Current {
		ent.Set(bar)
		if bar.Event != nil {
			world.Signal(bar.Event)
		}
	}
}

func (uim *uiManager) progressBarsMouseUp(world *goecs.World, mue events.MouseUpEvent) {
	for it := world.Iterator(ui.TYPE.ProgressBar, ui.TYPE.ProgressBarHoverColor, geometry.TYPE.Point,
		shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		bar := ui.Get.ProgressBar(ent)
		state := ui.Get.ControlState(ent)
		if state.Clicked {
			state.Clicked = false
			uim.clicked = nil
			ent.Set(bar)
			ent.Set(state)

			uim.calculateBarCurrent(world, ent, mue.Point)

			if bar.Sound != "" {
				world.Signal(events.PlaySoundEvent{Name: bar.Sound, Volume: bar.Volume})
			}

			uim.progressBarsColors(ent)
		}
	}
}

func (uim uiManager) spriteButtons(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.SpriteButton); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.NotContains(ui.TYPE.ControlState) {
			ent.Set(ui.ControlState{})
		}
		sb := ui.Get.SpriteButton(ent)
		// add sprite if we have not one yet
		if ent.NotContains(sprite.TYPE) {
			spr := sprite.Sprite{
				Sheet: sb.Sheet,
				Name:  sb.Normal,
				Scale: sb.Scale,
			}
			ent.Set(spr)
		} else {
			uim.refreshSpriteButton(ent)
		}
	}
}

func (uim uiManager) refreshSpriteButton(ent *goecs.Entity) {
	spr := sprite.Get(ent)
	sbn := ui.Get.SpriteButton(ent)
	state := ui.Get.ControlState(ent)

	if state.Disabled {
		spr.Name = sbn.Disabled
	} else if state.Clicked {
		spr.Name = sbn.Clicked
	} else if state.Hover {
		spr.Name = sbn.Hover
	} else {
		spr.Name = sbn.Normal
	}

	if state.Focused && !state.Clicked {
		spr.Name = sbn.Focused
	}

	ent.Set(spr)
}

func (uim *uiManager) spriteButtonsMouseMove(world *goecs.World, mme events.MouseMoveEvent) {
	if uim.clicked != nil {
		return
	}
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		sbn := ui.Get.SpriteButton(ent)
		state := ui.Get.ControlState(ent)
		if state.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		spr := sprite.Get(ent)
		state.Hover = uim.cm.SpriteAtContains(spr, pos, mme.Point)
		ent.Set(sbn)
		ent.Set(state)
		uim.refreshSpriteButton(ent)
	}
}

func (uim *uiManager) spriteButtonsMouseDown(world *goecs.World, mde events.MouseDownEvent) {
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		sbn := ui.Get.SpriteButton(ent)
		state := ui.Get.ControlState(ent)
		if state.Disabled {
			continue
		}
		spr := sprite.Get(ent)
		pos := geometry.Get.Point(ent)

		if uim.cm.SpriteAtContains(spr, pos, mde.Point) {
			uim.clicked = ent
			state.Clicked = true
			ent.Set(sbn)
			ent.Set(state)
			uim.refreshSpriteButton(ent)
		}
	}
}

func (uim *uiManager) spriteButtonsMouseUp(world *goecs.World, _ events.MouseUpEvent) {
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		sbn := ui.Get.SpriteButton(ent)
		state := ui.Get.ControlState(ent)
		if state.Clicked {
			uim.clicked = nil
			state.Clicked = false
			ent.Set(sbn)
			ent.Set(state)
			uim.refreshSpriteButton(ent)
			if sbn.Sound != "" {
				world.Signal(events.PlaySoundEvent{Name: sbn.Sound, Volume: sbn.Volume})
			}
			world.Signal(sbn.Event)
		}
	}
}

func (uim *uiManager) focusControl(world *goecs.World, control *goecs.Entity) {
	for it := world.Iterator(ui.TYPE.ControlState); it != nil; it = it.Next() {
		ent := it.Value()
		state := ui.Get.ControlState(ent)
		state.Focused = ent.ID() == control.ID()
		if state.Focused {
			clr := effects.AlternateColor{
				From:  color.White,
				To:    color.White.Alpha(170),
				Time:  0.35,
				Delay: 0.15,
			}
			ent.Add(clr)
		} else {
			if ent.Contains(effects.TYPE.AlternateColorState) {
				ent.Remove(effects.TYPE.AlternateColorState)
			}
			if ent.Contains(effects.TYPE.AlternateColor) {
				ent.Remove(effects.TYPE.AlternateColor)
			}
		}
		ent.Set(state)
	}
	uim.focus = control
}

func (uim *uiManager) clearFocus(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.ControlState); it != nil; it = it.Next() {
		ent := it.Value()
		state := ui.Get.ControlState(ent)
		state.Focused = false
		if ent.Contains(effects.TYPE.AlternateColorState) {
			ent.Remove(effects.TYPE.AlternateColorState)
		}
		if ent.Contains(effects.TYPE.AlternateColor) {
			ent.Remove(effects.TYPE.AlternateColor)
		}
		ent.Set(state)
	}
	uim.focus = nil
}

func (uim *uiManager) handleFocus(world *goecs.World, key device.Key) {
	if key == device.KeyDown || key == device.KeyUp || key == device.KeyLeft || key == device.KeyRight {
		if uim.focus.Contains(ui.TYPE.ProgressBar) && (key == device.KeyLeft || key == device.KeyRight) {
			uim.moveProgressBarWithKeys(world, key)
		} else {
			uim.selectNextControl(world, key)
		}
	} else if key == device.KeySpace {
		uim.activateFocus(world)
	}
}

func (uim *uiManager) handleKeys(world *goecs.World, delta float32) {
	if uim.focus == nil || uim.lastKey == device.FirstKey {
		return
	}
	uim.keyDelay -= delta
	if uim.keyDelay <= 0 {
		uim.handleFocus(world, uim.lastKey)
		uim.keyDelay = keyWaitDelay
	}
}

func (uim *uiManager) keyDown(world *goecs.World, key device.Key) {
	if uim.focus == nil {
		return
	}
	uim.lastKey = key
	uim.keyDelay = keyWaitDelay
	uim.handleFocus(world, key)
}

func (uim *uiManager) keyUp(key device.Key) {
	if uim.lastKey == key {
		uim.lastKey = device.FirstKey
		uim.keyDelay = 0
	}
}

func (uim *uiManager) moveProgressBarWithKeys(world *goecs.World, key device.Key) {
	bar := ui.Get.ProgressBar(uim.focus)
	amt := (bar.Max - bar.Min) * 0.05
	if key == device.KeyLeft {
		bar.Current = float32(math.Max(float64(bar.Current-amt), float64(bar.Min)))
	} else {
		bar.Current = float32(math.Min(float64(bar.Current+amt), float64(bar.Max)))
	}
	if bar.Event != nil {
		world.Signal(bar.Event)
		if bar.Sound != "" {
			world.Signal(events.PlaySoundEvent{Name: bar.Sound, Volume: bar.Volume})
		}
	}
	uim.focus.Set(bar)
}

func (uim *uiManager) selectNextControl(world *goecs.World, key device.Key) {
	focusPos := geometry.Get.Point(uim.focus)
	focusRect := uim.getControlRect(uim.focus, focusPos)
	smallerDiff := float32(math.MaxFloat32)
	var possible *goecs.Entity
	for it := world.Iterator(ui.TYPE.ControlState); it != nil; it = it.Next() {
		ent := it.Value()
		state := ui.Get.ControlState(ent)
		if ent.ID() == uim.focus.ID() || ent.Contains(effects.TYPE.Hide) || state.Disabled {
			continue
		}
		possiblePos := geometry.Get.Point(ent)
		possibleRect := uim.getControlRect(ent, possiblePos)
		distance := uim.getControlDistance(focusRect, possibleRect, key)

		if distance < smallerDiff {
			smallerDiff = distance
			possible = ent
		}
	}
	if possible != nil {
		if possible.Contains(ui.TYPE.FlatButton) {
			fbt := ui.Get.FlatButton(possible)
			if fbt.Group != "" {
				possible = uim.handleGroupSelection(world, possible)
			}
		}
		uim.focusControl(world, possible)
	}
}

func (uim *uiManager) handleGroupSelection(world *goecs.World, possible *goecs.Entity) *goecs.Entity {
	match := possible
	fbt := ui.Get.FlatButton(possible)
	if uim.focus.Contains(ui.TYPE.FlatButton) {
		fbtFocus := ui.Get.FlatButton(uim.focus)
		if fbtFocus.Group != fbt.Group {
			match = uim.getSelectedInGroup(world, fbt.Group)
		}
	} else {
		match = uim.getSelectedInGroup(world, fbt.Group)
	}
	return match
}

func (uim uiManager) getControlRect(control *goecs.Entity, at geometry.Point) geometry.Rect {
	var rect geometry.Rect
	if control.Contains(sprite.TYPE) {
		spr := sprite.Get(control)
		rect = uim.cm.getSpriteRectAt(spr, at)
	} else {
		box := shapes.Get.Box(control)
		rect = box.GetReactAt(at)
	}
	return rect
}

func (uim uiManager) getControlDistance(rect1 geometry.Rect, rect2 geometry.Rect, key device.Key) float32 {
	var point1 geometry.Point
	var point2 geometry.Point

	switch key {
	case device.KeyUp:
		point1.X = rect1.From.X + (rect1.Size.Width / 2)
		point1.Y = rect1.From.Y

		point2.X = rect2.From.X + (rect2.Size.Width / 2)
		point2.Y = rect2.From.Y + (rect2.Size.Height)
		if point1.Y <= point2.Y {
			return math.MaxFloat32
		}
	case device.KeyDown:
		point1.X = rect1.From.X + (rect1.Size.Width / 2)
		point1.Y = rect1.From.Y + (rect1.Size.Height)

		point2.X = rect2.From.X + (rect2.Size.Width / 2)
		point2.Y = rect2.From.Y
		if point1.Y >= point2.Y {
			return math.MaxFloat32
		}
	case device.KeyLeft:
		point1.X = rect1.From.X
		point1.Y = rect1.From.Y + (rect1.Size.Height / 2)

		point2.X = rect2.From.X + (rect2.Size.Width)
		point2.Y = rect2.From.Y + (rect2.Size.Height / 2)

		if point1.X <= point2.X {
			return math.MaxFloat32
		}
	case device.KeyRight:
		point1.X = rect1.From.X + (rect1.Size.Width)
		point1.Y = rect1.From.Y + (rect1.Size.Height / 2)

		point2.X = rect2.From.X
		point2.Y = rect2.From.Y + (rect2.Size.Height / 2)

		if point1.X >= point2.X {
			return math.MaxFloat32
		}
	}
	return rect1.From.Distance(rect2.From)
}

func (uim *uiManager) getSelectedInGroup(world *goecs.World, group string) *goecs.Entity {
	for it := world.Iterator(ui.TYPE.FlatButton); it != nil; it = it.Next() {
		ent := it.Value()
		fbt := ui.Get.FlatButton(ent)
		if fbt.Group == group {
			state := ui.Get.ControlState(ent)
			if state.Checked {
				return ent
			}
		}
	}
	return nil
}

func (uim *uiManager) activateFocus(world *goecs.World) {
	if uim.focus.Contains(ui.TYPE.FlatButton) {
		state := ui.Get.ControlState(uim.focus)
		state.Clicked = true
		uim.focus.Set(state)
		uim.clicked = uim.focus
		uim.flatButtonsMouseUp(world, events.MouseUpEvent{})
	} else if uim.focus.Contains(ui.TYPE.SpriteButton) {
		state := ui.Get.ControlState(uim.focus)
		state.Clicked = true
		uim.focus.Set(state)
		uim.clicked = uim.focus
		uim.spriteButtonsMouseUp(world, events.MouseUpEvent{})
	}
}

func (uim *uiManager) gamepadDown(world *goecs.World, _ int32, button device.GamepadButton) {
	switch button {
	case device.GamepadUp:
		uim.keyDown(world, device.KeyUp)
	case device.GamepadDown:
		uim.keyDown(world, device.KeyDown)
	case device.GamepadLeft:
		uim.keyDown(world, device.KeyLeft)
	case device.GamepadRight:
		uim.keyDown(world, device.KeyRight)
	case device.GamepadButton3:
		uim.keyDown(world, device.KeySpace)
	}
}

func (uim *uiManager) gamepadUp(_ int32, button device.GamepadButton) {
	switch button {
	case device.GamepadUp:
		uim.keyUp(device.KeyUp)
	case device.GamepadDown:
		uim.keyUp(device.KeyDown)
	case device.GamepadLeft:
		uim.keyUp(device.KeyLeft)
	case device.GamepadRight:
		uim.keyUp(device.KeyRight)
	case device.GamepadButton3:
		uim.keyUp(device.KeySpace)
	}
}

func (uim *uiManager) gamepadStickMove(world *goecs.World, _ int32, _ device.GamepadStick, movement geometry.Point) {
	if movement.Y == -1 && uim.lastKey != device.KeyUp {
		uim.keyDown(world, device.KeyUp)
	} else if movement.Y > -1 && uim.lastKey == device.KeyUp {
		uim.keyUp(device.KeyUp)
	}
	if movement.Y == 1 && uim.lastKey != device.KeyDown {
		uim.keyDown(world, device.KeyDown)
	} else if movement.Y < 1 && uim.lastKey == device.KeyDown {
		uim.keyUp(device.KeyDown)
	}
	if movement.X == -1 && uim.lastKey != device.KeyLeft {
		uim.keyDown(world, device.KeyLeft)
	} else if movement.X > -1 && uim.lastKey == device.KeyLeft {
		uim.keyUp(device.KeyLeft)
	}
	if movement.X == 1 && uim.lastKey != device.KeyRight {
		uim.keyDown(world, device.KeyRight)
	} else if movement.Y < 1 && uim.lastKey == device.KeyRight {
		uim.keyUp(device.KeyRight)
	}
}

// UI returns a managers.WithSystemAndListener that handle ui components
func UI(dm DeviceManager, cm *CollisionManager) WithSystemAndListener {
	return &uiManager{
		dm:       dm,
		cm:       cm,
		clicked:  nil,
		focus:    nil,
		lastKey:  device.FirstKey,
		keyDelay: 0,
	}
}
