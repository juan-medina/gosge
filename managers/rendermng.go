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
)

type renderingManager struct {
	dm DeviceManager
	sm *StorageManager
}

var noTint = color.White

func (rdm renderingManager) renderSprite(ent *goecs.Entity) error {
	spr := sprite.Get(ent)
	pos := geometry.Get.Point(ent)

	var tint color.Solid
	if ent.Contains(color.TYPE.Solid) {
		tint = color.Get.Solid(ent)
	} else {
		tint = noTint
	}

	if def, err := rdm.sm.GetSpriteDef(spr.Sheet, spr.Name); err == nil {
		if err := rdm.dm.DrawSprite(def, spr, pos, tint); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (rdm renderingManager) renderBox(ent *goecs.Entity) error {
	pos := geometry.Get.Point(ent)
	box := shapes.Get.Box(ent)
	clr := color.Get.Solid(ent)
	rdm.dm.DrawBox(pos, box, clr)
	return nil
}

func (rdm renderingManager) renderSolidBox(ent *goecs.Entity) error {
	pos := geometry.Get.Point(ent)
	box := shapes.Get.SolidBox(ent)
	if ent.Contains(color.TYPE.Solid) {
		clr := color.Get.Solid(ent)
		rdm.dm.DrawSolidBox(pos, box, clr)
	} else if ent.Contains(color.TYPE.Gradient) {
		gra := color.Get.Gradient(ent)
		rdm.dm.DrawGradientBox(pos, box, gra)
	}
	return nil
}

func (rdm renderingManager) renderLine(ent *goecs.Entity) error {
	pos := geometry.Get.Point(ent)
	line := shapes.Get.Line(ent)
	clr := color.White

	if ent.Contains(color.TYPE.Solid) {
		clr = color.Get.Solid(ent)
	}

	rdm.dm.DrawLine(pos, line.To, line.Thickness, clr)

	return nil
}

func (rdm renderingManager) renderFlatButton(ent *goecs.Entity) error {
	pos := geometry.Get.Point(ent)
	box := shapes.Get.SolidBox(ent)
	fb := ui.Get.FlatButton(ent)

	if fb.Shadow.Width > 0 || fb.Shadow.Height > 0 {
		shadowPos := geometry.Point{
			X: pos.X + fb.Shadow.Width,
			Y: pos.Y + fb.Shadow.Height,
		}
		sc := color.Gray.Alpha(127)
		rdm.dm.DrawSolidBox(shadowPos, box, sc)
	}

	// draw the box
	if ent.Contains(color.TYPE.Solid) {
		clr := color.Get.Solid(ent)
		rdm.dm.DrawSolidBox(pos, box, clr)
	} else if ent.Contains(color.TYPE.Gradient) {
		gra := color.Get.Gradient(ent)
		rdm.dm.DrawGradientBox(pos, box, gra)
	}

	// if contains a text component
	if ent.Contains(ui.TYPE.Text) {
		// get the text component
		textCmp := ui.Get.Text(ent)

		// calculate position on the center of the button
		textPos := geometry.Point{
			X: pos.X + ((box.Size.Width / 2) * box.Scale),
			Y: pos.Y + ((box.Size.Height / 2) * box.Scale),
		}

		// default color white
		colorCmp := color.White

		// if the button has a solid color
		if ent.Contains(color.TYPE.Solid) {
			// text color is inverse
			colorCmp = color.Get.Solid(ent).Inverse()
			// if is has a gradient color
		} else if ent.Contains(color.TYPE.Gradient) {
			// get the gradient
			gra := color.Get.Gradient(ent)
			// text color if the inverse of the from
			colorCmp = gra.From.Inverse()
		}

		if ftd, err := rdm.sm.GetFontDef(textCmp.Font); err == nil {
			// draw the text
			rdm.dm.DrawText(ftd, textCmp, textPos, colorCmp)
		} else {
			return err
		}
	}

	return nil
}

func (rdm renderingManager) renderProgressBar(v *goecs.Entity) error {
	box := shapes.Get.Box(v)
	pos := geometry.Get.Point(v)
	pro := ui.Get.ProgressBar(v)
	clr := ui.Get.ProgressBarColor(v)

	sb := shapes.SolidBox{
		Size:  box.Size,
		Scale: box.Scale,
	}

	// draw shadow
	if pro.Shadow.Width > 0 || pro.Shadow.Height > 0 {
		shadowPos := geometry.Point{
			X: pos.X + pro.Shadow.Width,
			Y: pos.Y + pro.Shadow.Height,
		}

		sc := color.Gray.Alpha(127)
		rdm.dm.DrawSolidBox(shadowPos, sb, sc)
	}

	// draw background
	rdm.dm.DrawSolidBox(pos, sb, clr.Empty)

	total := pro.Max - pro.Min
	per := pro.Current / total

	// if we need to fill anything
	if per > 0 {
		// if we need scissor
		if per < 1 {
			size := box.Size.Scale(box.Scale)
			size.Width = size.Width * per
			rdm.dm.BeginScissor(pos, size)
		}

		// draw fill
		if clr.Gradient.From.Equals(clr.Gradient.To) {
			rdm.dm.DrawSolidBox(pos, sb, clr.Solid)
		} else {
			rdm.dm.DrawGradientBox(pos, sb, clr.Gradient)
		}

		// if we need scissor
		if per < 1 {
			rdm.dm.EndScissor()
		}
	}

	// draw border
	if box.Thickness > 0 {
		rdm.dm.DrawBox(pos, box, clr.Border)
	}

	return nil
}

func (rdm renderingManager) renderText(v *goecs.Entity) error {
	textCmp := ui.Get.Text(v)
	posCmp := geometry.Get.Point(v)
	colorCmp := color.Get.Solid(v)

	if ftd, err := rdm.sm.GetFontDef(textCmp.Font); err == nil {
		// draw the text
		rdm.dm.DrawText(ftd, textCmp, posCmp, colorCmp)
	} else {
		return err
	}

	return nil
}

func (rdm renderingManager) isRenderable(ent *goecs.Entity) bool {
	return ent.Contains(geometry.TYPE.Point) &&
		(ent.Contains(sprite.TYPE) || ent.Contains(ui.TYPE.Text) || ent.Contains(shapes.TYPE.Box) ||
			ent.Contains(shapes.TYPE.SolidBox) || ent.Contains(ui.TYPE.FlatButton) ||
			ent.Contains(ui.TYPE.ProgressBar) || ent.Contains(shapes.TYPE.Line))
}

func (rdm renderingManager) sortRenderable(first, second *goecs.Entity) bool {
	if !rdm.isRenderable(first) {
		return false
	} else if !rdm.isRenderable(second) {
		return true
	} else {
		firstDepth := DefaultLayer
		if first.Contains(effects.TYPE.Layer) {
			firstDepth = effects.Get.Layer(first).Depth
		}

		secondDepth := DefaultLayer
		if second.Contains(effects.TYPE.Layer) {
			secondDepth = effects.Get.Layer(second).Depth
		}
		if firstDepth == secondDepth {
			return first.ID() < second.ID()
		}
		return firstDepth > secondDepth
	}
}

func (rdm renderingManager) System(world *goecs.World, _ float32) error {
	// sort by renderable in-place
	world.Sort(rdm.sortRenderable)

	// go trough all the world
	for it := world.Iterator(); it != nil; it = it.Next() {
		v := it.Value()
		if !rdm.isRenderable(v) {
			break // since is sorted by renderable we don't have nothing more to render
		}
		if v.Contains(sprite.TYPE) {
			if err := rdm.renderSprite(v); err != nil {
				return err
			}
		} else if v.Contains(ui.TYPE.FlatButton) {
			if err := rdm.renderFlatButton(v); err != nil {
				return err
			}
		} else if v.Contains(ui.TYPE.ProgressBar) {
			if err := rdm.renderProgressBar(v); err != nil {
				return err
			}
		} else if v.Contains(shapes.TYPE.Box) {
			if err := rdm.renderBox(v); err != nil {
				return err
			}
		} else if v.Contains(shapes.TYPE.SolidBox) {
			if err := rdm.renderSolidBox(v); err != nil {
				return err
			}
		} else if v.Contains(ui.TYPE.Text, color.TYPE.Solid) {
			if err := rdm.renderText(v); err != nil {
				return err
			}
		} else if v.Contains(shapes.TYPE.Line) {
			if err := rdm.renderLine(v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Rendering returns a managers.WithSystem that will handle rendering
func Rendering(dm DeviceManager, sm *StorageManager) WithSystem {
	return &renderingManager{
		dm: dm,
		sm: sm,
	}
}
