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
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/juan-medina/goecs/pkg/system"
	"github.com/juan-medina/goecs/pkg/view"
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/render"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
)

type textRenderingSystem struct {
	faces map[float64]font.Face
}

func (r textRenderingSystem) getFace(size float64) font.Face {
	if value, ok := r.faces[size]; ok {
		return value
	} else {
		tt, _ := truetype.Parse(goregular.TTF)
		face := truetype.NewFace(tt, &truetype.Options{
			Size:    size,
			Hinting: font.HintingFull,
		})
		r.faces[size] = face
		return face
	}
}

func (r textRenderingSystem) Update(view *view.View) {
	context := view.Entity(render.ContextType).Get(render.ContextType).(render.Context)
	for _, v := range view.Entities(components.TextType, components.PosType, components.ColorType) {
		textCmp := v.Get(components.TextType).(components.Text)
		posCmp := v.Get(components.PosType).(components.Pos)
		colorCmp := v.Get(components.ColorType).(color.Color)
		ttFace := r.getFace(float64(textCmp.Size))
		text.Draw(context.Image, textCmp.String, ttFace, int(posCmp.X), int(posCmp.Y), colorCmp)
	}
}

func TextRenderingSystem() system.System {
	return textRenderingSystem{
		faces: make(map[float64]font.Face),
	}
}
