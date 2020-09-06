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
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/sprite"
)

// CollisionManager is a manager for manage collisions
type CollisionManager struct {
	sm *StorageManager
}

// getSpriteRect return a geometry.Rect for a given sprite.Sprite at a geometry.Point
func (cm CollisionManager) getSpriteRect(spr sprite.Sprite, at geometry.Point) geometry.Rect {
	def, _ := cm.sm.GetSpriteDef(spr.Sheet, spr.Name)
	size := def.Origin.Size.Scale(spr.Scale)

	return geometry.Rect{
		From: geometry.Point{
			X: at.X - (size.Width * def.Pivot.X),
			Y: at.Y - (size.Height * def.Pivot.Y),
		},
		Size: size,
	}
}

// SpriteAtContains indicates if a sprite.Sprite at a given geometry.Point contains a geometry.Point
func (cm CollisionManager) SpriteAtContains(spr sprite.Sprite, at geometry.Point, point geometry.Point) bool {
	return cm.getSpriteRect(spr, at).IsPointInRect(point)
}

// Collisions returns a CollisionManager
func Collisions(sm *StorageManager) *CollisionManager {
	return &CollisionManager{sm: sm}
}
