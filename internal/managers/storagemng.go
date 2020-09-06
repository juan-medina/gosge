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
	"encoding/json"
	"fmt"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/internal/components"
	"github.com/lafriks/go-tiled"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

// StorageManager is a storage for our game data
type StorageManager interface {
	// LoadSpriteSheet preloads a sprite.Sprite sheet
	LoadSpriteSheet(name string) error
	// LoadSingleSprite preloads a sprite.Sprite in a sheet
	LoadSingleSprite(sheet string, name string, pivot geometry.Point) error
	//GetSpriteSize returns the geometry.Size of a given sprite
	GetSpriteSize(sheet string, name string) (geometry.Size, error)
	//GetSpriteDef returns the components.SpriteDef for an sprite
	GetSpriteDef(sheet string, name string) (components.SpriteDef, error)
	// LoadFont preloads a font
	LoadFont(name string) error
	//GetFontDef returns the components.FontDef for a font
	GetFontDef(name string) (components.FontDef, error)
	//LoadMusic preload a music stream
	LoadMusic(name string) error
	//GetMusicDef returns the components.MusicDef for a music stream
	GetMusicDef(name string) (components.MusicDef, error)
	//LoadSound preload a sound wave
	LoadSound(name string) error
	//GetSoundDef returns the components.SoundDef for a wave sound
	GetSoundDef(name string) (components.SoundDef, error)
	// LoadTiledMap preload a tiled map
	LoadTiledMap(name string) error
	//GetTiledMapDef returns the components.TiledMapDef for a tiled map
	GetTiledMapDef(name string) (components.TiledMapDef, error)
	//Clear all loaded data
	Clear()
	// SpriteAtContains indicates if a sprite.Sprite at a given geometry.Point contains a geometry.Point
	SpriteAtContains(spr sprite.Sprite, at geometry.Point, point geometry.Point) bool
}

type spriteSheetData struct {
	Frames []struct {
		Name  string `json:"filename"`
		Frame struct {
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

type storageManager struct {
	sheets    map[string]spriteSheet
	textures  map[string]components.TextureDef
	fonts     map[string]components.FontDef
	musics    map[string]components.MusicDef
	sounds    map[string]components.SoundDef
	tiledMaps map[string]components.TiledMapDef
	dm        DeviceManager
}

func (sm *storageManager) LoadTiledMap(name string) (err error) {
	var tiledMap components.TiledMapDef
	if _, ok := sm.tiledMaps[name]; !ok {
		if tiledMap, err = sm.loadTileMap(name); err == nil {
			sm.tiledMaps[name] = tiledMap
		}
	}
	return
}

func (sm *storageManager) GetTiledMapDef(name string) (components.TiledMapDef, error) {
	if _, ok := sm.tiledMaps[name]; ok {
		return sm.tiledMaps[name], nil
	}
	return components.TiledMapDef{}, fmt.Errorf("can not find tiled map %q", name)
}

func (sm storageManager) loadTileMap(name string) (result components.TiledMapDef, err error) {
	var tiledMap *tiled.Map
	if tiledMap, err = tiled.LoadFromFile(name); err == nil {

		result.Cols = int32(tiledMap.Width)
		result.Rows = int32(tiledMap.Height)
		result.TileSize = geometry.Size{
			Width:  float32(tiledMap.TileWidth),
			Height: float32(tiledMap.TileHeight),
		}
		result.Size = geometry.Size{
			Width:  float32(result.Cols) * result.TileSize.Width,
			Height: float32(result.Rows) * result.TileSize.Height,
		}

		result.Data = tiledMap
		dir := filepath.Dir(name)
		var texture components.TextureDef
		for _, ts := range tiledMap.Tilesets {
			texturePath := path.Join(dir, ts.Image.Source)
			if texture, err = sm.dm.LoadTexture(texturePath); err != nil {
				return
			}
			st := make(spriteSheet, 0)
			sm.sheets[name] = st

			tilesetTileCount := ts.TileCount
			tilesetColumns := ts.Columns
			margin := ts.Margin
			spacing := ts.Spacing

			if tilesetColumns == 0 {
				tilesetColumns = ts.Image.Width / (ts.TileWidth + spacing)
			}
			if tilesetTileCount == 0 {
				tilesetTileCount = (ts.Image.Height / (ts.TileHeight + spacing)) * tilesetColumns
			}

			for i := ts.FirstGID; i < ts.FirstGID+uint32(tilesetTileCount); i++ {
				x := int(i-ts.FirstGID) % tilesetColumns
				y := int(i-ts.FirstGID) / tilesetColumns
				xOffset := int(x)*spacing + margin
				yOffset := int(y)*spacing + margin
				origin := geometry.Rect{
					From: geometry.Point{
						X: float32(x*ts.TileWidth + xOffset),
						Y: float32(y*ts.TileHeight + yOffset),
					},
					Size: geometry.Size{
						Width:  float32(ts.TileWidth),
						Height: float32(ts.TileHeight),
					},
				}

				sprName := strconv.Itoa(int(i - 1))
				st[sprName] = components.SpriteDef{
					Texture: texture,
					Origin:  origin,
					Pivot: geometry.Point{
						X: 0,
						Y: 0,
					},
				}
			}
		}
	}
	return
}

func (sm *storageManager) LoadSound(name string) (err error) {
	var sound components.SoundDef
	if _, ok := sm.sounds[name]; !ok {
		if sound, err = sm.dm.LoadSound(name); err == nil {
			sm.sounds[name] = sound
		}
	}
	return err
}

func (sm *storageManager) GetSoundDef(name string) (components.SoundDef, error) {
	if _, ok := sm.sounds[name]; ok {
		return sm.sounds[name], nil
	}
	return components.SoundDef{}, fmt.Errorf("can not find sound %q", name)
}

func (sm *storageManager) GetMusicDef(name string) (components.MusicDef, error) {
	if _, ok := sm.musics[name]; ok {
		return sm.musics[name], nil
	}
	return components.MusicDef{}, fmt.Errorf("can not find music %q", name)
}

func (sm *storageManager) LoadMusic(name string) (err error) {
	var music components.MusicDef
	if _, ok := sm.musics[name]; !ok {
		if music, err = sm.dm.LoadMusic(name); err == nil {
			sm.musics[name] = music
		}
	}
	return err
}

func (sm *storageManager) LoadFont(name string) (err error) {
	var font components.FontDef
	if _, ok := sm.fonts[name]; !ok {
		if font, err = sm.dm.LoadFont(name); err == nil {
			sm.fonts[name] = font
		}
	}
	return
}

func (sm *storageManager) GetFontDef(name string) (components.FontDef, error) {
	if _, ok := sm.fonts[name]; ok {
		return sm.fonts[name], nil
	}
	return components.FontDef{}, fmt.Errorf("can not find font %q", name)
}

func (sm *storageManager) loadTexture(name string) (def components.TextureDef, err error) {
	if _, ok := sm.textures[name]; !ok {
		if texture, err := sm.dm.LoadTexture(name); err == nil {
			sm.textures[name] = texture
		} else {
			return texture, err
		}
	}
	return sm.textures[name], nil
}

func (sm *storageManager) handleSheet(data spriteSheetData, name string) (err error) {
	var texture components.TextureDef
	st := make(spriteSheet, 0)
	sm.sheets[name] = st
	dir := filepath.Dir(name)
	texturePath := path.Join(dir, data.Meta.Image)
	if texture, err = sm.dm.LoadTexture(texturePath); err == nil {
		for _, spr := range data.Frames {
			st[spr.Name] = components.SpriteDef{
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

func (sm *storageManager) LoadSpriteSheet(name string) (err error) {
	data := spriteSheetData{}
	var jsonFile *os.File
	if jsonFile, err = os.Open(name); err == nil {
		//goland:noinspection GoUnhandledErrorResult
		defer jsonFile.Close()
		var bytes []byte
		if bytes, err = ioutil.ReadAll(jsonFile); err == nil {
			if err = json.Unmarshal(bytes, &data); err == nil {
				return sm.handleSheet(data, name)
			}
		}
	}
	return
}

func (sm *storageManager) LoadSingleSprite(sheet string, name string, pivot geometry.Point) error {
	if texture, err := sm.loadTexture(name); err == nil {
		if _, ok := sm.sheets[sheet]; !ok {
			st := make(spriteSheet, 0)
			sm.sheets[sheet] = st
		}
		sm.sheets[sheet][name] = components.SpriteDef{
			Texture: texture,
			Origin: geometry.Rect{
				From: geometry.Point{
					X: 0,
					Y: 0,
				},
				Size: texture.Size,
			},
			Pivot: pivot,
		}
	} else {
		return err
	}
	return nil
}

func (sm storageManager) GetSpriteDef(sheet string, name string) (components.SpriteDef, error) {
	if sh, ok := sm.sheets[sheet]; ok {
		if def, ok := sh[name]; ok {
			return def, nil
		}
		return components.SpriteDef{}, fmt.Errorf("can not find sprite %q in sheet %q", name, sheet)
	}
	return components.SpriteDef{}, fmt.Errorf("can not find sprite sheet %q", sheet)
}

func (sm storageManager) GetSpriteSize(sheet string, name string) (geometry.Size, error) {
	def, err := sm.GetSpriteDef(sheet, name)
	//goland:noinspection GoNilness
	return def.Origin.Size, err
}

func (sm *storageManager) Clear() {
	sm.sheets = make(map[string]spriteSheet, 0)

	for _, v := range sm.textures {
		sm.dm.UnloadTexture(v)
	}
	sm.textures = make(map[string]components.TextureDef, 0)

	for _, v := range sm.fonts {
		sm.dm.UnloadFont(v)
	}
	sm.fonts = make(map[string]components.FontDef, 0)

	for _, v := range sm.musics {
		sm.dm.UnloadMusic(v)
	}
	sm.musics = make(map[string]components.MusicDef, 0)

	for _, v := range sm.sounds {
		sm.dm.UnloadSound(v)
	}
	sm.sounds = make(map[string]components.SoundDef, 0)
	sm.tiledMaps = make(map[string]components.TiledMapDef, 0)
}

// getSpriteRect return a geometry.Rect for a given sprite.Sprite at a geometry.Point
func (sm storageManager) getSpriteRect(spr sprite.Sprite, at geometry.Point) geometry.Rect {
	def, _ := sm.GetSpriteDef(spr.Sheet, spr.Name)
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
func (sm storageManager) SpriteAtContains(spr sprite.Sprite, at geometry.Point, point geometry.Point) bool {
	return sm.getSpriteRect(spr, at).IsPointInRect(point)
}

// Storage returns a new storage.StorageManager
func Storage(dm DeviceManager) StorageManager {
	return &storageManager{
		sheets:    make(map[string]spriteSheet, 0),
		textures:  make(map[string]components.TextureDef, 0),
		fonts:     make(map[string]components.FontDef, 0),
		musics:    make(map[string]components.MusicDef, 0),
		sounds:    make(map[string]components.SoundDef, 0),
		tiledMaps: make(map[string]components.TiledMapDef, 0),
		dm:        dm,
	}
}
