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
	"fmt"
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/tiled"
	"github.com/juan-medina/gosge/internal/components"
	"strconv"
)

type tiledManager struct {
	sm StorageManager
}

const (
	rightDown = "right-down"
)

func (tm tiledManager) System(world *goecs.World, _ float32) (err error) {
	for it := world.Iterator(tiled.TYPE.Map, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		tiledMap := tiled.Get.Map(ent)
		pos := geometry.Get.Point(ent)
		if ent.NotContains(tiled.TYPE.MapState) {
			depth := DefaultLayer
			if ent.Contains(effects.TYPE.Layer) {
				depth = effects.Get.Layer(ent).Depth
			}

			if err = tm.addSpriteFromTiledMap(world, tiledMap, depth, pos); err == nil {
				pos := geometry.Get.Point(ent)
				ent.Set(tiled.MapState{Position: pos, Scale: tiledMap.Scale})
			}
		} else {
			state := tiled.Get.MapState(ent)
			if state.Position.X != pos.X || state.Position.Y != pos.Y || state.Scale != tiledMap.Scale {
				// if the state need to change
				diff := geometry.Point{
					X: state.Position.X - pos.X,
					Y: state.Position.Y - pos.Y,
				}
				tm.updateSprites(world, tiledMap, diff)
				state.Position = pos
				state.Scale = tiledMap.Scale
				ent.Set(state)
			}
		}
	}
	return
}

func (tm *tiledManager) GetTilePosition(x, y int, def components.TiledMapDef) geometry.Point {
	return geometry.Point{
		X: float32(x * def.Data.TileWidth),
		Y: float32(y * def.Data.TileHeight),
	}
}

func (tm tiledManager) addSpriteFromTiledMap(world *goecs.World, tiledMap tiled.Map, depth float32, mapPos geometry.Point) (err error) {
	if mapDef, err := tm.sm.GetTiledMapDef(tiledMap.Name); err == nil {
		if !(mapDef.Data.RenderOrder == "" || mapDef.Data.RenderOrder == rightDown) {
			return fmt.Errorf("unsupported tiled render order : got %q, want %q", mapDef.Data.RenderOrder, rightDown)
		}

		ld := depth
		for _, l := range mapDef.Data.Layers {
			if !l.Visible {
				ld--
				continue
			}
			var xs, xe, xi, ys, ye, yi int
			xs = 0
			xe = mapDef.Data.Width
			xi = 1
			ys = 0
			ye = mapDef.Data.Height
			yi = 1

			i := 0
			for y := ys; y*yi < ye; y = y + yi {
				for x := xs; x*xi < xe; x = x + xi {
					if l.Tiles[i].IsNil() {
						i++
						continue
					}
					sprName := strconv.Itoa(int(l.Tiles[i].ID))
					pos := tm.GetTilePosition(x, y, mapDef)
					pos.X = (pos.X * tiledMap.Scale) + mapPos.X
					pos.Y = (pos.Y * tiledMap.Scale) + mapPos.Y
					world.AddEntity(
						sprite.Sprite{
							Sheet: tiledMap.Name,
							Name:  sprName,
							Scale: tiledMap.Scale,
							FlipX: l.Tiles[i].HorizontalFlip,
							FlipY: l.Tiles[i].VerticalFlip,
						},
						pos,
						effects.Layer{Depth: ld},
					)
					i++
				}
			}
			ld--
		}
	}

	return
}

func (tm tiledManager) updateSprites(world *goecs.World, tiledMap tiled.Map, diff geometry.Point) {
	for it := world.Iterator(sprite.TYPE, geometry.TYPE.Point); it != nil; it = it.Next() {
		ent := it.Value()
		spr := sprite.Get(ent)
		if spr.Sheet == tiledMap.Name {
			pos := geometry.Get.Point(ent)
			pos.X += diff.X
			pos.Y -= diff.Y
			ent.Set(pos)
		}
	}
}

// TiledMaps returns a managers.WithSystem that handle tiled maps
func TiledMaps(sm StorageManager) WithSystem {
	return tiledManager{
		sm: sm,
	}
}
