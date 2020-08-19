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
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"reflect"
)

var opt = options.Options{
	Title:      "Hello Game",
	Width:      1920,
	Height:     1080,
	ClearColor: components.NewColor(0, 0, 0, 255),
}

type stickPosition int

const (
	stickToCenter = iota
	stickToBottom
)

type stickyText struct {
	size    float64
	spacing float64
	stick   stickPosition
}

var stickyTextType = reflect.TypeOf(stickyText{})

type stickyTextSystem struct{}

func (cts stickyTextSystem) Update(_ *world.World, _ float64) error {
	return nil
}

func (cts stickyTextSystem) Notify(wld *world.World, event interface{}, _ float64) error {
	switch e := event.(type) {
	case events.ScreenSizeChangeEvent:
		// get all our texts
		for _, v := range wld.Entities(components.PosType, components.UiTextType, stickyTextType) {
			// get the text components
			pos := v.Get(components.PosType).(components.Pos)
			text := v.Get(components.UiTextType).(components.UiText)
			st := v.Get(stickyTextType).(stickyText)

			// calculate position based on current screen size and stick
			switch st.stick {
			case stickToCenter:
				pos.X = float64(e.Current.Width) / 2
				pos.Y = float64(e.Current.Height) / 2
			case stickToBottom:
				pos.X = float64(e.Current.Width) / 2
				pos.Y = float64(e.Current.Height)
			}
			v.Set(pos)

			// change text size & spacing from current scale
			text.Size = st.size * e.Scale
			text.Spacing = st.spacing * e.Scale
			v.Set(text)
		}
	}

	return nil
}

func loadGame(gWorld *world.World) error {
	gWorld.Add(entitiy.New(
		components.UiText{
			String:     "Hello World",
			Size:       300,
			Spacing:    10,
			HAlignment: components.CenterHAlignment,
			VAlignment: components.MiddleVAlignment,
		},
		components.Pos{X: float64(opt.Width / 2), Y: float64(opt.Height / 2)},
		components.AlternateColor{
			Time: 1,
			From: components.NewColor(0, 255, 255, 255),
			To:   components.NewColor(255, 0, 0, 255),
		},
		stickyText{size: 300, spacing: 10, stick: stickToCenter},
	))
	gWorld.Add(entitiy.New(
		components.UiText{
			String:     "press <ESC> to close",
			Size:       60,
			Spacing:    10,
			HAlignment: components.CenterHAlignment,
			VAlignment: components.BottomVAlignment,
		},
		components.Pos{X: float64(opt.Width / 2), Y: float64(opt.Height)},
		components.AlternateColor{
			Time: .25,
			From: components.NewColor(255, 255, 255, 255),
			To:   components.NewColor(255, 255, 255, 0),
		},
		stickyText{size: 60, spacing: 10, stick: stickToBottom},
	))
	gWorld.AddSystem(stickyTextSystem{})
	return nil
}

func main() {
	if err := engine.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}
