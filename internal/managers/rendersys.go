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
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/internal/storage"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/ui"
	"sort"
)

type renderingManager struct {
	rdr render.Render
	ds  storage.Storage
}

var noTint = color.White

func (rdm renderingManager) renderSprite(ent *entity.Entity) error {
	spr := sprite.Get(ent)
	pos := geometry.Get.Point(ent)

	var tint color.Solid
	if ent.Contains(color.TYPE.Solid) {
		tint = color.Get.Solid(ent)
	} else {
		tint = noTint
	}

	if def, err := rdm.ds.GetSpriteDef(spr.Sheet, spr.Name); err == nil {
		if err := rdm.rdr.DrawSprite(def, spr, pos, tint); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (rdm renderingManager) renderShape(ent *entity.Entity) error {
	pos := geometry.Get.Point(ent)
	box := shapes.Get.Box(ent)
	if ent.Contains(color.TYPE.Solid) {
		clr := color.Get.Solid(ent)
		rdm.rdr.DrawSolidBox(pos, box, clr)
	} else if ent.Contains(color.TYPE.Gradient) {
		gra := color.Get.Gradient(ent)
		rdm.rdr.DrawGradientBox(pos, box, gra)
	}
	return nil
}

func (rdm renderingManager) renderFlatButton(ent *entity.Entity) error {
	pos := geometry.Get.Point(ent)
	box := shapes.Get.Box(ent)
	fb := ui.Get.FlatButton(ent)

	if fb.Shadow.Width > 0 || fb.Shadow.Height > 0 {
		shadowPos := geometry.Point{
			X: pos.X + fb.Shadow.Width,
			Y: pos.Y + fb.Shadow.Height,
		}
		sc := color.Gray.Alpha(127)
		rdm.rdr.DrawSolidBox(shadowPos, box, sc)
	}

	// draw the box
	if ent.Contains(color.TYPE.Solid) {
		clr := color.Get.Solid(ent)
		rdm.rdr.DrawSolidBox(pos, box, clr)
	} else if ent.Contains(color.TYPE.Gradient) {
		gra := color.Get.Gradient(ent)
		rdm.rdr.DrawGradientBox(pos, box, gra)
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

		if ftd, err := rdm.ds.GetFontDef(textCmp.Font); err == nil {
			// draw the text
			rdm.rdr.DrawText(ftd, textCmp, textPos, colorCmp)
		} else {
			return err
		}
	}

	return nil
}

func (rdm renderingManager) renderText(v *entity.Entity) error {
	textCmp := ui.Get.Text(v)
	posCmp := geometry.Get.Point(v)
	colorCmp := color.Get.Solid(v)

	if ftd, err := rdm.ds.GetFontDef(textCmp.Font); err == nil {
		// draw the text
		rdm.rdr.DrawText(ftd, textCmp, posCmp, colorCmp)
	} else {
		return err
	}

	return nil
}

func (rdm renderingManager) isRenderable(ent *entity.Entity) bool {
	return ent.Contains(geometry.TYPE.Point) &&
		(ent.Contains(sprite.TYPE) || ent.Contains(ui.TYPE.Text) || ent.Contains(shapes.TYPE.Box) ||
			ent.Contains(ui.TYPE.FlatButton))
}

func (rdm renderingManager) getSortedByLayers(world *world.World) []*entity.Entity {
	entities := make([]*entity.Entity, world.Size())
	i := 0
	for it := world.Iterator(); it != nil; it = it.Next() {
		e := it.Value()
		if rdm.isRenderable(e) {
			entities[i] = e
			i++
		}
	}
	entities = entities[:i]

	sort.Slice(entities, func(i, j int) bool {
		first := entities[i]
		second := entities[j]

		firstDepth := render.DefaultLayer
		if first.Contains(effects.TYPE.Layer) {
			firstDepth = effects.Get.Layer(first).Depth
		}

		secondDepth := render.DefaultLayer
		if second.Contains(effects.TYPE.Layer) {
			secondDepth = effects.Get.Layer(second).Depth
		}
		if firstDepth == secondDepth {
			return first.ID() < second.ID()
		}
		return firstDepth > secondDepth
	})
	return entities
}

func (rdm renderingManager) System(world *world.World, _ float32) error {
	for _, v := range rdm.getSortedByLayers(world) {
		if v.Contains(sprite.TYPE) {
			if err := rdm.renderSprite(v); err != nil {
				return err
			}
		} else if v.Contains(ui.TYPE.FlatButton) {
			if err := rdm.renderFlatButton(v); err != nil {
				return err
			}
		} else if v.Contains(shapes.TYPE.Box) {
			if err := rdm.renderShape(v); err != nil {
				return err
			}
		} else if v.Contains(ui.TYPE.Text, color.TYPE.Solid) {
			if err := rdm.renderText(v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Rendering returns a managers.WithSystem that will handle rendering
func Rendering(rdr render.Render, ds storage.Storage) WithSystem {
	return &renderingManager{
		rdr: rdr,
		ds:  ds,
	}
}
