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
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
)

const (
	normalColorDarkenFactor = 0.25 // normalColorDarkenFactor is how much dark the color will be on normal
	hoverColorDarkenFactor  = 0.10 // hoverColorDarkenFactor is how much dark the color will be on hover
)

type uiManager struct {
	dm      DeviceManager
	cm      *CollisionManager
	clicked *goecs.Entity
}

func (uim *uiManager) System(world *goecs.World, _ float32) error {
	uim.flatButtons(world)
	uim.progressBars(world)
	uim.spriteButtons(world)
	return nil
}

func (uim *uiManager) Listener(world *goecs.World, event interface{}, _ float32) error {
	switch v := event.(type) {
	case events.MouseMoveEvent:
		uim.flatButtonsMouseMove(world, v)
		if err := uim.progressBarsMouseMove(world, v); err != nil {
			return err
		}
		uim.spriteButtonsMouseMove(world, v)
	case events.MouseDownEvent:
		if err := uim.flatButtonsMouseDown(world, v); err != nil {
			return err
		}
		if err := uim.progressBarsMouseDown(world, v); err != nil {
			return err
		}
		if err := uim.spriteButtonsMouseDown(world, v); err != nil {
			return err
		}
	case events.MouseUpEvent:
		if err := uim.flatButtonsMouseUp(world, v); err != nil {
			return err
		}
		if err := uim.progressBarsMouseUp(world, v); err != nil {
			return err
		}
		if err := uim.spriteButtonsMouseUp(world, v); err != nil {
			return err
		}
	}
	return nil
}

func (uim *uiManager) flatButtons(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.FlatButton); it != nil; it = it.Next() {
		ent := it.Value()
		// calculate state if is not done yet
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
					To:        clr.Gradient.To.Blend(color.Black, normalColorDarkenFactor),
					Direction: clr.Gradient.Direction,
				}
				bcl.Normal.Gradient = color.Gradient{
					From:      clr.Gradient.From.Blend(color.Black, hoverColorDarkenFactor),
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
			btn := ui.Get.FlatButton(ent)
			uim.flatButtonsColors(ent, btn)
		}
	}
}

func (uim uiManager) flatButtonsColors(ent *goecs.Entity, btn ui.FlatButton) {
	bcl := ui.Get.ButtonHoverColors(ent)
	if btn.State.Disabled {
		ent.Set(bcl.Disabled)
	} else if btn.State.Clicked {
		ent.Set(bcl.Clicked)
	} else if btn.State.Hover {
		ent.Set(bcl.Hover)
	} else {
		ent.Set(bcl.Normal)
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
		if btn.State.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)
		btn.State.Hover = box.Contains(pos, mme.Point)
		uim.flatButtonsColors(ent, btn)
		ent.Set(btn)
	}
}

func (uim *uiManager) flatButtonsMouseDown(world *goecs.World, mde events.MouseDownEvent) error {
	for it := world.Iterator(ui.TYPE.FlatButton, geometry.TYPE.Point, shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		btn := ui.Get.FlatButton(ent)
		if btn.State.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)
		if box.Contains(pos, mde.Point) {
			btn.State.Clicked = true
			uim.clicked = ent
		}
		uim.flatButtonsColors(ent, btn)
		ent.Set(btn)
	}
	return nil
}

func (uim *uiManager) flatButtonsMouseUp(world *goecs.World, _ events.MouseUpEvent) error {
	for it := world.Iterator(ui.TYPE.FlatButton, geometry.TYPE.Point, shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		btn := ui.Get.FlatButton(ent)
		if btn.State.Clicked {
			if btn.CheckBox {
				btn.State.Checked = !btn.State.Checked
			}
			btn.State.Clicked = false
			if btn.Sound != "" {
				if err := world.Signal(events.PlaySoundEvent{Name: btn.Sound, Volume: btn.Volume}); err != nil {
					return err
				}
			}
			if err := world.Signal(btn.Event); err != nil {
				return err
			}
			uim.clicked = nil
		}
		uim.flatButtonsColors(ent, btn)
		ent.Set(btn)
	}
	return nil
}

func (uim uiManager) progressBars(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.ProgressBar); it != nil; it = it.Next() {
		ent := it.Value()
		// calculate state if is not done yet
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
			bar := ui.Get.ProgressBar(ent)
			uim.progressBarsColors(ent, bar)
		}
	}
}

func (uim uiManager) progressBarsColors(ent *goecs.Entity, bar ui.ProgressBar) {
	bcl := ui.Get.ProgressBarHoverColor(ent)
	if bar.State.Disabled {
		ent.Set(bcl.Disabled)
	} else if bar.State.Clicked {
		ent.Set(bcl.Clicked)
	} else if bar.State.Hover {
		ent.Set(bcl.Hover)
	} else {
		ent.Set(bcl.Normal)
	}
}

func (uim *uiManager) progressBarsMouseMove(world *goecs.World, mme events.MouseMoveEvent) error {
	if uim.clicked != nil {
		if uim.clicked.Contains(ui.TYPE.ProgressBar) {
			bar := ui.Get.ProgressBar(uim.clicked)
			if !bar.State.Disabled {
				if err := uim.calculateBarCurrent(world, uim.clicked, mme.Point); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for it := world.Iterator(ui.TYPE.ProgressBar, ui.TYPE.ProgressBarHoverColor, geometry.TYPE.Point,
		shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		bar := ui.Get.ProgressBar(ent)
		if bar.Event != nil && !bar.State.Disabled {
			pos := geometry.Get.Point(ent)
			box := shapes.Get.Box(ent)
			bar.State.Hover = box.Contains(pos, mme.Point)
			uim.progressBarsColors(ent, bar)
			ent.Set(bar)
		}
	}
	return nil
}

func (uim *uiManager) progressBarsMouseDown(world *goecs.World, mde events.MouseDownEvent) error {
	for it := world.Iterator(ui.TYPE.ProgressBar, ui.TYPE.ProgressBarHoverColor, geometry.TYPE.Point,
		shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		bar := ui.Get.ProgressBar(ent)
		if bar.State.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		box := shapes.Get.Box(ent)
		if box.Contains(pos, mde.Point) {
			bar.State.Clicked = true
			uim.clicked = ent
			uim.progressBarsColors(ent, bar)
			ent.Set(bar)
		}
	}
	return nil
}

func (uim uiManager) calculateBarCurrent(world *goecs.World, ent *goecs.Entity, mouse geometry.Point) error {
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
			if err := world.Signal(bar.Event); err != nil {
				return err
			}
		}
	}

	return nil
}

func (uim *uiManager) progressBarsMouseUp(world *goecs.World, mue events.MouseUpEvent) error {
	for it := world.Iterator(ui.TYPE.ProgressBar, ui.TYPE.ProgressBarHoverColor, geometry.TYPE.Point,
		shapes.TYPE.Box); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		bar := ui.Get.ProgressBar(ent)

		if bar.State.Clicked {
			bar.State.Clicked = false
			uim.clicked = nil
			ent.Set(bar)

			if err := uim.calculateBarCurrent(world, ent, mue.Point); err != nil {
				return err
			}

			if bar.Sound != "" {
				if err := world.Signal(events.PlaySoundEvent{Name: bar.Sound, Volume: bar.Volume}); err != nil {
					return err
				}
			}

			uim.progressBarsColors(ent, bar)
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
			uim.refreshSpriteButton(ent)
		}
	}
}

func (uim uiManager) refreshSpriteButton(ent *goecs.Entity) {
	spr := sprite.Get(ent)
	sbn := ui.Get.SpriteButton(ent)

	if sbn.State.Disabled {
		spr.Name = sbn.Disabled
	} else if sbn.State.Clicked {
		spr.Name = sbn.Clicked
	} else if sbn.State.Hover {
		spr.Name = sbn.Hover
	} else {
		spr.Name = sbn.Normal
	}

	ent.Set(sbn)
	ent.Set(spr)
}

func (uim *uiManager) spriteButtonsMouseMove(world *goecs.World, mme events.MouseMoveEvent) {
	if uim.clicked != nil {
		return
	}
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		sbn := ui.Get.SpriteButton(ent)
		if sbn.State.Disabled {
			continue
		}
		pos := geometry.Get.Point(ent)
		spr := sprite.Get(ent)
		sbn.State.Hover = uim.cm.SpriteAtContains(spr, pos, mme.Point)
		ent.Set(sbn)
		uim.refreshSpriteButton(ent)
	}
}

func (uim *uiManager) spriteButtonsMouseDown(world *goecs.World, mde events.MouseDownEvent) error {
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		sbn := ui.Get.SpriteButton(ent)
		if sbn.State.Disabled {
			continue
		}
		spr := sprite.Get(ent)
		pos := geometry.Get.Point(ent)

		if uim.cm.SpriteAtContains(spr, pos, mde.Point) {
			uim.clicked = ent
			sbn.State.Clicked = true
			ent.Set(sbn)
			uim.refreshSpriteButton(ent)
		}
	}
	return nil
}

func (uim *uiManager) spriteButtonsMouseUp(world *goecs.World, _ events.MouseUpEvent) error {
	for it := world.Iterator(ui.TYPE.SpriteButton, sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			continue
		}
		sbn := ui.Get.SpriteButton(ent)
		if sbn.State.Clicked {
			uim.clicked = nil
			sbn.State.Clicked = false
			ent.Set(sbn)
			uim.refreshSpriteButton(ent)
			if sbn.Sound != "" {
				if err := world.Signal(events.PlaySoundEvent{Name: sbn.Sound, Volume: sbn.Volume}); err != nil {
					return err
				}
			}
			return world.Signal(sbn.Event)
		}
	}
	return nil
}

// UI returns a managers.WithSystemAndListener that handle ui components
func UI(dm DeviceManager, cm *CollisionManager) WithSystemAndListener {
	return &uiManager{dm: dm, cm: cm, clicked: nil}
}
