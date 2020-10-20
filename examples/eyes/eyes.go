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
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
	"math"
)

// our game options
var opt = options.Options{
	Title:      "GOSGE Eyes Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed: true,
	// Width:    2048,
	// Height:   1536,
}

// entities that we are going to create
var (
	leftExterior  goecs.EntityID
	rightExterior goecs.EntityID
	dizzyBar      goecs.EntityID
	dizzyText     goecs.EntityID

	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}

	// dizzy how much dizzy we are
	dizzy = float32(0)

	// lastMousePos is the last mouse position
	lastMousePos geometry.Point

	// dizzyBarWith is the dizzy bar with
	dizzyBarWith float32
)

// game constants
const (
	noseVerticalGap = 300
	eyesGap         = 400
	fontName        = "resources/go_regular.fnt"
	textSmallSize   = 60
	textBigSize     = 100
	dizzyBarHeight  = 110
	dizzyBarGap     = 50
	maxDizzy        = 2.5
	dizzyGainRate   = 1.0
	dizzyReduceRate = 0.45
)

func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng *gosge.Engine) error {
	var err error

	// Preload font
	if err = eng.LoadFont(fontName); err != nil {
		return err
	}

	// pre load sprites
	if err = eng.LoadSpriteSheet("resources/gopher.json"); err != nil {
		return err
	}

	// get the world
	world := eng.World()

	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// Get the eyes size
	var eyeSize geometry.Size
	var eyeRadius geometry.Point
	if eyeSize, err = eng.GetSpriteSize("resources/gopher.json", "eye_exterior.png"); err == nil {
		eyeRadius = geometry.Point{X: eyeSize.Width / 4 * gameScale.Point.X, Y: eyeSize.Height / 4 * gameScale.Point.X}
	} else {
		return err
	}

	// the nose is in the middle and a bit down
	nosePos := geometry.Point{
		X: (designResolution.Width / 2) * gameScale.Point.X,
		Y: ((designResolution.Height / 2) + noseVerticalGap) * gameScale.Point.Y,
	}

	// left eye is a bit up left of the nose
	leftEyePos := geometry.Point{
		X: nosePos.X - (eyesGap * gameScale.Point.X),
		Y: nosePos.Y - (eyesGap * gameScale.Point.Y),
	}

	// right eye is a bit up right of the nose
	rightEyePos := geometry.Point{
		X: nosePos.X + (eyesGap * gameScale.Point.Y),
		Y: leftEyePos.Y,
	}

	// add the nose sprite
	world.AddEntity(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "nose.png", Scale: gameScale.Max},
		nosePos,
	)

	// add the left exterior eye
	leftExterior = world.AddEntity(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior.png", Scale: gameScale.Max},
		leftEyePos,
	)

	// add the lef exterior eye
	world.AddEntity(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior.png", Scale: gameScale.Max},
		leftEyePos,
		lookAtMouse{pivot: leftExterior, radius: eyeRadius},
	)

	// add the right exterior eye
	rightExterior = world.AddEntity(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior.png", Scale: gameScale.Max},
		rightEyePos,
	)

	// add the right interior eye
	world.AddEntity(
		sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior.png", Scale: gameScale.Max},
		rightEyePos,
		lookAtMouse{pivot: rightExterior, radius: eyeRadius},
	)

	// the text is bottom center
	textPos := geometry.Point{
		X: (designResolution.Width / 2) * gameScale.Point.X,
		Y: designResolution.Height * gameScale.Point.Y,
	}

	// add our text
	world.AddEntity(
		ui.Text{
			String:     "press <ESC> to close",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       textSmallSize * gameScale.Max,
		},
		textPos,
		color.White,
	)

	ss := eng.GetScreenSize()

	dizzyBarWith = ss.Width / gameScale.Max * 0.75

	// Point of the bar
	dizzyBarPoint := geometry.Point{
		X: ((designResolution.Width * gameScale.Point.X) - (dizzyBarWith * gameScale.Max)) / 2,
		Y: dizzyBarGap * gameScale.Max,
	}

	// Point the dizzy text
	dizzyTextPoint := geometry.Point{
		X: designResolution.Width / 2 * gameScale.Point.X,
		Y: (dizzyBarGap + (dizzyBarHeight / 2)) * gameScale.Max,
	}

	// add the bar
	dizzyBar = world.AddEntity(
		ui.ProgressBar{
			Min:     0,
			Max:     maxDizzy,
			Current: maxDizzy,
			Shadow: geometry.Size{
				Width:  5 * gameScale.Max,
				Height: 5 * gameScale.Max,
			},
		},
		dizzyBarPoint,
		shapes.Box{
			Size: geometry.Size{
				Width:  dizzyBarWith,
				Height: dizzyBarHeight,
			},
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		ui.ProgressBarColor{
			Gradient: color.Gradient{
				From:      color.Green,
				To:        color.Red,
				Direction: color.GradientHorizontal,
			},
			Border: color.White,
			Empty:  color.Black,
		},
	)

	// add the dizzy text
	dizzyText = world.AddEntity(
		ui.Text{
			String:     "Dizzy! Level",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.MiddleVAlignment,
			Font:       fontName,
			Size:       textBigSize * gameScale.Max,
		},
		color.Green,
		dizzyTextPoint,
	)

	// listen to mouse moves
	world.AddListener(mouseMoveListener, events.TYPE.MouseMoveEvent)
	// add the system that decrease how dizzy we are
	world.AddSystem(decreaseDizzySystem)
	// add our dizzy bar system
	world.AddSystem(updateDizzyBarSystem)

	// set last mouse pos
	lastMousePos = geometry.Point{
		X: -1,
		Y: -1,
	}

	return nil
}

// component to make an entity to look at mouse with a pivot
type lookAtMouse struct {
	pivot  goecs.EntityID
	radius geometry.Point
}

func (l lookAtMouse) Type() goecs.ComponentType {
	return types.lookAtMouse
}

var types = struct{ lookAtMouse goecs.ComponentType }{lookAtMouse: goecs.NewComponentType()}

func getLookAtMouse(e *goecs.Entity) lookAtMouse {
	return e.Get(types.lookAtMouse).(lookAtMouse)
}

func mouseMoveListener(gw *goecs.World, event goecs.Component, delta float32) error {
	switch ev := event.(type) {
	// if we move the mouse
	case events.MouseMoveEvent:
		// get the entities that look at the mouse
		for it := gw.Iterator(types.lookAtMouse); it != nil; it = it.Next() {
			v := it.Value()
			la := getLookAtMouse(v)
			// make this entity to look at the mouse
			lookAt(gw, v, la, ev.Point)
		}
		pos := ev.Point
		if lastMousePos.X == -1 && lastMousePos.Y == -1 {
			lastMousePos = pos
		} else {
			if lastMousePos.X != pos.X || lastMousePos.Y != lastMousePos.Y {
				// increase how dizzy we are
				dizzy = float32(math.Min(float64(dizzy+(delta*dizzyGainRate)), 2.5))
				lastMousePos = pos
			}
		}
	}
	return nil
}

func lookAt(world *goecs.World, ent *goecs.Entity, la lookAtMouse, mouse geometry.Point) {
	pivotEnt := world.Get(la.pivot)
	pos := geometry.Get.Point(pivotEnt)

	dx := mouse.X - pos.X
	dy := mouse.Y - pos.Y

	angle := float32(math.Atan2(float64(dy), float64(dx)))

	ax := la.radius.X * float32(math.Cos(float64(angle)))
	ay := la.radius.Y * float32(math.Sin(float64(angle)))

	np := geometry.Point{
		X: pos.X + ax,
		Y: pos.Y + ay,
	}

	ent.Set(np)
}

func decreaseDizzySystem(_ *goecs.World, delta float32) error {
	// in each frame we reduce how dizzy we are
	dizzy = float32(math.Max(float64(dizzy-(delta*dizzyReduceRate)), 0))
	return nil
}

func updateDizzyBarSystem(world *goecs.World, _ float32) error {
	dizzyBarEnt := world.Get(dizzyBar)
	bar := ui.Get.ProgressBar(dizzyBarEnt)
	bar.Current = dizzy
	dizzyBarEnt.Set(bar)

	// calculate how dizzy we are in 0..1
	percent := 1 - (dizzy / maxDizzy)

	// make the dizzy text color change from green to blend depending on how dizzy we are
	cl := color.Green.Blend(color.Red, 1-percent)
	dizzyTextEnt := world.Get(dizzyText)
	dizzyTextEnt.Set(cl)

	// color the eyes red
	clr := color.White.Blend(color.Red, (1-percent)/2)
	leftExteriorEnt := world.Get(leftExterior)
	rightExteriorEnt := world.Get(rightExterior)
	leftExteriorEnt.Set(clr)
	rightExteriorEnt.Set(clr)

	return nil
}
