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
	"github.com/juan-medina/gosge/pkg/components/device"
	"github.com/juan-medina/gosge/pkg/options"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"strings"
)

// Init the rendering device
func (rr *RenderImpl) Init(opt options.Options) {
	rr.saveOpts = opt
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

	w := rl.GetMonitorWidth(opt.Monitor)
	h := rl.GetMonitorHeight(opt.Monitor)
	rl.InitWindow(int32(w), int32(h), opt.Title)

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
		if !strings.Contains(runtime.GOOS, "linux") {
			rl.ToggleFullscreen()
		}
	}
}

// End the rendering device
func (rr RenderImpl) End() {
	rl.CloseWindow()
}

// BeginFrame for rendering
func (rr RenderImpl) BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(rr.color2RayColor(rr.saveOpts.BackGround))
}

// EndFrame for rendering
func (rr RenderImpl) EndFrame() {
	rl.EndDrawing()
}

// ShouldClose returns if th engine should close
func (rr RenderImpl) ShouldClose() bool {
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
}

// GetKeyStatus returns the device.KeyStatus for a given device.Key
func (rr RenderImpl) GetKeyStatus(key device.Key) device.KeyStatus {
	status := device.KeyStatus{}

	if v, ok := engineKeyToRayKey[key]; ok {
		status.Down = rl.IsKeyDown(v)
		status.Up = rl.IsKeyUp(v)
		status.Pressed = rl.IsKeyPressed(v)
		status.Released = rl.IsKeyReleased(v)
	}

	return status
}
