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
	"fmt"
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/internal/components"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/internal/storage"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/tiled"
	"strconv"
)

type tiledSystem struct {
	ds storage.Storage
}

const (
	rightDown = "right-down"
)

func (ts tiledSystem) Update(wld *world.World, _ float32) (err error) {
	for it := wld.Iterator(tiled.TYPE.Map, geometry.TYPE.Point); it.HasNext(); {
		ent := it.Value()
		tiledMap := tiled.Get.Map(ent)
		if ent.NotContains(tiled.TYPE.MapState) {
			depth := int32(render.DefaultLayer)
			if ent.Contains(effects.TYPE.Layer) {
				depth = effects.Get.Layer(ent).Depth
			}

			if err = ts.addSpriteFromTiledMap(wld, tiledMap, depth); err == nil {
				pos := geometry.Get.Point(ent)
				ent.Set(tiled.MapState{Pos: pos})
			}
		}
	}
	return
}

func (ts tiledSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}

func (ts *tiledSystem) GetTilePosition(x, y int, def components.TiledMapDef) geometry.Point {
	return geometry.Point{
		X: float32(x * def.Data.TileWidth),
		Y: float32(y * def.Data.TileHeight),
	}
}

func (ts tiledSystem) addSpriteFromTiledMap(wld *world.World, tiledMap tiled.Map, depth int32) (err error) {
	if mapDef, err := ts.ds.GetTiledMapDef(tiledMap.Name); err == nil {
		if !(mapDef.Data.RenderOrder == "" || mapDef.Data.RenderOrder == rightDown) {
			return fmt.Errorf("unsupported tiled render order : got %q, want %q", mapDef.Data.RenderOrder, rightDown)
		}

		tl := len(mapDef.Data.Layers)
		for ln, l := range mapDef.Data.Layers {
			if !l.Visible {
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
					pos := ts.GetTilePosition(x, y, mapDef)
					wld.Add(entity.New(
						sprite.Sprite{
							Sheet: tiledMap.Name,
							Name:  sprName,
							Scale: 1,
						},
						pos,
						effects.Layer{Depth: depth + int32(tl-ln)},
					))
					i++
				}
			}
		}
	}

	return
}

// TiledSystem returns a world.System that handle tiled maps
func TiledSystem(ds storage.Storage) world.System {
	return tiledSystem{
		ds: ds,
	}
}
