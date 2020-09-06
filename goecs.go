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

// Package gosge is a opinionated 2D only game engine that use an ECS
//
// gosge uses https://github.com/juan-medina/goecs) and https://github.com/gen2brain/raylib-go, the go port of
// https://www.raylib.com/ for rendering and device manager
package gosge

import (
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// Run the game given the game options.Options and the loading InitFunc
func Run(opt options.Options, init InitFunc) error {
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	log.Logger = log.Output(writer).Level(zerolog.TraceLevel)
	return New(opt, init).Run()
}
