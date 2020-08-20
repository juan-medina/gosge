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
	"github.com/juan-medina/goecs/pkg/entitiy"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/position"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
)

var opt = options.Options{
	Title:      "Eyes Game",
	Width:      1920,
	Height:     1080,
	ClearColor: color.Gopher,
}

func loadGame(eng engine.Engine) error {
	if err := eng.LoadTexture("resources/nose.png"); err != nil {
		return err
	}
	if err := eng.LoadTexture("resources/eye_exterior.png"); err != nil {
		return err
	}
	if err := eng.LoadTexture("resources/eye_interior.png"); err != nil {
		return err
	}

	nosePos := position.Position{
		X: float64(opt.Width / 2),
		Y: float64(opt.Height - 300),
	}
	leftEyePos := position.Position{
		X: nosePos.X - 400,
		Y: nosePos.Y - 400,
	}
	rightEyePos := position.Position{
		X: nosePos.X + 400,
		Y: leftEyePos.Y,
	}

	world := eng.World()
	world.Add(entitiy.New(
		sprite.Sprite{
			FileName: "resources/nose.png",
			Scale:    1,
			Rotation: 0,
		},
		nosePos,
	))

	world.Add(entitiy.New(
		sprite.Sprite{
			FileName: "resources/eye_exterior.png",
			Scale:    1,
			Rotation: 0,
		},
		leftEyePos,
	))

	world.Add(entitiy.New(
		sprite.Sprite{
			FileName: "resources/eye_interior.png",
			Scale:    1,
			Rotation: 0,
		},
		leftEyePos,
	))

	world.Add(entitiy.New(
		sprite.Sprite{
			FileName: "resources/eye_exterior.png",
			Scale:    1,
			Rotation: 0,
		},
		rightEyePos,
	))

	world.Add(entitiy.New(
		sprite.Sprite{
			FileName: "resources/eye_interior.png",
			Scale:    1,
			Rotation: 0,
		},
		rightEyePos,
	))

	return nil
}

func main() {
	if err := game.Run(opt, loadGame); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}
