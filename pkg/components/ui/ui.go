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

// Package ui contains the components for displaying ui elements
package ui

import (
	"github.com/juan-medina/goecs/pkg/entity"
	"reflect"
)

type types struct {
	// FlatButton is the reflect.Type for ui.FlatButton
	FlatButton reflect.Type
	// Text is the reflect.Type for ui.Text
	Text reflect.Type
}

// TYPE hold the reflect.Type for our ui components
var TYPE = types{
	FlatButton: reflect.TypeOf(FlatButton{}),
	Text:       reflect.TypeOf(Text{}),
}

type gets struct {
	// FlatButton gets a ui.FlatButton from a entity.Entity
	FlatButton func(e *entity.Entity) FlatButton
	// FlatButton gets a ui.FlatButton from a entity.Entity
	Text func(e *entity.Entity) Text
}

// Get a ui component
var Get = gets{
	// FlatButton gets a ui.FlatButton from a entity.Entity
	FlatButton: func(e *entity.Entity) FlatButton {
		return e.Get(TYPE.FlatButton).(FlatButton)
	},
	// Text gets a ui.Text from a entity.Entity
	Text: func(e *entity.Entity) Text {
		return e.Get(TYPE.Text).(Text)
	},
}
