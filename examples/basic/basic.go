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

func (m *myGame) Init(e engine.Engine) {
	e.World().Add(entitiy.New().
		Add(components.GameSettings{Width: 640, Height: 480, Title: "Basic Game"}),
	)
	e.World().Add(entitiy.New().
		Add(components.Text{String: "Hello world", Size: 50, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}}).
		Add(components.Pos{X: 100, Y: 100}),
	)
}

func main() {
	engine.Run(&myGame{})
}
