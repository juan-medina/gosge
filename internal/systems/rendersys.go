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

type renderingSystem struct {
	rdr render.Render
	ds  storage.Storage
}

var noTint = color.White

const (
	deepestLayer = 1000000
)

func (rs renderingSystem) renderSprite(ent *entity.Entity) error {
	spr := sprite.Get(ent)
	pos := geometry.Get.Point(ent)

	var tint color.Solid
	if ent.Contains(color.TYPE.Solid) {
		tint = color.Get.Solid(ent)
	} else {
		tint = noTint
	}

	if def, err := rs.ds.GetSpriteDef(spr.Sheet, spr.Name); err == nil {
		if err := rs.rdr.DrawSprite(def, spr, pos, tint); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (rs renderingSystem) renderShape(ent *entity.Entity) error {
	pos := geometry.Get.Point(ent)
	box := shapes.Get.Box(ent)
	if ent.Contains(color.TYPE.Solid) {
		clr := color.Get.Solid(ent)
		rs.rdr.DrawSolidBox(pos, box, clr)
	} else if ent.Contains(color.TYPE.Gradient) {
		gra := color.Get.Gradient(ent)
		rs.rdr.DrawGradientBox(pos, box, gra)
	}
	return nil
}

func (rs renderingSystem) renderFlatButton(ent *entity.Entity) error {
	pos := geometry.Get.Point(ent)
	box := shapes.Get.Box(ent)
	fb := ui.Get.FlatButton(ent)

	if fb.Shadow.Width > 0 || fb.Shadow.Height > 0 {
		shadowPos := geometry.Point{
			X: pos.X + fb.Shadow.Width,
			Y: pos.Y + fb.Shadow.Height,
		}
		sc := color.Gray.Alpha(127)
		rs.rdr.DrawSolidBox(shadowPos, box, sc)
	}

	// draw the box
	if ent.Contains(color.TYPE.Solid) {
		clr := color.Get.Solid(ent)
		rs.rdr.DrawSolidBox(pos, box, clr)
	} else if ent.Contains(color.TYPE.Gradient) {
		gra := color.Get.Gradient(ent)
		rs.rdr.DrawGradientBox(pos, box, gra)
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

		if ftd, err := rs.ds.GetFontDef(textCmp.Font); err == nil {
			// draw the text
			rs.rdr.DrawText(ftd, textCmp, textPos, colorCmp)
		} else {
			return err
		}
	}

	return nil
}

func (rs renderingSystem) renderText(v *entity.Entity) error {
	textCmp := ui.Get.Text(v)
	posCmp := geometry.Get.Point(v)
	colorCmp := color.Get.Solid(v)

	if ftd, err := rs.ds.GetFontDef(textCmp.Font); err == nil {
		// draw the text
		rs.rdr.DrawText(ftd, textCmp, posCmp, colorCmp)
	} else {
		return err
	}

	return nil
}

func (rs renderingSystem) isRenderable(ent *entity.Entity) bool {
	return ent.Contains(geometry.TYPE.Point) &&
		(ent.Contains(sprite.TYPE) || ent.Contains(ui.TYPE.Text) || ent.Contains(shapes.TYPE.Box) ||
			ent.Contains(ui.TYPE.FlatButton))
}

func (rs renderingSystem) getSortedByLayers(world *world.World) []*entity.Entity {
	entities := make([]*entity.Entity, world.Size())
	i := 0
	for it := world.Iterator(); it.HasNext(); {
		e := it.Value()
		if rs.isRenderable(e) {
			entities[i] = e
			i++
		}
	}
	entities = entities[:i]

	sort.Slice(entities, func(i, j int) bool {
		first := entities[i]
		second := entities[j]

		firstDepth := int32(deepestLayer)
		if first.Contains(effects.TYPE.Layer) {
			firstDepth = effects.Get.Layer(first).Depth
		}

		secondDepth := int32(deepestLayer)
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

func (rs renderingSystem) Update(world *world.World, _ float32) error {
	for _, v := range rs.getSortedByLayers(world) {
		if v.Contains(sprite.TYPE) {
			if err := rs.renderSprite(v); err != nil {
				return err
			}
		} else if v.Contains(ui.TYPE.FlatButton) {
			if err := rs.renderFlatButton(v); err != nil {
				return err
			}
		} else if v.Contains(shapes.TYPE.Box) {
			if err := rs.renderShape(v); err != nil {
				return err
			}
		} else if v.Contains(ui.TYPE.Text, color.TYPE.Solid) {
			if err := rs.renderText(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rs renderingSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}

// RenderingSystem returns a world.System that will handle rendering
func RenderingSystem(rdr render.Render, ds storage.Storage) world.System {
	return &renderingSystem{
		rdr: rdr,
		ds:  ds,
	}
}
