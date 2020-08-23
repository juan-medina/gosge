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
	Title:      "Eyes Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
}

// entities that we are going to create
var (
	leftExterior  *entity.Entity
	rightExterior *entity.Entity
	dizzyBar      *entity.Entity
	dizzyBarEmpty *entity.Entity
	dizzyText     *entity.Entity

	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
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
	dizzyGainRate   = 5
	dizzyReduceRate = 0.5
)

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

func loadGame(eng engine.Engine) error {
	var err error
	// pre load sprites
	if err = eng.LoadSpriteSheet("resources/gopher.json"); err != nil {
		return err
	}

	// get the world
	gw := eng.World()

	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// Get the eyes size
	var eyeSize geometry.Size
	var eyeRadius geometry.Position
	if eyeSize, err = eng.GetSpriteSize("resources/gopher.json", "eye_exterior.png"); err == nil {
		eyeRadius = geometry.Position{X: eyeSize.Width / 4 * gameScale.Point.X, Y: eyeSize.Height / 4 * gameScale.Point.X}
	} else {
		return err
	}

	// the nose is in the middle and a bit down
	nosePos := geometry.Position{
		X: (designResolution.Width / 2) * gameScale.Point.X,
		Y: ((designResolution.Height / 2) + noseVerticalGap) * gameScale.Point.Y,
	}

	// left eye is a bit up left of the nose
	leftEyePos := geometry.Position{
		X: nosePos.X - (eyesGap * gameScale.Point.X),
		Y: nosePos.Y - (eyesGap * gameScale.Point.Y),
	}

	// right eye is a bit up right of the nose
	rightEyePos := geometry.Position{
		X: nosePos.X + (eyesGap * gameScale.Point.Y),
		Y: leftEyePos.Y,
	}

	// add the nose sprite
	gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "nose.png", Scale: gameScale.Min},
		nosePos,
	))

	// add the left exterior eye
	leftExterior = gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior.png", Scale: gameScale.Min},
		leftEyePos,
	))

	// add the lef exterior eye
	gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior.png", Scale: gameScale.Min},
		leftEyePos,
		lookAtMouse{pivot: leftExterior, radius: eyeRadius},
	))

	// add the right exterior eye
	rightExterior = gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior.png", Scale: gameScale.Min},
		rightEyePos,
	))

	// add the right interior eye
	gw.Add(entity.New(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior.png", Scale: gameScale.Min},
		rightEyePos,
		lookAtMouse{pivot: rightExterior, radius: eyeRadius},
	))

	// the text is bottom center
	textPos := geometry.Position{
		X: (designResolution.Width / 2) * gameScale.Point.X,
		Y: designResolution.Height * gameScale.Point.Y,
	}

	// add our text
	gw.Add(entity.New(
		text.Text{
			String:     "press <ESC> to close",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.BottomVAlignment,
			Size:       textSmallSize * gameScale.Min,
			Spacing:    (textSmallSize / 4) * gameScale.Min,
		},
		textPos,
		color.Black,
	))

	// our bar shape
	box := shapes.Box{
		Size: geometry.Size{
			Width:  dizzyBarWith,
			Height: dizzyBarHeight,
		},
		Scale: gameScale.Min,
	}

	// position of the bar
	dizzyBarPosition := geometry.Position{
		X: (designResolution.Width - dizzyBarWith) / 2 * gameScale.Min,
		Y: dizzyBarGap * gameScale.Min,
	}

	// position the dizzy text
	dizzyTextPosition := geometry.Position{
		X: designResolution.Width / 2 * gameScale.Min,
		Y: (dizzyBarGap + (dizzyBarHeight / 2)) * gameScale.Min,
	}

	// Add the dizzy bar
	dizzyBar = gw.Add(entity.New(
		color.Gradient{From: color.Green, To: color.Red},
		box,
		dizzyBarPosition,
	))

	// Add the empty dizzy bar
	dizzyBarEmpty = gw.Add(entity.New(
		color.Black,
		box,
		dizzyBarPosition,
	))

	// add the dizzy text
	dizzyText = gw.Add(entity.New(
		text.Text{
			String:     "Dizzy! Level",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.MiddleVAlignment,
			Size:       textBigSize * gameScale.Min,
			Spacing:    (textBigSize / 4) * gameScale.Min,
		},
		color.Green,
		dizzyTextPosition,
	))

	// add our look at mouse system
	gw.AddSystem(&lookAtMouseSystem{})
	// add our dizzy bar system
	gw.AddSystem(&dizzyBarSystem{})

	return nil
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
type lookAtMouseSystem struct{}

func (lam lookAtMouseSystem) Update(_ *world.World, _ float32) error {
	return nil
}

func (lam *lookAtMouseSystem) Notify(gw *world.World, event interface{}, _ float32) error {
	switch ev := event.(type) {
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

	ax := la.radius.X * float32(math.Cos(float64(angle)))
	ay := la.radius.Y * float32(math.Sin(float64(angle)))

	np := geometry.Position{
		X: pos.X + ax,
		Y: pos.Y + ay,
	}

	ent.Set(np)
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

	// make the dizzy text color change from green to blend depending on how dizzy we are
	cl := color.Green.Blend(color.Red, 1-percent)
	dizzyText.Set(cl)

	// color the eyes red
	clr := color.White.Blend(color.Red, (1-percent)/2)
	leftExterior.Set(clr)
	rightExterior.Set(clr)

	return nil
}

func (dbs *dizzyBarSystem) Notify(_ *world.World, event interface{}, delta float32) error {
	switch event.(type) {
	case events.MouseMoveEvent:
		dbs.dizzy = float32(math.Min(float64(dbs.dizzy+(delta*dizzyGainRate)), 2.5))
	}
	return nil
}
