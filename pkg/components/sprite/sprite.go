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

// Package sprite handle the Sprite component
package sprite

import (
	"github.com/juan-medina/goecs/pkg/entity"
	"reflect"
)

// Sprite is a graphic image that will drawn on the screen
type Sprite struct {
	Sheet    string  // Sheet is the Sprite sheet
	Name     string  // Name is the Sprite name
	Rotation float32 // Rotation is the Sprite rotation
	Scale    float32 // Scale is the Sprite scale
	FlipX    bool    // FlipX indicates if the Sprite is flipped in the X-Assis
	FlipY    bool    // FlipY indicates if the Sprite is flipped in the Y-Assis
}

// TYPE is the reflect.Type of the Sprite
var TYPE = reflect.TypeOf(Sprite{})

// Get gets a Sprite from a entity.Entity
func Get(e *entity.Entity) Sprite {
	return e.Get(TYPE).(Sprite)
}
