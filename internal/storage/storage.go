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

// Package storage x
package storage

import (
	"encoding/json"
	"fmt"
	"github.com/juan-medina/gosge/internal/components"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/pkg/components/geometry"
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

type dataStorage struct {
	sheets   map[string]spriteSheet
	textures map[string]components.TextureDef
	rdr      render.Render
}

// Storage is a storage for our game data
type Storage interface {
	// LoadSpriteSheet preloads a sprite.Sprite sheet
	LoadSpriteSheet(fileName string) error
	//GetSpriteSize returns the geometry.Size of a given sprite
	GetSpriteSize(sheet string, name string) (geometry.Size, error)
	//GetSpriteDef returns the components.SpriteDef for an sprite
	GetSpriteDef(sheet string, name string) (components.SpriteDef, error)
	//Clear all loaded data
	Clear()
}

func (ds *dataStorage) loadTexture(name string) (def components.TextureDef, err error) {
	if _, ok := ds.textures[name]; !ok {
		if texture, err := ds.rdr.LoadTexture(name); err != nil {
			ds.textures[name] = texture
		} else {
			return texture, err
		}
	}
	return ds.textures[name], nil
}

func (ds *dataStorage) handleSheet(data spriteSheetData, name string) (err error) {
	var texture components.TextureDef
	st := make(spriteSheet, 0)
	ds.sheets[name] = st
	dir := filepath.Dir(name)
	texturePath := path.Join(dir, data.Meta.Image)
	if texture, err = ds.rdr.LoadTexture(texturePath); err == nil {
		for _, spr := range data.Frames {
			st[spr.Filename] = components.SpriteDef{
				Texture: texture,
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

func (ds *dataStorage) LoadSpriteSheet(fileName string) (err error) {
	data := spriteSheetData{}
	var jsonFile *os.File
	if jsonFile, err = os.Open(fileName); err == nil {
		//goland:noinspection GoUnhandledErrorResult
		defer jsonFile.Close()
		var bytes []byte
		if bytes, err = ioutil.ReadAll(jsonFile); err == nil {
			if err = json.Unmarshal(bytes, &data); err == nil {
				return ds.handleSheet(data, fileName)
			}
		}
	}
	return
}

func (ds dataStorage) GetSpriteDef(sheet string, name string) (components.SpriteDef, error) {
	if sh, ok := ds.sheets[sheet]; ok {
		if def, ok := sh[name]; ok {
			return def, nil
		}
		return components.SpriteDef{}, fmt.Errorf("can not find sprite %q in sheet %q", name, sheet)
	}
	return components.SpriteDef{}, fmt.Errorf("can not find sprite sheet %q", sheet)
}

func (ds dataStorage) GetSpriteSize(sheet string, name string) (geometry.Size, error) {
	def, err := ds.GetSpriteDef(sheet, name)
	//goland:noinspection GoNilness
	return def.Origin.Size, err
}

func (ds *dataStorage) Clear() {
	ds.sheets = make(map[string]spriteSheet, 0)

	for _, v := range ds.textures {
		ds.rdr.UnloadTexture(v)
	}
	ds.textures = make(map[string]components.TextureDef, 0)
}

// New returns a new storage.Storage
func New(rdr render.Render) Storage {
	return &dataStorage{
		sheets:   make(map[string]spriteSheet, 0),
		textures: make(map[string]components.TextureDef, 0),
		rdr:      rdr,
	}
}
