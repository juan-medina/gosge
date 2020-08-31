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

func (rcs alternateColorSystem) Update(world *world.World, delta float32) error {
	for it := world.Iterator(effects.TYPE.AlternateColor); it.HasNext(); {
		ent := it.Value()
		ac := effects.Get.AlternateColor(ent)

		// init the state or get it from entity
		var state effects.AlternateColorState
		if ent.NotContains(effects.TYPE.AlternateColorState) {
			state = effects.AlternateColorState{
				CurrentTime: 0,
				Running:     false,
				From:        ac.From,
				To:          ac.To,
			}
		} else {
			state = effects.Get.AlternateColorState(ent)
		}

		// calculate color base on the state
		var clr color.Solid
		if !state.Running {
			clr = state.From
		} else {
			clr = state.From.Blend(state.To, state.CurrentTime/ac.Time)
		}

		// move current time
		state.CurrentTime += delta

		// if we are running
		if state.Running {
			// when we are on the time stop and swap colors in the state
			if state.CurrentTime > ac.Time {
				state.Running = false
				state.CurrentTime = 0
				aux := state.From
				state.From = state.To
				state.To = aux
			}
		} else if state.CurrentTime > ac.Delay { // wait the delay time
			state.Running = true
		}

		ent.Set(state)
		ent.Set(clr)
	}
	return nil
}

func (rcs alternateColorSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}

// AlternateColorSystem returns a world.System that handle effects.AlternateColor
func AlternateColorSystem() world.System {
	return &alternateColorSystem{}
}
