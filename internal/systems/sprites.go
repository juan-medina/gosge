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
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/position"
	"github.com/juan-medina/gosge/pkg/components/sprite"
)

type spriteRenderingSystem struct{}

var noTint = color.White

func (s spriteRenderingSystem) Update(world *world.World, _ float64) error {
	for _, v := range world.Entities(sprite.TYPE, position.TYPE) {
		spr := sprite.Get(v)
		pos := position.Get(v)

		var tint color.Color
		if v.Contains(color.TYPE) {
			tint = color.Get(v)
		} else {
			tint = noTint
		}

		if err := render.DrawSprite(spr, pos, tint); err != nil {
			return err
		}
	}
	return nil
}

func (s spriteRenderingSystem) Notify(_ *world.World, _ interface{}, _ float64) error {
	return nil
}

// SpriteRenderingSystem returns a world.System that will handle sprite.Sprite rendering
func SpriteRenderingSystem() world.System {
	return spriteRenderingSystem{}
}
