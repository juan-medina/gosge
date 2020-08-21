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
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
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
	Title: "Eyes Game",
	Size: geometry.Size{
		Width:  1920,
		Height: 1080,
	},
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
	dizzyBar      *entity.Entity
	dizzyBarEmpty *entity.Entity
	dizzyText     *entity.Entity
)

// game constants
const (
	noseVerticalGap = 300
	eyesGap         = 400
	textSmallSize   = 40
	textBigSize     = 70
	dizzyBarWith    = 1500
	dizzyBarHeight  = 100
	dizzyBarGap     = 50
	maxDizzy        = 2.5
	dizzyGainRate   = 6
	dizzyReduceRate = 1
)

func loadGame(eng engine.Engine) error {
	var err error
	// pre load sprites
	if err = eng.LoadSpriteSheet("resources/gopher.json"); err != nil {
		return err
	}

	// get the world
	gw := eng.World()

	// Get the eyes size
	var eyeSize geometry.Size
	var eyeRadius geometry.Position
	if eyeSize, err = eng.GetSpriteSize("resources/gopher.json", "eye_exterior"); err == nil {
		eyeRadius = geometry.Position{X: eyeSize.Width / 4, Y: eyeSize.Height / 4}
	} else {
		return err
	}

	// add the sprites
	nose = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "nose"}))
	leftExterior = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior"}))
	leftInterior = gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior"},
		lookAtMouse{pivot: leftExterior, radius: eyeRadius},
	))
	rightExterior = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior"}))
	rightInterior = gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior"},
		lookAtMouse{pivot: rightExterior, radius: eyeRadius},
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

	// Add the dizzy bar
	dizzyBar = gw.Add(entity.New(
		shapes.Box{},
		color.Gradient{From: color.Green, To: color.Red},
	))

	dizzyBarEmpty = gw.Add(entity.New(
		shapes.Box{},
		color.Black,
	))

	// add the dizzy text
	dizzyText = gw.Add(entity.New(
		text.Text{
			String:     "Dizzy! Level",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.MiddleVAlignment,
		},
		effects.AlternateColor{
			Time:  0.25,
			Delay: 0.25,
			From:  color.Red,
			To:    color.Yellow,
		},
	))

	// add our layout system
	gw.AddSystem(&layoutSystem{})
	// add our look at mouse system
	gw.AddSystem(&lookAtMouseSystem{})
	// add our dizzy bar system
	gw.AddSystem(&dizzyBarSystem{})

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
	radius geometry.Position
}

var types = struct{ lookAtMouse reflect.Type }{lookAtMouse: reflect.TypeOf(lookAtMouse{})}

func getLookAtMouse(e *entity.Entity) lookAtMouse {
	return e.Get(types.lookAtMouse).(lookAtMouse)
}

// system that make entities to look at the mouse
type lookAtMouseSystem struct {
	scaleX float32
	scaleY float32
}

func (lam lookAtMouseSystem) Update(_ *world.World, _ float32) error {
	return nil
}

func (lam *lookAtMouseSystem) Notify(gw *world.World, event interface{}, _ float32) error {
	switch ev := event.(type) {
	// if the screen has change size
	case events.ScreenSizeChangeEvent:
		// save the scale
		lam.scaleX = ev.Scale.Point.X
		lam.scaleY = ev.Scale.Point.Y
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

func (lam lookAtMouseSystem) lookAt(ent *entity.Entity, la lookAtMouse, mouse geometry.Position) {
	pos := geometry.Get.Position(la.pivot)

	dx := mouse.X - pos.X
	dy := mouse.Y - pos.Y
	angle := float32(math.Atan2(float64(dy), float64(dx)))

	ax := la.radius.X * lam.scaleX * float32(math.Cos(float64(angle)))
	ay := la.radius.Y * lam.scaleY * float32(math.Sin(float64(angle)))

	np := geometry.Position{
		X: pos.X + ax,
		Y: pos.Y + ay,
	}

	ent.Set(np)
}

// system that layout/scale our entities when the screen size change
type layoutSystem struct{}

func (ls layoutSystem) Update(_ *world.World, _ float32) error {
	return nil
}

func (ls *layoutSystem) Notify(_ *world.World, event interface{}, _ float32) error {
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
	nosePos := geometry.Position{
		X: event.Current.Width / 2,
		Y: event.Current.Height/2 + (noseVerticalGap * scale),
	}

	// left eye is a bit up left of the nose
	leftEyePos := geometry.Position{
		X: nosePos.X - (eyesGap * scale),
		Y: nosePos.Y - (eyesGap * scale),
	}

	// right eye is a bit up right of the nose
	rightEyePos := geometry.Position{
		X: nosePos.X + (eyesGap * scale),
		Y: leftEyePos.Y,
	}

	// update sprites pos and scale
	ls.setPosAndScaleSprite(nose, nosePos, scale)
	ls.setPosAndScaleSprite(leftExterior, leftEyePos, scale)
	ls.setPosAndScaleSprite(leftInterior, leftEyePos, scale)
	ls.setPosAndScaleSprite(rightExterior, rightEyePos, scale)
	ls.setPosAndScaleSprite(rightInterior, rightEyePos, scale)

	// the text is bottom center
	textPos := geometry.Position{
		X: event.Current.Width / 2,
		Y: event.Current.Height,
	}
	// update text pos and scale
	ls.setPosAndScaleText(bottomText, textPos, scale, textSmallSize)

	box := shapes.Get.Box(dizzyBar)
	box.Size.Width = dizzyBarWith
	box.Size.Height = dizzyBarHeight
	box.Scale = scale

	dizzyBarPosition := geometry.Position{
		X: (event.Current.Width - (box.Size.Width * scale)) / 2,
		Y: dizzyBarGap * scale,
	}
	dizzyBar.Set(dizzyBarPosition)
	dizzyBar.Set(box)

	dizzyTextPosition := geometry.Position{
		X: event.Current.Width / 2,
		Y: dizzyBarPosition.Y + (box.Size.Height / 2 * scale),
	}

	ls.setPosAndScaleText(dizzyText, dizzyTextPosition, scale, textBigSize)
}

// set the position and scale of the objects
func (ls layoutSystem) setPosAndScaleSprite(ent *entity.Entity, pos geometry.Position, scale float32) {
	// set pos
	ent.Set(pos)

	// update sprite scale
	sp := sprite.Get(ent)
	sp.Scale = scale
	ent.Set(sp)

}

// set the position and scale of the objects
func (ls layoutSystem) setPosAndScaleText(ent *entity.Entity, pos geometry.Position, size float32, scale float32) {
	// set pos
	ent.Set(pos)

	// update text scale
	tx := text.Get(ent)
	tx.Size = size * scale
	tx.Spacing = size / 4 * scale
	ent.Set(tx)

}

type dizzyBarSystem struct {
	dizzy    float32
	lasMouse geometry.Position
}

func (dbs *dizzyBarSystem) Update(_ *world.World, delta float32) error {
	// in each frame we reduce how dizzy we are
	dbs.dizzy = float32(math.Max(float64(dbs.dizzy-(delta/dizzyReduceRate)), 0))

	// calculate how dizzy we are in 0..1
	percent := 1 - (dbs.dizzy / maxDizzy)

	// if we have position
	if dizzyBar.Contains(geometry.TYPE.Position) {
		// get the position of the regular dizzy bar
		dizzyBarPosition := geometry.Get.Position(dizzyBar)
		// get the position
		box := shapes.Get.Box(dizzyBar)

		// calculate position and width
		box.Size.Width = box.Size.Width * percent
		dizzyBarPosition.X = dizzyBarPosition.X - (dizzyBarWith * box.Scale * (percent - 1))

		// set components
		dizzyBarEmpty.Set(dizzyBarPosition)
		dizzyBarEmpty.Set(box)

		// make the dizzy text to blink faster
		ac := effects.Get.AlternateColor(dizzyText)
		ac.Delay = 0.25 * percent
		ac.Time = 0.25 * (percent * 8)
		dizzyText.Set(ac)

		// color the eyes red
		clr := color.White.Blend(color.Red, (1-percent)/2)
		leftExterior.Set(clr)
		rightExterior.Set(clr)
	}

	return nil
}

func (dbs *dizzyBarSystem) Notify(_ *world.World, event interface{}, delta float32) error {
	switch event.(type) {
	case events.MouseMoveEvent:
		dbs.dizzy = float32(math.Min(float64(dbs.dizzy+(delta*dizzyGainRate)), 2.5))
	}
	return nil
}
