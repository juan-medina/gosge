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
)

// ControlState is the state of UI element
type ControlState struct {
	Hover    bool
	Clicked  bool
	Disabled bool
	Checked  bool
	Focused  bool
}

// Type return this goecs.ComponentType
func (c ControlState) Type() goecs.ComponentType {
	return TYPE.ControlState
}

type types struct {
	// FlatButton is the goecs.ComponentType for ui.FlatButton
	FlatButton goecs.ComponentType
	// ButtonColor is the goecs.ComponentType for ui.ButtonColor
	ButtonColor goecs.ComponentType
	// ButtonHoverColors is the goecs.ComponentType for ui.ButtonHoverColors
	ButtonHoverColors goecs.ComponentType
	// SpriteButton is the goecs.ComponentType for ui.SpriteButton
	SpriteButton goecs.ComponentType
	// Text is the goecs.ComponentType for ui.Text
	Text goecs.ComponentType
	// ProgressBar is the goecs.ComponentType for ui.ProgressBar
	ProgressBar goecs.ComponentType
	// ProgressBarColor is the goecs.ComponentType for ui.ProgressBarColor
	ProgressBarColor goecs.ComponentType
	// ProgressBarHoverColor is the goecs.ComponentType for ui.ProgressBarHoverColor
	ProgressBarHoverColor goecs.ComponentType
	// ControlState is the goecs.ComponentType for ui.ControlState
	ControlState goecs.ComponentType
}

// TYPE hold the goecs.ComponentType for our ui components
var TYPE = types{
	FlatButton:            goecs.NewComponentType(),
	ButtonColor:           goecs.NewComponentType(),
	ButtonHoverColors:     goecs.NewComponentType(),
	SpriteButton:          goecs.NewComponentType(),
	Text:                  goecs.NewComponentType(),
	ProgressBar:           goecs.NewComponentType(),
	ProgressBarColor:      goecs.NewComponentType(),
	ProgressBarHoverColor: goecs.NewComponentType(),
	ControlState:          goecs.NewComponentType(),
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
