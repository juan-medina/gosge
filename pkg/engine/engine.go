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

package engine

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/juan-medina/goecs/pkg/entitiy"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components"
	"github.com/juan-medina/gosge/pkg/render"
	"github.com/juan-medina/gosge/pkg/systems"
	"log"
)

const (
	renderingGroup = "RENDERING_GROUP"
)

type gameEngine struct {
	game  Game
	world *world.World
}

//goland:noinspection GoUnusedParameter
func (eng *gameEngine) layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	settings := eng.world.Entity(components.GameSettingsType).Get(components.GameSettingsType).(components.GameSettings)
	return settings.Width, settings.Height
}

//goland:noinspection GoUnusedParameter
func (eng *gameEngine) update(screen *ebiten.Image) error {
	eng.World().Update()
	return nil
}

func (eng *gameEngine) draw(screen *ebiten.Image) {
	context := eng.World().Entity(render.ContextType)
	context.Set(render.NewContext(screen))

	eng.World().UpdateGroup(renderingGroup)
}

func (eng *gameEngine) run() {
	eng.world.Add(entitiy.New().Add(render.NewContext(nil)))
	eng.world.AddSystemToGroup(systems.TextRenderingSystem(), renderingGroup)

	eng.game.Init(eng)

	settings := components.GameSettings{
		Width:  640,
		Height: 480,
		Title:  "Go Simple Game Engine",
	}

	if settingsEnt := eng.world.Entity(components.GameSettingsType); settingsEnt == nil {
		eng.world.Add(entitiy.New().Add(settings))
	} else {
		settings = settingsEnt.Get(components.GameSettingsType).(components.GameSettings)
	}

	ebiten.SetWindowSize(settings.Width, settings.Height)
	ebiten.SetWindowTitle(settings.Title)

	if err := ebiten.RunGame(newEbitenGameWrapper(eng)); err != nil {
		log.Fatal(err)
	}
}

func (eng *gameEngine) World() *world.World {
	return eng.world
}

func newEngine(game Game) Engine {
	return &gameEngine{
		game:  game,
		world: world.New(),
	}
}

func Run(game Game) {
	newEngine(game).run()
}
