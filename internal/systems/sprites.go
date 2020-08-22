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
	"encoding/json"
	"fmt"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/internal/components"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type spriteSheetData struct {
	Frames []struct {
		Filename string `json:"filename"`
		Frame    struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
			W float32 `json:"w"`
			H float32 `json:"h"`
		} `json:"frame"`
		Pivot struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		} `json:"pivot"`
	} `json:"frames"`
	Meta struct {
		Image string `json:"image"`
	} `json:"meta"`
}

type spriteSheet map[string]components.SpriteDef

type spriteRenderingSystem struct {
	sheets map[string]spriteSheet
	rdr    render.Render
}

var noTint = color.White

// SpriteRendering is a world.System for rendering sprite.Sprite
type SpriteRendering interface {
	//Update the world.World
	Update(world *world.World, delta float32) error
	//Notify is trigger when new events are received
	Notify(world *world.World, event interface{}, delta float32) error
	// LoadSpriteSheet preloads a sprite.Sprite sheet
	LoadSpriteSheet(fileName string) error
	//GetSpriteSize returns the geometry.Size of a given sprite
	GetSpriteSize(sheet string, name string) (geometry.Size, error)
}

func (srs spriteRenderingSystem) Update(world *world.World, _ float32) error {
	for _, v := range world.Entities(sprite.TYPE, geometry.TYPE.Position) {
		spr := sprite.Get(v)
		pos := geometry.Get.Position(v)

		var tint color.Solid
		if v.Contains(color.TYPE.Solid) {
			tint = color.Get.Solid(v)
		} else {
			tint = noTint
		}

		if def, err := srs.getSpriteDef(spr.Sheet, spr.Name); err == nil {
			if err := srs.rdr.DrawSprite(def, spr, pos, tint); err != nil {
				return err
			}
		} else {
			return err
		}

	}
	return nil
}

func (srs spriteRenderingSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}

func (srs *spriteRenderingSystem) handleSheet(data spriteSheetData, name string) (err error) {
	st := make(spriteSheet, 0)
	srs.sheets[name] = st
	dir := filepath.Dir(name)
	texturePath := path.Join(dir, data.Meta.Image)
	if err = srs.rdr.LoadTexture(texturePath); err == nil {
		for _, spr := range data.Frames {
			st[spr.Filename] = components.SpriteDef{
				Texture: texturePath,
				Origin: geometry.Rect{
					From: geometry.Point{
						X: spr.Frame.X,
						Y: spr.Frame.Y,
					},
					Size: geometry.Size{
						Width:  spr.Frame.W,
						Height: spr.Frame.H,
					},
				},
				Pivot: geometry.Point{
					X: spr.Pivot.X,
					Y: spr.Pivot.Y,
				},
			}
		}
	}

	return
}

func (srs *spriteRenderingSystem) LoadSpriteSheet(fileName string) (err error) {
	data := spriteSheetData{}
	var jsonFile *os.File
	if jsonFile, err = os.Open(fileName); err == nil {
		//goland:noinspection GoUnhandledErrorResult
		defer jsonFile.Close()
		var bytes []byte
		if bytes, err = ioutil.ReadAll(jsonFile); err == nil {
			if err = json.Unmarshal(bytes, &data); err == nil {
				return srs.handleSheet(data, fileName)
			}
		}
	}
	return
}

func (srs spriteRenderingSystem) getSpriteDef(sheet string, name string) (components.SpriteDef, error) {
	if sh, ok := srs.sheets[sheet]; ok {
		if def, ok := sh[name]; ok {
			return def, nil
		}
		return components.SpriteDef{}, fmt.Errorf("can not find sprite %q in sheet %q", name, sheet)
	}
	return components.SpriteDef{}, fmt.Errorf("can not find sprite sheet %q", sheet)
}

func (srs spriteRenderingSystem) GetSpriteSize(sheet string, name string) (geometry.Size, error) {
	def, err := srs.getSpriteDef(sheet, name)
	//goland:noinspection GoNilness
	return def.Origin.Size, err
}

// SpriteRenderingSystem returns a world.System that will handle sprite.Sprite rendering
func SpriteRenderingSystem(rdr render.Render) SpriteRendering {
	return &spriteRenderingSystem{
		sheets: make(map[string]spriteSheet, 0),
		rdr:    rdr,
	}
}
