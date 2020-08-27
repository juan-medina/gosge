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
	"github.com/juan-medina/gosge/pkg/components/ui"
)

func (rr *RenderImpl) color2RayColor(color color.Solid) rl.Color {
	return rl.NewColor(color.R, color.G, color.B, color.A)
}

var (
	emptyTexture = components.TextureDef{}
)

// LoadTexture giving it file name into VRAM
func (rr RenderImpl) LoadTexture(fileName string) (components.TextureDef, error) {
	if t := rl.LoadTexture(fileName); t.ID != 0 {
		return components.TextureDef{Data: t, Size: geometry.Size{Width: float32(t.Width), Height: float32(t.Height)}}, nil
	}
	return emptyTexture, fmt.Errorf("error loading texture: %q", fileName)
}

// UnloadTexture from VRAM
func (rr RenderImpl) UnloadTexture(textureDef components.TextureDef) {
	rl.UnloadTexture(textureDef.Data.(rl.Texture2D))
}

// DrawText will draw a text.Text in the given geometry.Point with the correspondent color.Color
func (rr RenderImpl) DrawText(txt ui.Text, pos geometry.Point, color color.Solid) {
	font := rl.GetFontDefault()

	vec := rl.Vector2{
		X: pos.X,
		Y: pos.Y,
	}

	if txt.HAlignment != ui.LeftHAlignment || txt.VAlignment != ui.BottomVAlignment {
		av := rl.MeasureTextEx(font, txt.String, txt.Size, txt.Spacing)

		switch txt.HAlignment {
		case ui.LeftHAlignment:
			av.X = 0
		case ui.CenterHAlignment:
			av.X = -av.X / 2
		case ui.RightHAlignment:
			av.X = -av.X
		}

		switch txt.VAlignment {
		case ui.BottomVAlignment:
			av.Y = -av.Y
		case ui.MiddleVAlignment:
			av.Y = -av.Y / 2
		case ui.TopVAlignment:
			av.Y = 0
			break
		}
		vec.X += av.X
		vec.Y += av.Y
	}

	rl.DrawTextEx(font, txt.String, vec, txt.Size, txt.Spacing, rr.color2RayColor(color))
}

// DrawSprite draws a sprite.Sprite in the given geometry.Point with the tint color.Color
func (rr RenderImpl) DrawSprite(def components.SpriteDef, sprite sprite.Sprite, pos geometry.Point, tint color.Solid) error {
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
	texture := def.Texture.Data.(rl.Texture2D)
	rl.DrawTexturePro(texture, sourceRec, destRec, origin, rotation, rc)

	return nil
}

// DrawSolidBox draws a solid box with an color.Solid and a scale
func (rr RenderImpl) DrawSolidBox(pos geometry.Point, box shapes.Box, solid color.Solid) {
	rec := rl.Rectangle{
		X:      pos.X,
		Y:      pos.Y,
		Width:  box.Size.Width * box.Scale,
		Height: box.Size.Height * box.Scale,
	}

	rl.DrawRectangleRec(rec, rr.color2RayColor(solid))
}

// DrawGradientBox draws a solid box with an color.Solid and a scale
func (rr RenderImpl) DrawGradientBox(pos geometry.Point, box shapes.Box, gradient color.Gradient) {
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

// SetBackgroundColor changes the current background color.Solid
func (rr *RenderImpl) SetBackgroundColor(color color.Solid) {
	rr.saveOpts.BackGround = color
}

// MeasureText return the geometry.Size of a string with a defined size and spacing
func (rr *RenderImpl) MeasureText(str string, size, spacing float32) geometry.Size {
	font := rl.GetFontDefault()
	av := rl.MeasureTextEx(font, str, size, spacing)
	return geometry.Size{
		Width:  av.X,
		Height: av.Y,
	}
}
