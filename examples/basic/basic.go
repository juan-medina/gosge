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
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/engine"
	"image/color"
)

type myGame struct {
}

func (m *myGame) Load(eng engine.Engine) {
	eng.World().Add(entitiy.New(
		components.GameSettings{Width: 1920, Height: 1080, Title: "Basic Game"}),
	)
}

func (m *myGame) Init(eng engine.Engine) {
	settings := eng.World().Entity(components.GameSettingsType).Get(components.GameSettingsType).(components.GameSettings)
	eng.World().Add(entitiy.New(
		components.Text{
			String:     "Hello world",
			Size:       100,
			HAlignment: components.CenterHAlignment,
			VAlignment: components.MiddleVAlignment,
		},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		components.Pos{X: float64(settings.Width) / 2, Y: float64(settings.Height) / 2}),
	)
}

func main() {
	engine.Run(&myGame{})
}
