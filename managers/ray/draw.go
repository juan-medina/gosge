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
	"github.com/juan-medina/gosge/components"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
)

func (dmi *DeviceManagerImpl) color2RayColor(color color.Solid) rl.Color {
	return rl.NewColor(color.R, color.G, color.B, color.A)
}

var (
	emptyTexture = components.TextureDef{}
	emptyFont    = components.FontDef{}
)

// LoadTexture giving it file name into VRAM
func (dmi DeviceManagerImpl) LoadTexture(fileName string) (components.TextureDef, error) {
	if t := rl.LoadTexture(fileName); t.ID != 0 {
		return components.TextureDef{Data: t, Size: geometry.Size{Width: float32(t.Width), Height: float32(t.Height)}}, nil
	}
	return emptyTexture, fmt.Errorf("error loading texture: %q", fileName)
}

// LoadFont giving it file name into VRAM
func (dmi DeviceManagerImpl) LoadFont(fileName string) (components.FontDef, error) {
	if f := rl.LoadFont(fileName); f.Texture.ID != 0 {
		return components.FontDef{Data: f}, nil
	}
	return emptyFont, fmt.Errorf("error loading font: %q", fileName)
}

// UnloadFont from VRAM
func (dmi DeviceManagerImpl) UnloadFont(textureDef components.FontDef) {
	rl.UnloadFont(textureDef.Data.(rl.Font))
}

// UnloadTexture from VRAM
func (dmi DeviceManagerImpl) UnloadTexture(textureDef components.TextureDef) {
	rl.UnloadTexture(textureDef.Data.(rl.Texture2D))
}

// DrawText will draw a text.Text in the given geometry.Point with the correspondent color.Color
func (dmi DeviceManagerImpl) DrawText(ftd components.FontDef, txt ui.Text, pos geometry.Point, color color.Solid) {
	font := ftd.Data.(rl.Font)

	vec := rl.Vector2{
		X: pos.X,
		Y: pos.Y,
	}

	if txt.HAlignment != ui.LeftHAlignment || txt.VAlignment != ui.BottomVAlignment {
		av := rl.MeasureTextEx(font, txt.String, txt.Size, 0)

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

	rl.DrawTextEx(font, txt.String, vec, txt.Size, 0, dmi.color2RayColor(color))
}

// DrawSprite draws a sprite.Sprite in the given geometry.Point with the tint color.Color
func (dmi DeviceManagerImpl) DrawSprite(def components.SpriteDef, sprite sprite.Sprite, pos geometry.Point, tint color.Solid) error {
	scale := sprite.Scale
	px := def.Origin.Size.Width * def.Pivot.X
	py := def.Origin.Size.Height * def.Pivot.Y

	sourceFlip := geometry.Size{
		Width:  def.Origin.Size.Width,
		Height: def.Origin.Size.Height,
	}

	if sprite.FlipX {
		sourceFlip.Width *= -1
	}
	if sprite.FlipY {
		sourceFlip.Height *= -1
	}

	rc := dmi.color2RayColor(tint)
	rotation := sprite.Rotation
	sourceRec := rl.Rectangle{
		X:      def.Origin.From.X,
		Y:      def.Origin.From.Y,
		Width:  sourceFlip.Width,
		Height: sourceFlip.Height,
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
func (dmi DeviceManagerImpl) DrawSolidBox(pos geometry.Point, box shapes.SolidBox, solid color.Solid) {
	rec := rl.Rectangle{
		X:      pos.X,
		Y:      pos.Y,
		Width:  box.Size.Width * box.Scale,
		Height: box.Size.Height * box.Scale,
	}

	rl.DrawRectangleRec(rec, dmi.color2RayColor(solid))
}

// BeginScissor start a scissor draw (define screen area for following drawing)
func (dmi DeviceManagerImpl) BeginScissor(from geometry.Point, size geometry.Size) {
	rl.BeginScissorMode(int32(from.X), int32(from.Y), int32(size.Width), int32(size.Height))
}

// EndScissor end scissor
func (dmi DeviceManagerImpl) EndScissor() {
	rl.EndScissorMode()
}

// DrawBox draws a box outline with an color.Solid and a scale
func (dmi DeviceManagerImpl) DrawBox(pos geometry.Point, box shapes.Box, solid color.Solid) {
	rec := rl.Rectangle{
		X:      pos.X,
		Y:      pos.Y,
		Width:  box.Size.Width * box.Scale,
		Height: box.Size.Height * box.Scale,
	}

	rl.DrawRectangleLinesEx(rec, box.Thickness, dmi.color2RayColor(solid))
}

// DrawGradientBox draws a solid box with an color.Solid and a scale
func (dmi DeviceManagerImpl) DrawGradientBox(pos geometry.Point, box shapes.SolidBox, gradient color.Gradient) {
	x := int32(pos.X)
	y := int32(pos.Y)
	w := int32(box.Size.Width * box.Scale)
	h := int32(box.Size.Height * box.Scale)
	c1 := dmi.color2RayColor(gradient.From)
	c2 := dmi.color2RayColor(gradient.To)

	if gradient.Direction == color.GradientHorizontal {
		rl.DrawRectangleGradientH(x, y, w, h, c1, c2)
	} else {
		rl.DrawRectangleGradientV(x, y, w, h, c1, c2)
	}
}

// SetBackgroundColor changes the current background color.Solid
func (dmi *DeviceManagerImpl) SetBackgroundColor(color color.Solid) {
	dmi.saveOpts.BackGround = color
}

// MeasureText return the geometry.Size of a string with a defined size and spacing
func (dmi *DeviceManagerImpl) MeasureText(fnt components.FontDef, str string, size float32) geometry.Size {
	fray := fnt.Data.(rl.Font)
	av := rl.MeasureTextEx(fray, str, size, 0)
	return geometry.Size{
		Width:  av.X,
		Height: av.Y,
	}
}

// DrawLine between from and to with a given thickness and color.Solid
func (dmi DeviceManagerImpl) DrawLine(from, to geometry.Point, thickness float32, color color.Solid) {
	rf := rl.Vector2{
		X: from.X,
		Y: from.Y,
	}
	rt := rl.Vector2{
		X: to.X,
		Y: to.Y,
	}
	rl.DrawLineEx(rf, rt, thickness, dmi.color2RayColor(color))
}
