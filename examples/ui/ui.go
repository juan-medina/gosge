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
	"fmt"
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/events"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
	"strconv"
)

var opt = options.Options{
	Title:      "GOSGE UI Game",
	BackGround: color.Gopher,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed:   true,
	// Width:      2048,
	// Height:     1536,
}

const (
	fontName    = "resources/go_regular.fnt"
	fontSmall   = 25
	fontBig     = 40
	columnGap   = 350
	finalColumn = 900
	rowGap      = 50
	spriteSheet = "resources/ui.json"
	spriteScale = 0.25
	hint        = "press <F1> to hide/unhide, <F2> to disable/enable, <ESC> to close"
)

var (
	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
	message          goecs.EntityID
	gEng             *gosge.Engine
)

func main() {
	if err := gosge.Run(opt, loadGame); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
	}
}

func loadGame(eng *gosge.Engine) error {
	var err error
	gEng = eng

	// Preload font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	// Preload sprite sheet
	if err := eng.LoadSpriteSheet(spriteSheet); err != nil {
		return err
	}

	world := eng.World()

	// gameScale has a geometry.Scale from the real screen size to our designResolution
	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// element pos
	pos := geometry.Point{
		X: 10 * gameScale.Max,
		Y: 10 * gameScale.Max,
	}

	// add the flat button
	addFlatButton(world, gameScale, pos, false, true)
	pos.Y += rowGap * gameScale.Max
	addFlatButton(world, gameScale, pos, true, false)

	// add check boxes
	pos.Y += rowGap * gameScale.Max
	addCheckBox(world, gameScale, pos, false)
	pos.Y += rowGap * gameScale.Max
	addCheckBox(world, gameScale, pos, true)

	// add option group
	pos.Y += rowGap * gameScale.Max
	addOptionGroup(world, gameScale, pos, false)
	pos.Y += rowGap * gameScale.Max
	addOptionGroup(world, gameScale, pos, true)

	// add the progress bar
	pos.Y += rowGap * gameScale.Max
	addProgressBar(world, gameScale, pos, false, true)
	pos.Y += rowGap * gameScale.Max
	addProgressBar(world, gameScale, pos, true, true)
	pos.Y += rowGap * gameScale.Max
	addProgressBar(world, gameScale, pos, false, false)
	pos.Y += rowGap * gameScale.Max
	addProgressBar(world, gameScale, pos, true, false)

	// add sprite button
	pos.Y += rowGap * gameScale.Max
	if err := addSpriteButton(world, gameScale, pos); err != nil {
		return err
	}

	// add the bottom text
	message = world.AddEntity(
		ui.Text{
			String:     "",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       fontBig * gameScale.Max,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height / 2 * gameScale.Point.Y,
		},
		color.White,
	)

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     hint,
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       fontBig * gameScale.Max,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height * gameScale.Point.Y,
		},
		effects.AlternateColor{
			Time: .25,
			From: color.White,
			To:   color.White.Alpha(0),
		},
	)

	// listen tu ui events
	world.AddListener(uiEvents, uiDemoEventType, progressBarEventType, checkBoxEventType, optionEventType)

	// listen to keys
	world.AddListener(keyEvents, events.TYPE.KeyUpEvent)

	// listen to game pad
	world.AddListener(padEvents, events.TYPE.GamePadButtonUpEvent)

	return err
}

func padEvents(world *goecs.World, signal goecs.Component, _ float32) error {
	switch v := signal.(type) {
	case events.GamePadButtonUpEvent:
		switch v.Button {
		case device.GamepadStart:
			hideUnhide(world)
		case device.GamepadSelect:
			disableEnable(world)
		}
	}
	return nil
}

func keyEvents(world *goecs.World, signal goecs.Component, _ float32) error {
	switch e := signal.(type) {
	case events.KeyUpEvent:
		switch e.Key {
		case device.KeyF1:
			hideUnhide(world)
		case device.KeyF2:
			disableEnable(world)
		}
	}
	return nil
}

func disableEnable(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.ControlState); it != nil; it = it.Next() {
		ent := it.Value()
		state := ui.Get.ControlState(ent)
		state.Disabled = !state.Disabled
		ent.Set(state)
	}
}

func hideUnhide(world *goecs.World) {
	for it := world.Iterator(ui.TYPE.ControlState); it != nil; it = it.Next() {
		ent := it.Value()
		if ent.Contains(effects.TYPE.Hide) {
			ent.Remove(effects.TYPE.Hide)
		} else {
			ent.Add(effects.Hide{})
		}
	}
}

func addOptionGroup(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point, gradient bool) {
	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max),
		Y: labelPos.Y,
	}

	text := "Group [Solid]"
	group := "A"

	clr := ui.ButtonColor{
		Solid:  color.SkyBlue,
		Border: color.White,
		Text:   color.White,
	}

	if gradient {
		clr.Gradient = color.Gradient{
			From:      color.SkyBlue,
			To:        color.DarkBlue,
			Direction: color.GradientHorizontal,
		}
		text = "Group [Gradient]"
		group = "B"
	}

	// add a label
	world.AddEntity(
		ui.Text{
			String:     text,
			HAlignment: ui.LeftHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		labelPos,
		color.White,
	)

	valuePos := geometry.Point{
		X: finalColumn * gameScale.Max,
		Y: controlPos.Y,
	}

	valueID := world.AddEntity(
		ui.Text{
			String:     group + " 1",
			Size:       fontSmall * gameScale.Max,
			Font:       fontName,
			VAlignment: ui.TopVAlignment,
			HAlignment: ui.LeftHAlignment,
		},
		valuePos,
		color.White,
	)

	for c := 0; c < 3; c++ {
		checked := c == 0
		check := ui.FlatButton{
			Shadow: geometry.Size{
				Width:  2 * gameScale.Max,
				Height: 2 * gameScale.Max,
			},
			CheckBox: true,
			Group:    group,
		}

		// add a control : flat button with checkbox
		checkID := world.AddEntity(
			check,
			clr,
			ui.Text{
				String:     "   " + group + "  " + strconv.Itoa(c+1),
				HAlignment: ui.CenterHAlignment,
				VAlignment: ui.MiddleVAlignment,
				Font:       fontName,
				Size:       fontSmall * gameScale.Max,
			},
			shapes.Box{
				Size: geometry.Size{
					Width:  70 * gameScale.Max,
					Height: 20 * gameScale.Max,
				},
				Scale:     gameScale.Max,
				Thickness: int32(2 * gameScale.Max),
			},
			ui.ControlState{
				Checked: checked,
			},
			controlPos,
		)

		check.Event = optionEvent{
			valueID: valueID,
			Message: text + " clicked",
			value:   group + " " + strconv.Itoa(c+1),
		}
		checkEnt := world.Get(checkID)
		checkEnt.Set(check)
		controlPos.X += 150 * gameScale.Max
	}
}

func addCheckBox(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point, gradient bool) {
	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max),
		Y: labelPos.Y,
	}

	text := "Check [Solid]"

	clr := ui.ButtonColor{
		Solid:  color.SkyBlue,
		Border: color.White,
		Text:   color.White,
	}

	if gradient {
		clr.Gradient = color.Gradient{
			From:      color.SkyBlue,
			To:        color.DarkBlue,
			Direction: color.GradientHorizontal,
		}
		text = "Check [Gradient]"
	}

	// add a label
	world.AddEntity(
		ui.Text{
			String:     text,
			HAlignment: ui.LeftHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		labelPos,
		color.White,
	)

	check := ui.FlatButton{
		Shadow: geometry.Size{
			Width:  2 * gameScale.Max,
			Height: 2 * gameScale.Max,
		},
		CheckBox: true,
	}

	// add a control : flat button with checkbox
	checkID := world.AddEntity(
		check,
		clr,
		ui.Text{
			String:     "   Check",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.MiddleVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  70 * gameScale.Max,
				Height: 20 * gameScale.Max,
			},
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		controlPos,
	)

	controlPos.X = finalColumn * gameScale.Max

	valueID := world.AddEntity(
		ui.Text{
			String:     "Not checked",
			Size:       fontSmall * gameScale.Max,
			Font:       fontName,
			VAlignment: ui.TopVAlignment,
			HAlignment: ui.LeftHAlignment,
		},
		controlPos,
		color.White,
	)

	check.Event = checkBoxEvent{
		checkID: checkID,
		valueID: valueID,
		Message: text + " clicked",
	}

	checkEnt := world.Get(checkID)
	checkEnt.Set(check)
}

func addFlatButton(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point, gradient bool, focus bool) {
	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max),
		Y: labelPos.Y,
	}

	text := "FlatButton [Solid]"

	clr := ui.ButtonColor{
		Solid:  color.SkyBlue,
		Border: color.White,
		Text:   color.White,
	}

	if gradient {
		clr.Gradient = color.Gradient{
			From:      color.SkyBlue,
			To:        color.DarkBlue,
			Direction: color.GradientHorizontal,
		}
		text = "FlatButton [Gradient]"
	}

	// add a label
	world.AddEntity(
		ui.Text{
			String:     text,
			HAlignment: ui.LeftHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		labelPos,
		color.White,
	)

	// add a control : flat button
	ctl := world.AddEntity(
		ui.FlatButton{
			Shadow: geometry.Size{
				Width:  2 * gameScale.Max,
				Height: 2 * gameScale.Max,
			},
			Event: uiDemoEvent{Message: text + " clicked"},
		},
		clr,
		ui.Text{
			String:     "Click Me",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.MiddleVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		shapes.Box{
			Size: geometry.Size{
				Width:  70 * gameScale.Max,
				Height: 20 * gameScale.Max,
			},
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		controlPos,
	)
	if focus {
		world.Signal(events.FocusOnControlEvent{Control: ctl})
	}
}

func addProgressBar(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point, gradient bool, focusable bool) {
	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max),
		Y: labelPos.Y,
	}

	text := "ProgressBar [Solid]"

	clr := ui.ProgressBarColor{
		Solid:  color.SkyBlue,
		Border: color.White,
		Empty:  color.Blue.Blend(color.White, 0.65),
	}

	if gradient {
		clr.Gradient = color.Gradient{
			From:      color.SkyBlue,
			To:        color.DarkBlue,
			Direction: color.GradientHorizontal,
		}
		text = "ProgressBar [Gradient]"
	}

	if !focusable {
		text = text + " [No focus]"
	}

	// add a label
	world.AddEntity(
		ui.Text{
			String:     text,
			HAlignment: ui.LeftHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		labelPos,
		color.White,
	)

	bar := ui.ProgressBar{
		Min:     0,
		Max:     100,
		Current: 50,
		Shadow: geometry.Size{
			Width:  2 * gameScale.Max,
			Height: 2 * gameScale.Max,
		},
	}

	barID := world.AddEntity(
		bar,
		shapes.Box{
			Size: geometry.Size{
				Width:  220 * gameScale.Max,
				Height: 20 * gameScale.Max,
			},
			Scale:     gameScale.Max,
			Thickness: int32(2 * gameScale.Max),
		},
		clr,
		controlPos,
	)

	controlPos.X = finalColumn * gameScale.Max

	valueID := world.AddEntity(
		ui.Text{
			String:     "50",
			Size:       fontSmall * gameScale.Max,
			Font:       fontName,
			VAlignment: ui.TopVAlignment,
			HAlignment: ui.LeftHAlignment,
		},
		controlPos,
		color.White,
	)

	if focusable {
		bar.Event = progressBarEvent{
			barID:   barID,
			valueID: valueID,
			Message: text + " clicked",
		}
	}

	barEnt := world.Get(barID)
	barEnt.Set(bar)
}

func addSpriteButton(world *goecs.World, gameScale geometry.Scale, labelPos geometry.Point) error {
	text := "SpriteButton"

	// add a label
	world.AddEntity(
		ui.Text{
			String:     text,
			HAlignment: ui.LeftHAlignment,
			VAlignment: ui.TopVAlignment,
			Font:       fontName,
			Size:       fontSmall * gameScale.Max,
		},
		labelPos,
		color.White,
	)
	var err error
	var size geometry.Size
	if size, err = gEng.GetSpriteSize(spriteSheet, "normal.png"); err != nil {
		return err
	}

	// control pos
	controlPos := geometry.Point{
		X: labelPos.X + (columnGap * gameScale.Max) + (size.Width * 0.5 * gameScale.Max * spriteScale),
		Y: labelPos.Y + (size.Height * 0.5 * gameScale.Max * spriteScale),
	}

	// add the sprite button
	world.AddEntity(
		ui.SpriteButton{
			Sheet:    spriteSheet,
			Normal:   "normal.png",
			Hover:    "hover.png",
			Clicked:  "click.png",
			Disabled: "locked.png",
			Focused:  "focused.png",
			Scale:    gameScale.Max * spriteScale,
			Event:    uiDemoEvent{Message: text + " clicked"},
		},
		controlPos,
	)
	return nil
}

func uiEvents(world *goecs.World, signal goecs.Component, _ float32) error {
	var err error
	messageEnt := world.Get(message)
	switch e := signal.(type) {
	case uiDemoEvent:
		text := ui.Get.Text(messageEnt)
		text.String = e.Message
		messageEnt.Set(text)
	case progressBarEvent:
		barEnt := world.Get(e.barID)
		labelEnt := world.Get(e.valueID)
		bar := ui.Get.ProgressBar(barEnt)
		label := ui.Get.Text(labelEnt)
		label.String = fmt.Sprintf("%d", int(bar.Current))
		labelEnt.Set(label)
		text := ui.Get.Text(messageEnt)
		text.String = e.Message
		messageEnt.Set(text)
	case checkBoxEvent:
		valueEnt := world.Get(e.valueID)
		checkEnt := world.Get(e.checkID)
		label := ui.Get.Text(valueEnt)
		state := ui.Get.ControlState(checkEnt)
		if state.Checked {
			label.String = "Checked"
		} else {
			label.String = "Not checked"
		}
		valueEnt.Set(label)
		text := ui.Get.Text(messageEnt)
		text.String = e.Message
		messageEnt.Set(text)
	case optionEvent:
		valueEnt := world.Get(e.valueID)
		label := ui.Get.Text(valueEnt)
		label.String = e.value

		valueEnt.Set(label)
		text := ui.Get.Text(messageEnt)
		text.String = e.Message
		messageEnt.Set(text)
	}
	return err
}

type uiDemoEvent struct {
	Message string
}

func (u uiDemoEvent) Type() goecs.ComponentType {
	return uiDemoEventType
}

var uiDemoEventType = goecs.NewComponentType()

type progressBarEvent struct {
	Message string
	barID   goecs.EntityID
	valueID goecs.EntityID
}

func (p progressBarEvent) Type() goecs.ComponentType {
	return progressBarEventType
}

var progressBarEventType = goecs.NewComponentType()

type checkBoxEvent struct {
	Message string
	checkID goecs.EntityID
	valueID goecs.EntityID
}

func (c checkBoxEvent) Type() goecs.ComponentType {
	return checkBoxEventType
}

var checkBoxEventType = goecs.NewComponentType()

type optionEvent struct {
	Message string
	value   string
	valueID goecs.EntityID
}

func (o optionEvent) Type() goecs.ComponentType {
	return optionEventType
}

var optionEventType = goecs.NewComponentType()
