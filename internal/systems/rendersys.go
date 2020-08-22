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
	"github.com/juan-medina/gosge/internal/store"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
	"sort"
)

type renderingSystem struct {
	rdr render.Render
	ss  store.SpriteStorage
}

var noTint = color.White

func (rs renderingSystem) renderSprite(ent *entity.Entity) error {
	spr := sprite.Get(ent)
	pos := geometry.Get.Position(ent)

	var tint color.Solid
	if ent.Contains(color.TYPE.Solid) {
		tint = color.Get.Solid(ent)
	} else {
		tint = noTint
	}

	if def, err := rs.ss.GetSpriteDef(spr.Sheet, spr.Name); err == nil {
		if err := rs.rdr.DrawSprite(def, spr, pos, tint); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (rs renderingSystem) renderShape(ent *entity.Entity) error {
	pos := geometry.Get.Position(ent)
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

func (rs renderingSystem) renderText(v *entity.Entity) error {
	textCmp := text.Get(v)
	posCmp := geometry.Get.Position(v)
	colorCmp := color.Get.Solid(v)

	rs.rdr.DrawText(textCmp, posCmp, colorCmp)
	return nil
}

func (rs renderingSystem) getSortedByLayers(world *world.World) []*entity.Entity {
	entities := world.Entities(geometry.TYPE.Position)

	sort.Slice(entities, func(i, j int) bool {
		first := entities[i]
		second := entities[j]

		firstDepth := int32(-1000)
		if first.Contains(effects.TYPE.Layer) {
			firstDepth = effects.Get.Layer(first).Depth
		}

		secondDepth := int32(-1000)
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
		} else if v.Contains(shapes.TYPE.Box) {
			if err := rs.renderShape(v); err != nil {
				return err
			}
		} else if v.Contains(text.TYPE, color.TYPE.Solid) {
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
func RenderingSystem(rdr render.Render, ss store.SpriteStorage) world.System {
	return &renderingSystem{
		rdr: rdr,
		ss:  ss,
	}
}
