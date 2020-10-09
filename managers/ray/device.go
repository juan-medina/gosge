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

package ray

import (
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/juan-medina/gosge/components/device"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
)

// Init the rendering device
func (dmi *DeviceManagerImpl) Init(opt options.Options) {
	dmi.saveOpts = opt
	rl.SetTraceLogCallback(func(logType int, str string) {
		switch logType {
		case rl.LogDebug:
			log.Debug().Msg(str)
		case rl.LogError:
			log.Error().Msg(str)
		case rl.LogInfo:
			log.Info().Msg(str)
		case rl.LogTrace:
			log.Trace().Msg(str)
		case rl.LogWarning:
			log.Warn().Msg(str)
		case rl.LogFatal:
			log.Fatal().Msg(str)
		default:
			log.Trace().Msg(str)
		}
	})

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(int32(opt.Width), int32(opt.Height), opt.Title)

	if opt.Icon != "" {
		if file, err := os.Open(opt.Icon); err == nil {
			_ = file.Close()
			img := rl.LoadImage(opt.Icon)
			rl.SetWindowIcon(*img)
			rl.UnloadImage(img)
		}
	}

	if !opt.Windowed {
		// TODO: Open a Issue in raylib repo
		// rl.ToggleFullscreen does not work ok on linux
		if runtime.GOOS != "linux" {
			rl.ToggleFullscreen()
			rl.SetWindowMonitor(opt.Monitor)
		}
	}
	rl.InitAudioDevice()
}

// End the rendering device
func (dmi DeviceManagerImpl) End() {
	dmi.StopAllSounds()
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

// BeginFrame for rendering
func (dmi DeviceManagerImpl) BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(dmi.color2RayColor(dmi.saveOpts.BackGround))
}

// EndFrame for rendering
func (dmi DeviceManagerImpl) EndFrame() {
	rl.EndDrawing()
}

// ShouldClose returns if th engine should close
func (dmi DeviceManagerImpl) ShouldClose() bool {
	return rl.WindowShouldClose()
}

var engineKeyToRayKey = map[device.Key]int32{
	device.KeyLeft:      rl.KeyLeft,
	device.KeyRight:     rl.KeyRight,
	device.KeyUp:        rl.KeyUp,
	device.KeyDown:      rl.KeyDown,
	device.KeySpace:     rl.KeySpace,
	device.KeyAltLeft:   rl.KeyLeftAlt,
	device.KeyCtrlLeft:  rl.KeyLeftControl,
	device.KeyAltRight:  rl.KeyRightAlt,
	device.KeyCtrlRight: rl.KeyRightControl,
	device.KeyEscape:    rl.KeyEscape,
	device.KeyReturn:    rl.KeyEnter,
	device.KeyF1:        rl.KeyF1,
	device.KeyF2:        rl.KeyF2,
	device.KeyF3:        rl.KeyF3,
	device.KeyF4:        rl.KeyF4,
	device.KeyF5:        rl.KeyF5,
	device.KeyF6:        rl.KeyF6,
	device.KeyF7:        rl.KeyF7,
	device.KeyF8:        rl.KeyF8,
	device.KeyF9:        rl.KeyF9,
	device.KeyF10:       rl.KeyF10,
	device.KeyF11:       rl.KeyF11,
	device.KeyF12:       rl.KeyF12,
}

// raylib axis
const (
	GamepadAxisLeftX  = 0 // Left stick
	GamepadAxisLeftY  = 1 // Left stick
	GamepadAxisRightX = 2 // Right stick
	GamepadAxisRightY = 3 // Right stick
)

var engineStickToRayXAxis = map[device.GamepadStick]int32{
	device.GamepadLeftStick:  GamepadAxisLeftX,
	device.GamepadRightStick: GamepadAxisRightX,
}

var engineStickToRayYAxis = map[device.GamepadStick]int32{
	device.GamepadLeftStick:  GamepadAxisLeftY,
	device.GamepadRightStick: GamepadAxisRightY,
}

// IsKeyPressed returns if given device.Key is pressed
func (dmi DeviceManagerImpl) IsKeyPressed(key device.Key) bool {
	if v, ok := engineKeyToRayKey[key]; ok {
		return rl.IsKeyPressed(v)
	}
	return false
}

// IsKeyReleased returns if given device.Key is released
func (dmi DeviceManagerImpl) IsKeyReleased(key device.Key) bool {
	if v, ok := engineKeyToRayKey[key]; ok {
		return rl.IsKeyReleased(v)
	}
	return false
}

// SetExitKey will set the exit key
func (dmi DeviceManagerImpl) SetExitKey(key device.Key) {
	rl.SetExitKey(engineKeyToRayKey[key])
}

// DisableExitKey will disable the exit key
func (dmi DeviceManagerImpl) DisableExitKey() {
	rl.SetExitKey(0)
}

// IsGamepadAvailable indicates if the game pad number is available
func (dmi DeviceManagerImpl) IsGamepadAvailable(gamePad int32) bool {
	return rl.IsGamepadAvailable(gamePad)
}

// IsGamepadButtonPressed returns if given gamepad button is pressed
func (dmi DeviceManagerImpl) IsGamepadButtonPressed(gamePad int32, button device.GamepadButton) bool {
	return rl.IsGamepadButtonPressed(gamePad, int32(button))
}

// IsGamepadButtonReleased returns if given gamepad button is released
func (dmi DeviceManagerImpl) IsGamepadButtonReleased(gamePad int32, button device.GamepadButton) bool {
	return rl.IsGamepadButtonReleased(gamePad, int32(button))
}

// GetGamepadStickMovement return the movement, -1..1, for a given gamepad stick
func (dmi DeviceManagerImpl) GetGamepadStickMovement(gamePad int32, stick device.GamepadStick) geometry.Point {
	axisX := engineStickToRayXAxis[stick]
	axisY := engineStickToRayYAxis[stick]
	x := rl.GetGamepadAxisMovement(gamePad, axisX)
	y := rl.GetGamepadAxisMovement(gamePad, axisY)
	return geometry.Point{X: x, Y: y}
}
