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

// Package components contains the components already defined
package components

import (
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/lafriks/go-tiled"
)

// TextureDef defines a texture
type TextureDef struct {
	Data interface{}   // Data is the texture data
	Size geometry.Size // Size is the texture size
}

// SpriteDef defines an sprite.Sprite
type SpriteDef struct {
	Texture TextureDef     // Texture is the TextureDef
	Origin  geometry.Rect  // Origin is where the sprite is on the texture
	Pivot   geometry.Point // Pivot is the relative pivot 0..1 in each axis
}

// FontDef defines a font
type FontDef struct {
	Data interface{} // Data is the font data
}

// MusicDef defines a music stream
type MusicDef struct {
	Data interface{}
}

// SoundDef defines a sound wave
type SoundDef struct {
	Data interface{}
}

// TiledMapDef defines a tiled map
type TiledMapDef struct {
	Data       *tiled.Map                   // Data is the internal map data
	Cols       int32                        // Cols is the number of cols in the map
	Rows       int32                        // Rows is the number of rows in the map
	Size       geometry.Size                // Size is the map size in pixels
	TileSize   geometry.Size                // TileSize is the size per each tile in pixels
	Properties map[uint32]map[string]string // Properties are properties per each global block id
}
