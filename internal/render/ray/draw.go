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

package ray

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/juan-medina/gosge/internal/components"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
)

func (rr *RenderImpl) color2RayColor(color color.Solid) rl.Color {
	return rl.NewColor(color.R, color.G, color.B, color.A)
}

// LoadTexture giving it file name into VRAM
func (rr RenderImpl) LoadTexture(fileName string) error {
	if t := rl.LoadTexture(fileName); t.ID != 0 {
		rr.textureHold[fileName] = t
		return nil
	}
	return fmt.Errorf("error loading texture: %q", fileName)
}

// UnloadAllTextures from VRAM
func (rr RenderImpl) UnloadAllTextures() {
	for k, v := range rr.textureHold {
		delete(rr.textureHold, k)
		rl.UnloadTexture(v)
	}
}

// DrawText will draw a text.Text in the given geometry.Position with the correspondent color.Color
func (rr RenderImpl) DrawText(txt text.Text, pos geometry.Position, color color.Solid) {
	font := rl.GetFontDefault()

	vec := rl.Vector2{
		X: pos.X,
		Y: pos.Y,
	}

	if txt.HAlignment != text.LeftHAlignment || txt.VAlignment != text.BottomVAlignment {
		av := rl.MeasureTextEx(font, txt.String, txt.Size, txt.Spacing)

		switch txt.HAlignment {
		case text.LeftHAlignment:
			av.X = 0
		case text.CenterHAlignment:
			av.X = -av.X / 2
		case text.RightHAlignment:
			av.X = -av.X
		}

		switch txt.VAlignment {
		case text.BottomVAlignment:
			av.Y = -av.Y
		case text.MiddleVAlignment:
			av.Y = -av.Y / 2
		case text.TopVAlignment:
			av.Y = 0
			break
		}
		vec.X += av.X
		vec.Y += av.Y
	}

	rl.DrawTextEx(font, txt.String, vec, txt.Size, txt.Spacing, rr.color2RayColor(color))
}

// DrawSprite draws a sprite.Sprite in the given geometry.Position with the tint color.Color
func (rr RenderImpl) DrawSprite(def components.SpriteDef, sprite sprite.Sprite, pos geometry.Position, tint color.Solid) error {
	if val, ok := rr.textureHold[def.Texture]; ok {
		scale := sprite.Scale
		px := def.Origin.Size.Width * def.Pivot.X
		py := def.Origin.Size.Height * def.Pivot.Y
		rc := rr.color2RayColor(tint)
		rotation := sprite.Rotation
		sourceRec := rl.Rectangle{
			X:      def.Origin.From.X,
			Y:      def.Origin.From.Y,
			Width:  def.Origin.Size.Width,
			Height: def.Origin.Size.Height,
		}
		destRec := rl.Rectangle{
			X:      pos.X - (px * scale),
			Y:      pos.Y - (py * scale),
			Width:  def.Origin.Size.Width * scale,
			Height: def.Origin.Size.Height * scale,
		}
		origin := rl.Vector2{X: 0, Y: 0}
		rl.DrawTexturePro(val, sourceRec, destRec, origin, rotation, rc)
	} else {
		return fmt.Errorf("error drawing sprite, texture not found: %q", sprite.Name)
	}
	return nil
}

// DrawSolidBox draws a solid box with an color.Solid and a scale
func (rr RenderImpl) DrawSolidBox(pos geometry.Position, box shapes.Box, solid color.Solid) {
	rec := rl.Rectangle{
		X:      pos.X,
		Y:      pos.Y,
		Width:  box.Size.Width * box.Scale,
		Height: box.Size.Height * box.Scale,
	}

	rl.DrawRectangleRec(rec, rr.color2RayColor(solid))
}

// DrawGradientBox draws a solid box with an color.Solid and a scale
func (rr RenderImpl) DrawGradientBox(pos geometry.Position, box shapes.Box, gradient color.Gradient) {
	x := int32(pos.X)
	y := int32(pos.Y)
	w := int32(box.Size.Width * box.Scale)
	h := int32(box.Size.Height * box.Scale)
	c1 := rr.color2RayColor(gradient.From)
	c2 := rr.color2RayColor(gradient.To)

	if gradient.Direction == color.GradientHorizontal {
		rl.DrawRectangleGradientH(x, y, w, h, c1, c2)
	} else {
		rl.DrawRectangleGradientV(x, y, w, h, c1, c2)
	}
}
