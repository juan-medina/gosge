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
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
)

type alternateColorSystem struct{}

func (rcs alternateColorSystem) Update(world *world.World, delta float64) error {
	for _, v := range world.Entities(effects.TYPES.AlternateColor) {
		ac := v.Get(effects.TYPES.AlternateColor).(effects.AlternateColor)

		var clr color.Color

		switch ac.State {
		case effects.NoState, effects.StateStopped:
			clr = ac.From
		case effects.StateRunning:
			clr = ac.From.Blend(ac.To, ac.Current/ac.Time)
		}

		ac.Current += delta

		switch ac.State {
		case effects.NoState, effects.StateStopped:
			if ac.Current > ac.Delay {
				ac.State = effects.StateRunning
				ac.Current = 0
			}
		case effects.StateRunning:
			if ac.Current > ac.Time {
				ac.State = effects.StateStopped
				ac.Current = 0
				aux := ac.From
				ac.From = ac.To
				ac.To = aux
			}
		}

		v.Set(clr)
		v.Set(ac)
	}
	return nil
}

func (rcs alternateColorSystem) Notify(_ *world.World, _ interface{}, _ float64) error {
	return nil
}

// AlternateColorSystem returns a world.System that handle effects.AlternateColor
func AlternateColorSystem() world.System {
	return &alternateColorSystem{}
}
