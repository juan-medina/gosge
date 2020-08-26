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
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/ui"
	"github.com/juan-medina/gosge/pkg/events"
)

const (
	// normalColorDarken is how much dark the color will be on normal (not hover)
	normalColorDarken = 0.25
)

type uiSystem struct{}

func (uis uiSystem) Update(wld *world.World, _ float32) error {
	for it := wld.Iterator(ui.TYPE.FlatButton); it.HasNext(); {
		ent := it.Value()
		btn := ui.Get.FlatButton(ent)

		// calculate hover color if is not done yet
		if !btn.Hover.Set {
			if ent.Contains(color.TYPE.Solid) {
				btn.Hover.Solid = color.Get.Solid(ent)
			} else if ent.Contains(color.TYPE.Gradient) {
				btn.Hover.Gradient = color.Get.Gradient(ent)
			}
			btn.Hover.Set = true
			ent.Set(btn)
		}

		// calculate normal (no hover) color if is not done yet
		if !btn.Normal.Set {
			if ent.Contains(color.TYPE.Solid) {
				btn.Normal.Solid = color.Get.Solid(ent).Blend(color.Black, normalColorDarken)
			} else if ent.Contains(color.TYPE.Gradient) {
				gra := color.Get.Gradient(ent)
				btn.Normal.Gradient = color.Gradient{
					From:      gra.From.Blend(color.Black, normalColorDarken),
					To:        gra.To.Blend(color.Black, normalColorDarken),
					Direction: gra.Direction,
				}
			}
			btn.Normal.Set = true
			ent.Set(btn)
		}
	}
	return nil
}

func (uis uiSystem) Notify(wld *world.World, event interface{}, _ float32) error {
	switch v := event.(type) {
	case events.MouseMoveEvent:
		for it := wld.Iterator(ui.TYPE.FlatButton); it.HasNext(); {
			ent := it.Value()
			btn := ui.Get.FlatButton(ent)

			pos := geometry.Get.Point(ent)
			box := shapes.Get.Box(ent)

			var clr interface{} = color.White

			if box.Contains(pos, v.Point) {
				if ent.Contains(color.TYPE.Solid) {
					clr = btn.Hover.Solid
				} else if ent.Contains(color.TYPE.Gradient) {
					clr = btn.Hover.Gradient
				}
			} else {
				if ent.Contains(color.TYPE.Solid) {
					clr = btn.Normal.Solid
				} else if ent.Contains(color.TYPE.Gradient) {
					clr = btn.Normal.Gradient
				}
			}
			if clr != nil {
				ent.Set(clr)
			}
		}
	case events.MouseUpEvent:
		for it := wld.Iterator(ui.TYPE.FlatButton); it.HasNext(); {
			ent := it.Value()
			btn := ui.Get.FlatButton(ent)

			pos := geometry.Get.Point(ent)
			box := shapes.Get.Box(ent)

			if box.Contains(pos, v.Point) {
				return wld.Notify(btn.Event)
			}
		}
	}
	return nil
}

// UISystem returns a world.System that handle ui components
func UISystem() world.System {
	return &uiSystem{}
}
