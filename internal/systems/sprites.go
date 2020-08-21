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
	"github.com/juan-medina/gosge/pkg/components/position"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type spriteSheetJSON struct {
	Atlas struct {
		ImagePath string `json:"imagePath"`
	} `json:"atlas"`
	Sprites []struct {
		NameID   string `json:"nameId"`
		Position struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"position"`
		TrimRec struct {
			X      int `json:"x"`
			Y      int `json:"y"`
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"trimRec"`
	} `json:"sprites"`
}

type spriteSheet map[string]components.SpriteDef

type spriteRenderingSystem struct {
	sheets map[string]spriteSheet
}

var noTint = color.White

// SpriteRendering is a world.System for rendering sprite.Sprite
type SpriteRendering interface {
	//Update the world.World
	Update(world *world.World, delta float64) error
	//Notify is trigger when new events are received
	Notify(world *world.World, event interface{}, delta float64) error
	// LoadSpriteSheet preloads a sprite.Sprite sheet
	LoadSpriteSheet(fileName string) error
}

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

		if sheet, ok := s.sheets[spr.Sheet]; ok {
			if def, ok := sheet[spr.Name]; ok {
				if err := render.DrawSprite(def, spr, pos, tint); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("can not find sprite %q in sheet %q", spr.Name, spr.Sheet)
			}
		} else {
			return fmt.Errorf("can not find sprite sheet %q", spr.Sheet)
		}

	}
	return nil
}

func (s spriteRenderingSystem) Notify(_ *world.World, _ interface{}, _ float64) error {
	return nil
}

func (s *spriteRenderingSystem) loadSpriteSheetFile(fileName string, sheet *spriteSheetJSON) (err error) {
	var jsonFile *os.File
	if jsonFile, err = os.Open(fileName); err == nil {
		dir := filepath.Dir(fileName)
		//goland:noinspection GoUnhandledErrorResult
		defer jsonFile.Close()
		var bytes []byte
		if bytes, err = ioutil.ReadAll(jsonFile); err == nil {
			if err = json.Unmarshal(bytes, &sheet); err == nil {
				st := make(spriteSheet, 0)
				s.sheets[fileName] = st
				texturePath := path.Join(dir, sheet.Atlas.ImagePath)
				if err = render.LoadTexture(texturePath); err == nil {
					for _, spr := range sheet.Sprites {
						st[spr.NameID] = components.SpriteDef{
							Texture: texturePath,
							X:       spr.TrimRec.X + spr.Position.X,
							Y:       spr.TrimRec.Y + spr.Position.Y,
							Width:   spr.TrimRec.Width,
							Height:  spr.TrimRec.Height,
						}
					}
				}
			}
		}
	}
	return
}

func (s *spriteRenderingSystem) LoadSpriteSheet(fileName string) (err error) {
	var ss spriteSheetJSON
	if err = s.loadSpriteSheetFile(fileName, &ss); err == nil {
	}
	return
}

// SpriteRenderingSystem returns a world.System that will handle sprite.Sprite rendering
func SpriteRenderingSystem() SpriteRendering {
	return &spriteRenderingSystem{
		sheets: make(map[string]spriteSheet, 0),
	}
}
