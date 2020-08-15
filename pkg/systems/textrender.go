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
	"github.com/hajimehoshi/ebiten"
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
	faces      map[float64]font.Face
	pt2PxRatio float64
}

func (r textRenderingSystem) getFace(size float64) font.Face {
	tt, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(tt, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
		DPI:     72,
	})

	return face
}

func (r textRenderingSystem) getCachedFace(size float64) font.Face {
	if value, ok := r.faces[size]; ok {
		return value
	} else {
		face := r.getFace(size)
		r.faces[size] = face
		return face
	}
}

func (r textRenderingSystem) drawBox(ctx render.Context, x float64, y float64, width float64, height float64, fill color.Color) {
	s := ebiten.DeviceScaleFactor()

	subImg, _ := ebiten.NewImage(int(width), int(height), ebiten.FilterDefault)
	_ = subImg.Fill(fill)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(s, s)

	_ = ctx.Image.DrawImage(subImg, op)
	_ = subImg.Clear()
}

func (r textRenderingSystem) drawColorBox(ctx render.Context, x float64, y float64, width float64, height float64) {
	w := width / 2
	h := height / 2

	r.drawBox(ctx, x, y, w, h, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	r.drawBox(ctx, x, y-h, w, h, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	r.drawBox(ctx, x-w, y, w, h, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	r.drawBox(ctx, x-w, y-h, w, h, color.RGBA{R: 255, G: 255, B: 0, A: 255})
}

func (r textRenderingSystem) Update(view *view.View) {
	context := view.Entity(render.ContextType).Get(render.ContextType).(render.Context)
	for _, v := range view.Entities(components.TextType, components.PosType, components.ColorType) {
		textCmp := v.Get(components.TextType).(components.Text)
		posCmp := v.Get(components.PosType).(components.Pos)
		colorCmp := v.Get(components.ColorType).(color.Color)

		s := ebiten.DeviceScaleFactor()
		x := posCmp.X * s
		y := posCmp.Y * s

		fs := float64(textCmp.Size) * r.pt2PxRatio * s
		ttFace := r.getFace(fs)

		text.Draw(context.Image, textCmp.String, ttFace, int(x), int(y), colorCmp)
	}
}

func (r textRenderingSystem) calculatePt2pxRatio() float64 {
	const samplePoints = 20
	face := r.getFace(samplePoints)
	bounds, _, _ := face.GlyphBounds('M')
	return samplePoints / float64((bounds.Max.Y - bounds.Min.Y).Ceil())
}

func TextRenderingSystem() system.System {
	trs := textRenderingSystem{
		faces: make(map[float64]font.Face),
	}

	trs.pt2PxRatio = trs.calculatePt2pxRatio()

	return trs
}
