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

package main

import (
	"github.com/juan-medina/goecs/pkg/entitiy"
	"github.com/juan-medina/goecs/pkg/view"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/options"
	"image/color"
)

var opt = options.Options{
	Title:      "Simple Game",
	Width:      1920,
	Height:     1080,
	ClearColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
}

var centerText *entitiy.Entity
var bottomText *entitiy.Entity

type centerTextSystem struct {
}

func (cts centerTextSystem) Update(view *view.View) {
	// get settings, including original and current screen size and the scale from original to current
	settings := view.Entity(components.GameSettingsType).Get(components.GameSettingsType).(components.GameSettings)

	// get the center text components
	pos := centerText.Get(components.PosType).(components.Pos)
	text := centerText.Get(components.UiTextType).(components.UiText)

	// calculate center text position base on current screen size
	pos.X = float64(settings.Current.Width) / 2
	pos.Y = float64(settings.Current.Height) / 2
	centerText.Set(pos)

	// change center text size & spacing from current scale
	text.Size = 300 * settings.Scale
	text.Spacing = 10 * settings.Scale
	centerText.Set(text)

	// get the bottom text components
	pos = bottomText.Get(components.PosType).(components.Pos)
	text = bottomText.Get(components.UiTextType).(components.UiText)

	// calculate bottom text position base on current screen size
	pos.X = float64(settings.Current.Width) / 2
	pos.Y = float64(settings.Current.Height)
	bottomText.Set(pos)

	// change bottom text size & spacing from current scale
	text.Size = 60 * settings.Scale
	text.Spacing = 10 * settings.Scale
	bottomText.Set(text)
}

func loadGame(gWorld *world.World) {
	centerText = gWorld.Add(entitiy.New(
		components.UiText{
			String:     "Hello world",
			Size:       300,
			Spacing:    10,
			HAlignment: components.CenterHAlignment,
			VAlignment: components.MiddleVAlignment,
		},
		components.Pos{X: 0, Y: 0},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	))
	bottomText = gWorld.Add(entitiy.New(
		components.UiText{
			String:     "press <ESC> to close",
			Size:       60,
			Spacing:    10,
			HAlignment: components.CenterHAlignment,
			VAlignment: components.BottomVAlignment,
		},
		components.Pos{X: 0, Y: 0},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	))
	gWorld.AddSystem(centerTextSystem{})
}

func main() {
	engine.Run(opt, loadGame)
}
