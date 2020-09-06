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

package managers

import (
	"fmt"
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge/components/animation"
	"github.com/juan-medina/gosge/components/sprite"
)

type animationManager struct{}

func (animationManager) System(world *goecs.World, delta float32) error {
	for it := world.Iterator(animation.TYPE.Animation); it != nil; it = it.Next() {
		ent := it.Value()

		anim := animation.Get.Animation(ent)

		var state animation.State

		if ent.Contains(animation.TYPE.State) {
			state = animation.Get.State(ent)
		} else {
			state = animation.State{}
		}

		if state.Current != anim.Current {
			state = animation.State{}
			state.Current = anim.Current
			if _, ok := anim.Sequences[anim.Current]; !ok {
				return fmt.Errorf("can not find animation: %q", anim.Current)
			}
			seq := anim.Sequences[anim.Current]
			ent.Set(seq)
			ent.Set(state)
			ent.Set(anim)
		}

		if state.Speed != anim.Speed {
			state.Speed = anim.Speed
			ent.Set(state)
		}
	}

	for it := world.Iterator(animation.TYPE.Sequence, animation.TYPE.Animation, animation.TYPE.State); it != nil; it = it.Next() {
		ent := it.Value()

		seq := animation.Get.Sequence(ent)
		anim := animation.Get.Animation(ent)
		state := animation.Get.State(ent)

		if state.Time += delta * state.Speed; state.Time > seq.Delay {
			state.Time = 0
			if state.Frame++; state.Frame >= seq.Frames {
				state.Frame = 0
			}
		}

		spriteName := fmt.Sprintf(seq.Base, state.Frame+1)

		spr := sprite.Sprite{
			Sheet:    seq.Sheet,
			Name:     spriteName,
			Scale:    seq.Scale,
			Rotation: seq.Rotation,
			FlipX:    anim.FlipX,
			FlipY:    anim.FlipY,
		}

		ent.Set(spr)
		ent.Set(seq)
		ent.Set(state)
		ent.Set(anim)
	}

	return nil
}

// Animation returns a managers.WithSystem for handling animations
func Animation() WithSystem {
	return &animationManager{}
}
