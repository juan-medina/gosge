package systems

import (
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/internal/render"
	"github.com/juan-medina/gosge/pkg/events"
	"math"
)

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

type eventSystem struct {
	tt  float64
	sse events.ScreenSizeChangeEvent
	wld *world.World
}

func (es eventSystem) Notify(_ *world.World, _ interface{}, _ float64) error {
	return nil
}

func (es *eventSystem) sendScreenSizeChange() error {
	w, h := render.GetScreenSize()
	es.sse.Current.Width = w
	es.sse.Current.Height = h

	sx := float64(es.sse.Current.Width) / float64(es.sse.Original.Width)
	sy := float64(es.sse.Current.Height) / float64(es.sse.Original.Height)
	es.sse.Scale.Min = math.Min(sx, sy)
	es.sse.Scale.Max = math.Max(sx, sy)
	es.sse.Scale.X = sx
	es.sse.Scale.Y = sy

	return es.wld.Notify(es.sse)
}

func (es eventSystem) sendGameClose() error {
	return es.wld.Notify(es.sse)
}

func (es *eventSystem) initialize(world *world.World) error {
	es.wld = world

	w, h := render.GetScreenSize()

	es.sse.Original.Width = w
	es.sse.Original.Height = h
	es.sse.Current.Width = w
	es.sse.Current.Height = h
	es.sse.Scale.Min = 1
	es.sse.Scale.Max = 1
	es.sse.Scale.X = 1
	es.sse.Scale.Y = 1

	return es.sendScreenSizeChange()
}

func (es *eventSystem) Update(world *world.World, delta float64) error {
	if es.tt == 0 {
		if err := es.initialize(world); err != nil {
			return err
		}
	}

	if render.ShouldClose() {
		if err := es.sendGameClose(); err != nil {
			return err
		}
	}

	if render.IsScreenScaleChange() {
		if err := es.sendScreenSizeChange(); err != nil {
			return err
		}
	}

	es.tt += delta
	return nil
}

// EventSystem returns a world.System that will handle events
func EventSystem() world.System {
	return &eventSystem{}
}
