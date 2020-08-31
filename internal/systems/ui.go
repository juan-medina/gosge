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
	// normalColorDarkenFactor is how much dark the color will be on normal (not hover)
	normalColorDarkenFactor = 0.25
)

type uiSystem struct{}

func (uis uiSystem) Update(wld *world.World, _ float32) error {
	uis.flatButtons(wld)
	return nil
}

func (uis uiSystem) Notify(wld *world.World, event interface{}, _ float32) error {
	switch v := event.(type) {
	case events.MouseMoveEvent:
		uis.flatButtonsMouseMove(wld, v)
	case events.MouseUpEvent:
		if err := uis.flatButtonsMouseUp(wld, v); err != nil {
			return err
		}
	}
	return nil
}

func (uis uiSystem) flatButtons(wld *world.World) {
	for it := wld.Iterator(ui.TYPE.FlatButton); it.HasNext(); {
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

func (uis uiSystem) flatButtonsMouseMove(wld *world.World, mme events.MouseMoveEvent) {
	for it := wld.Iterator(ui.TYPE.FlatButton, ui.TYPE.ButtonHoverColors); it.HasNext(); {
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

func (uis uiSystem) flatButtonsMouseUp(wld *world.World, mue events.MouseUpEvent) error {
	for it := wld.Iterator(ui.TYPE.FlatButton); it.HasNext(); {
		ent := it.Value()
		btn := ui.Get.FlatButton(ent)

		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)

		if box.Contains(pos, mue.Point) {
			return wld.Notify(btn.Event)
		}
	}
	return nil
}

// UISystem returns a world.System that handle ui components
func UISystem() world.System {
	return &uiSystem{}
}
