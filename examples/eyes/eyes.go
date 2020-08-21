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

// layout constants
const (
	noseVerticalGap = 300
	eyesGap         = 400
	textSize        = 40
)

// set the position and scale of the objects
func setPosAndScale(ent *entity.Entity, pos position.Position, scale float64) {
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

// layout our entities
func positionElements(width int, height int, scale float64) {
	// the nose is in the middle and a bit down
	nosePos := position.Position{
		X: float64(width / 2),
		Y: float64(height/2) + (noseVerticalGap * scale),
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
	setPosAndScale(nose, nosePos, scale)
	setPosAndScale(leftExterior, leftEyePos, scale)
	setPosAndScale(leftInterior, leftEyePos, scale)
	setPosAndScale(rightExterior, rightEyePos, scale)
	setPosAndScale(rightInterior, rightEyePos, scale)

	// the text is bottom center
	textPos := position.Position{
		X: float64(width) / 2,
		Y: float64(height),
	}
	// update text pos and scale
	setPosAndScale(bottomText, textPos, scale)
}

type layoutSystem struct{}

func (ls layoutSystem) Update(_ *world.World, _ float64) error {
	return nil
}

func (ls layoutSystem) Notify(_ *world.World, event interface{}, _ float64) error {
	switch ev := event.(type) {
	// if the screen has change size
	case events.ScreenSizeChangeEvent:
		// change layout
		positionElements(ev.Current.Width, ev.Current.Height, ev.Scale.Max)
	}

	return nil
}

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
	leftInterior = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior"}))
	rightExterior = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_exterior"}))
	rightInterior = gw.Add(entity.New(sprite.Sprite{Sheet: "resources/gopher.json", Name: "eye_interior"}))

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
	gw.AddSystem(layoutSystem{})

	return nil
}

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}
