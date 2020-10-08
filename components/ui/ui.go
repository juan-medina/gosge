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
	"github.com/juan-medina/goecs"
	"reflect"
)

// ControlState is the state of UI element
type ControlState struct {
	Hover    bool
	Clicked  bool
	Disabled bool
	Checked  bool
}

type types struct {
	// FlatButton is the reflect.Type for ui.FlatButton
	FlatButton reflect.Type
	// ButtonColor is the reflect.Type for ui.ButtonColor
	ButtonColor reflect.Type
	// ButtonHoverColors is the reflect.Type for ui.ButtonHoverColors
	ButtonHoverColors reflect.Type
	// SpriteButton is the reflect.Type for ui.SpriteButton
	SpriteButton reflect.Type
	// Text is the reflect.Type for ui.Text
	Text reflect.Type
	// ProgressBar is the reflect.Type for ui.ProgressBar
	ProgressBar reflect.Type
	// ProgressBarColor is the reflect.Type for ui.ProgressBarColor
	ProgressBarColor reflect.Type
	// ProgressBarHoverColor is the reflect.Type for ui.ProgressBarHoverColor
	ProgressBarHoverColor reflect.Type
	// ControlState is the reflect.Type for ui.ControlState
	ControlState reflect.Type
}

// TYPE hold the reflect.Type for our ui components
var TYPE = types{
	FlatButton:            reflect.TypeOf(FlatButton{}),
	ButtonColor:           reflect.TypeOf(ButtonColor{}),
	ButtonHoverColors:     reflect.TypeOf(ButtonHoverColors{}),
	SpriteButton:          reflect.TypeOf(SpriteButton{}),
	Text:                  reflect.TypeOf(Text{}),
	ProgressBar:           reflect.TypeOf(ProgressBar{}),
	ProgressBarColor:      reflect.TypeOf(ProgressBarColor{}),
	ProgressBarHoverColor: reflect.TypeOf(ProgressBarHoverColor{}),
	ControlState:          reflect.TypeOf(ControlState{}),
}

type gets struct {
	// FlatButton gets a ui.FlatButton from a goecs.Entity
	FlatButton func(e *goecs.Entity) FlatButton
	// ButtonColor gets a ui.ButtonColor from a goecs.Entity
	ButtonColor func(e *goecs.Entity) ButtonColor
	// ButtonHoverColors gets a ui.ButtonHoverColors from a goecs.Entity
	ButtonHoverColors func(e *goecs.Entity) ButtonHoverColors
	// SpriteButton gets a ui.SpriteButton from a goecs.Entity
	SpriteButton func(e *goecs.Entity) SpriteButton
	// Text gets a ui.Text from a goecs.Entity
	Text func(e *goecs.Entity) Text
	// ProgressBar gets a ui.ProgressBar from a goecs.Entity
	ProgressBar func(e *goecs.Entity) ProgressBar
	// ProgressBarColor gets a ui.ProgressBarColor from a goecs.Entity
	ProgressBarColor func(e *goecs.Entity) ProgressBarColor
	// ProgressBarHoverColor gets a ui.ProgressBarHoverColor from a goecs.Entity
	ProgressBarHoverColor func(e *goecs.Entity) ProgressBarHoverColor
	// ControlState gets a ui.ControlState from a goecs.Entity
	ControlState func(e *goecs.Entity) ControlState
}

// Get a ui component
var Get = gets{
	// FlatButton gets a ui.FlatButton from a goecs.Entity
	FlatButton: func(e *goecs.Entity) FlatButton {
		return e.Get(TYPE.FlatButton).(FlatButton)
	},
	// ButtonColor gets a ui.ButtonColor from a goecs.Entity
	ButtonColor: func(e *goecs.Entity) ButtonColor {
		return e.Get(TYPE.ButtonColor).(ButtonColor)
	},
	// ButtonHoverColors gets a ui.ButtonHoverColors from a goecs.Entity
	ButtonHoverColors: func(e *goecs.Entity) ButtonHoverColors {
		return e.Get(TYPE.ButtonHoverColors).(ButtonHoverColors)
	},
	// SpriteButton gets a ui.SpriteButton from a goecs.Entity
	SpriteButton: func(e *goecs.Entity) SpriteButton {
		return e.Get(TYPE.SpriteButton).(SpriteButton)
	},
	// Text gets a ui.Text from a goecs.Entity
	Text: func(e *goecs.Entity) Text {
		return e.Get(TYPE.Text).(Text)
	},
	// ProgressBar gets a ui.ProgressBar from a goecs.Entity
	ProgressBar: func(e *goecs.Entity) ProgressBar {
		return e.Get(TYPE.ProgressBar).(ProgressBar)
	},
	// ProgressBarColor gets a ui.ProgressBarColor from a goecs.Entity
	ProgressBarColor: func(e *goecs.Entity) ProgressBarColor {
		return e.Get(TYPE.ProgressBarColor).(ProgressBarColor)
	},
	// ProgressBarHoverColor gets a ui.ProgressBarHoverColor from a goecs.Entity
	ProgressBarHoverColor: func(e *goecs.Entity) ProgressBarHoverColor {
		return e.Get(TYPE.ProgressBarHoverColor).(ProgressBarHoverColor)
	},
	// ControlState gets a ui.ControlState from a goecs.Entity
	ControlState: func(e *goecs.Entity) ControlState {
		return e.Get(TYPE.ControlState).(ControlState)
	},
}
