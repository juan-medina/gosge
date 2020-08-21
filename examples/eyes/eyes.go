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
	"github.com/juan-medina/gosge/pkg/components/position"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"math"
	"reflect"
)

// our game options
var opt = options.Options{
	Title:      "Eyes Game",
	Width:      1920,
	Height:     1080,
	ClearColor: color.Gopher,
}

// entities that we are going to create
var (
	nose          *entity.Entity
	leftExterior  *entity.Entity
	leftInterior  *entity.Entity
	rightInterior *entity.Entity
	rightExterior *entity.Entity
	bottomText    *entity.Entity
)

// game constants
const (
	noseVerticalGap = 300
	eyesGap         = 400
	textSize        = 40
	eyeRadiusWidth  = 100
	eyeRadiusHeight = 90
)

func loadGame(eng engine.Engine) error {
	// pre load sprites
	if err := eng.LoadSpriteSheet("resources/gopher.json"); err != nil {
		return err
	}

	// get the world
	gw := eng.World()

	// add the sprites
	nose = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "nose"}))
	leftExterior = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior"}))
	leftInterior = gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior"},
		lookAtMouse{
			pivot:  leftExterior,
			radius: position.Position{X: eyeRadiusWidth, Y: eyeRadiusHeight},
		},
	))
	rightExterior = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior"}))
	rightInterior = gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior"},
		lookAtMouse{
			pivot:  rightExterior,
			radius: position.Position{X: eyeRadiusWidth, Y: eyeRadiusHeight},
		},
	))

	// add our text
	bottomText = gw.Add(entity.New(
		text.Text{
			String:     "press <ESC> to close",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.BottomVAlignment,
		},
		color.Black,
	))

	// add our layout system
	gw.AddSystem(&layoutSystem{})
	// add our look at mouse system
	gw.AddSystem(&lookAtMouseSystem{})

	return nil
}

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

// component to make an entity to look at mouse with a pivot
type lookAtMouse struct {
	pivot  *entity.Entity
	radius position.Position
}

var types = struct{ lookAtMouse reflect.Type }{lookAtMouse: reflect.TypeOf(lookAtMouse{})}

func getLookAtMouse(e *entity.Entity) lookAtMouse {
	return e.Get(types.lookAtMouse).(lookAtMouse)
}

// system that make entities to look at the mouse
type lookAtMouseSystem struct {
	scaleX float64
	scaleY float64
}

func (lam lookAtMouseSystem) Update(_ *world.World, _ float64) error {
	return nil
}

func (lam *lookAtMouseSystem) Notify(gw *world.World, event interface{}, _ float64) error {
	switch ev := event.(type) {
	// if the screen has change size
	case events.ScreenSizeChangeEvent:
		// save the scale
		lam.scaleX = ev.Scale.X
		lam.scaleY = ev.Scale.Y
	// if we move the mouse
	case events.MouseMoveEvent:
		// get the entities that look at the mouse
		for _, v := range gw.Entities(types.lookAtMouse) {
			la := getLookAtMouse(v)
			// make this entity to look at the mouse
			lam.lookAt(v, la, ev.Position)
		}
	}
	return nil
}

func (lam lookAtMouseSystem) lookAt(ent *entity.Entity, la lookAtMouse, mouse position.Position) {
	pos := position.Get(la.pivot)

	dx := mouse.X - pos.X
	dy := mouse.Y - pos.Y
	angle := math.Atan2(dy, dx)

	ax := la.radius.X * lam.scaleX * math.Cos(angle)
	ay := la.radius.Y * lam.scaleY * math.Sin(angle)

	np := position.Position{
		X: pos.X + ax,
		Y: pos.Y + ay,
	}

	ent.Set(np)
}

// system that layout/scale our entities when the screen size change
type layoutSystem struct{}

func (ls layoutSystem) Update(_ *world.World, _ float64) error {
	return nil
}

func (ls *layoutSystem) Notify(_ *world.World, event interface{}, _ float64) error {
	switch ev := event.(type) {
	// if the screen has change size
	case events.ScreenSizeChangeEvent:
		// change layout
		ls.positionElements(ev)
	}

	return nil
}

// layout our entities
func (ls layoutSystem) positionElements(event events.ScreenSizeChangeEvent) {
	scale := event.Scale.Max
	// the nose is in the middle and a bit down
	nosePos := position.Position{
		X: float64(event.Current.Width / 2),
		Y: float64(event.Current.Height/2) + (noseVerticalGap * scale),
	}

	// left eye is a bit up left of the nose
	leftEyePos := position.Position{
		X: nosePos.X - (eyesGap * scale),
		Y: nosePos.Y - (eyesGap * scale),
	}

	// right eye is a bit up right of the nose
	rightEyePos := position.Position{
		X: nosePos.X + (eyesGap * scale),
		Y: leftEyePos.Y,
	}

	// update sprites pos and scale
	ls.setPosAndScale(nose, nosePos, scale)
	ls.setPosAndScale(leftExterior, leftEyePos, scale)
	ls.setPosAndScale(leftInterior, leftEyePos, scale)
	ls.setPosAndScale(rightExterior, rightEyePos, scale)
	ls.setPosAndScale(rightInterior, rightEyePos, scale)

	// the text is bottom center
	textPos := position.Position{
		X: float64(event.Current.Width) / 2,
		Y: float64(event.Current.Height),
	}
	// update text pos and scale
	ls.setPosAndScale(bottomText, textPos, scale)
}

// set the position and scale of the objects
func (ls layoutSystem) setPosAndScale(ent *entity.Entity, pos position.Position, scale float64) {
	// set pos
	ent.Set(pos)

	if ent.Contains(sprite.TYPE) {
		// update sprite scale
		sp := sprite.Get(ent)
		sp.Scale = scale
		ent.Set(sp)
	}

	if ent.Contains(text.TYPE) {
		// update text scale
		tx := text.Get(ent)
		tx.Size = textSize * scale
		tx.Spacing = textSize / 4 * scale
		ent.Set(tx)
	}
}
