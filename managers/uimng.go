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
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
)

const (
	// normalColorDarkenFactor is how much dark the color will be on normal (not hover)
	normalColorDarkenFactor = 0.25
)

type uiManager struct {
	dm DeviceManager
	sm *StorageManager
}

func (uim uiManager) System(world *goecs.World, _ float32) error {
	uim.flatButtons(world)
	uim.spriteButtons(world)
	return nil
}

func (uim uiManager) Listener(world *goecs.World, event interface{}, _ float32) error {
	switch v := event.(type) {
	case events.MouseMoveEvent:
		uim.flatButtonsMouseMove(world, v)
		uim.spriteButtonsMouseMove(world, v)
	case events.MouseUpEvent:
		if err := uim.flatButtonsMouseUp(world, v); err != nil {
			return err
		}
		if err := uim.spriteButtonsMouseUp(world, v); err != nil {
			return err
		}
	}
	return nil
}

func (uim uiManager) flatButtons(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.FlatButton); it != nil; it = it.Next() {
		ent := it.Value()
		// calculate state if is not done yet
		if ent.NotContains(ui.TYPE.ButtonHoverColors) {
			bcl := ui.ButtonHoverColors{}
			if ent.Contains(color.TYPE.Solid) {
				bcl.Hover.Solid = color.Get.Solid(ent)
				bcl.Normal.Solid = color.Get.Solid(ent).Blend(color.Black, normalColorDarkenFactor)
			} else if ent.Contains(color.TYPE.Gradient) {
				gra := color.Get.Gradient(ent)
				bcl.Hover.Gradient = gra
				bcl.Normal.Gradient = color.Gradient{
					From:      gra.From.Blend(color.Black, normalColorDarkenFactor),
					To:        gra.To.Blend(color.Black, normalColorDarkenFactor),
					Direction: gra.Direction,
				}
			}
			ent.Set(bcl)
		}
	}
}

func (uim uiManager) flatButtonsMouseMove(world *goecs.World, mme events.MouseMoveEvent) {
	for it := world.Iterator(ui.TYPE.FlatButton, ui.TYPE.ButtonHoverColors, geometry.TYPE.Point, shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		bcl := ui.Get.ButtonHoverColors(ent)
		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)

		var clr interface{} = color.White

		if box.Contains(pos, mme.Point) {
			if ent.Contains(color.TYPE.Solid) {
				clr = bcl.Hover.Solid
			} else if ent.Contains(color.TYPE.Gradient) {
				clr = bcl.Hover.Gradient
			}
		} else {
			if ent.Contains(color.TYPE.Solid) {
				clr = bcl.Normal.Solid
			} else if ent.Contains(color.TYPE.Gradient) {
				clr = bcl.Normal.Gradient
			}
		}
		ent.Set(clr)
	}
}

func (uim uiManager) flatButtonsMouseUp(world *goecs.World, mue events.MouseUpEvent) error {
	for it := world.Iterator(ui.TYPE.FlatButton, geometry.TYPE.Point, shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		btn := ui.Get.FlatButton(ent)

		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)

		if box.Contains(pos, mue.Point) {
			return world.Signal(btn.Event)
		}
	}
	return nil
}

func (uim uiManager) spriteButtons(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.SpriteButton); it != nil; it = it.Next() {
		ent := it.Value()
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
			spr := sprite.Get(ent)
			// if the sprite button has change
			if !(spr.Name == sb.Normal || spr.Name == sb.Hover) {
				uim.refreshButtonOnPoint(ent, uim.dm.GetMousePoint())
			}
		}
	}
}

func (uim uiManager) refreshButtonOnPoint(ent *goecs.Entity, pnt geometry.Point) {
	spr := sprite.Get(ent)
	pos := geometry.Get.Point(ent)
	sbn := ui.Get.SpriteButton(ent)

	if uim.sm.SpriteAtContains(spr, pos, pnt) {
		spr.Name = sbn.Hover
	} else {
		spr.Name = sbn.Normal
	}
	ent.Set(spr)
}

func (uim uiManager) spriteButtonsMouseMove(world *goecs.World, mme events.MouseMoveEvent) {
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		uim.refreshButtonOnPoint(ent, mme.Point)
	}
}

func (uim uiManager) spriteButtonsMouseUp(world *goecs.World, mue events.MouseUpEvent) error {
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		spr := sprite.Get(ent)
		pos := geometry.Get.Point(ent)
		sbn := ui.Get.SpriteButton(ent)

		if uim.sm.SpriteAtContains(spr, pos, mue.Point) {
			if sbn.Sound != "" {
				if err := world.Signal(events.PlaySoundEvent{Name: sbn.Sound}); err != nil {
					return err
				}
			}
			return world.Signal(sbn.Event)
		}
	}
	return nil
}

// UI returns a managers.WithSystemAndListener that handle ui components
func UI(dm DeviceManager, sm *StorageManager) WithSystemAndListener {
	return &uiManager{dm: dm, sm: sm}
}
