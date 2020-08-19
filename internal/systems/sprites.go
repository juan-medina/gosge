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
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/pkg/components"
)

type spriteRenderingSystem struct{}

var noTint = components.WhiteColor

func (s spriteRenderingSystem) Update(world *world.World, _ float64) error {
	for _, v := range world.Entities(components.SpriteType, components.PosType) {
		sprite := v.Get(components.SpriteType).(components.Sprite)
		pos := v.Get(components.PosType).(components.Pos)

		var tint components.RGBAColor
		if v.Contains(components.RGBAColorType) {
			tint = v.Get(components.RGBAColorType).(components.RGBAColor)
		} else {
			tint = noTint
		}

		if err := render.DrawSprite(sprite, pos, tint); err != nil {
			return err
		}
	}
	return nil
}

func (s spriteRenderingSystem) Notify(_ *world.World, _ interface{}, _ float64) error {
	return nil
}

func SpriteRenderingSystem() world.System {
	return spriteRenderingSystem{}
}
