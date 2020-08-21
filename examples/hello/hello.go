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
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/position"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"reflect"
)

var opt = options.Options{
	Title:      "Hello Game",
	Width:      1920,
	Height:     1080,
	ClearColor: color.Black,
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

var types = struct{ stickyText reflect.Type }{stickyText: reflect.TypeOf(stickyText{})}

func getStickyText(e *entity.Entity) stickyText {
	return e.Get(types.stickyText).(stickyText)
}

type stickyTextSystem struct{}

func (sts stickyTextSystem) Update(_ *world.World, _ float64) error {
	return nil
}

func (sts stickyTextSystem) Notify(wld *world.World, event interface{}, _ float64) error {
	switch e := event.(type) {
	case events.ScreenSizeChangeEvent:
		// get all our texts
		for _, v := range wld.Entities(text.TYPE, types.stickyText) {
			// get the text components
			txt := text.Get(v)
			st := getStickyText(v)

			// change text size & spacing from current scale
			txt.Size = st.size * e.Scale.Min
			txt.Spacing = st.spacing * e.Scale.Min
			v.Set(txt)

			// calculate position based on current screen size and sticky
			pos := position.Position{}
			switch st.stick {
			case stickToCenter:
				pos.X = float64(e.Current.Width) / 2
				pos.Y = float64(e.Current.Height) / 2
			case stickToBottom:
				pos.X = float64(e.Current.Width) / 2
				pos.Y = float64(e.Current.Height)
			}
			v.Set(pos)
		}
	}

	return nil
}

func loadGame(eng engine.Engine) error {
	gWorld := eng.World()
	gWorld.Add(entity.New(
		text.Text{
			String:     "Hello World",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.MiddleVAlignment,
		},
		effects.AlternateColor{
			Time:  2,
			Delay: 1,
			From:  color.Red,
			To:    color.Yellow,
		},
		stickyText{size: 300, spacing: 10, stick: stickToCenter},
	))
	gWorld.Add(entity.New(
		text.Text{
			String:     "press <ESC> to close",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.BottomVAlignment,
		},
		effects.AlternateColor{
			Time: .25,
			From: color.White,
			To:   color.White.Alpha(0),
		},
		stickyText{size: 60, spacing: 10, stick: stickToBottom},
	))
	gWorld.AddSystem(stickyTextSystem{})
	return nil
}

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}
