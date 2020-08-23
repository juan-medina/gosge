package main

import (
	"fmt"
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/sprite"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"math/rand"
)

// game options
var opt = options.Options{
	Title:      "Layers Game",
	BackGround: color.Black,
}

func main() {
	// create and load the game
	if err := game.Run(opt, load); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

// group will contain a set of items
type group struct {
	clr   color.Solid   // clr is the group color
	name  string        // name is the group name
	layer effects.Layer // layer is group effect.Layer
}

var (
	backLayer   = effects.Layer{Depth: 3} // backLayer is at the far of the screen
	middleLayer = effects.Layer{Depth: 2} // middleLayer is at the middle of the screen
	topLayer    = effects.Layer{Depth: 1} // topLayer is at the near of the screen
	uiLayer     = effects.Layer{Depth: 0} // uiLayer is on top of everything

	// groups is an slice of groups
	groups = []group{
		{
			clr:   color.SkyBlue,
			name:  "SkyBlue",
			layer: backLayer,
		},
		{
			clr:   color.Red,
			name:  "Red",
			layer: topLayer,
		},
		{
			clr:   color.Green,
			name:  "Green",
			layer: middleLayer,
		},
	}

	// uiText contains the text on top of the screen
	uiText *entity.Entity

	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
)

// game constants
const (
	itemsToAdd   = 200         // itemsToAdd is how many items we are going to add
	timeToChange = float32(10) // timeToChange is how many seconds to change layers
	uiFontSize   = 75          // uiFontSize is the size for UI font
	itemFontSize = 50          // itemFontSize is the size for Items font
)

// load the game
func load(eng engine.Engine) error {
	// get the world
	wld := eng.World()

	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// preload sprite sheet
	if err := eng.LoadSpriteSheet("resources/gamer.json"); err != nil {
		return err
	}

	// calculate the UI positions
	boxWidth := designResolution.Width * 0.70
	boxStartX := (designResolution.Width / 2) - (boxWidth / 2)

	// set the position with scale
	boxSize := geometry.Size{Width: boxWidth * gameScale.Point.X, Height: uiFontSize * gameScale.Point.Y}
	uiPos := geometry.Position{X: boxStartX * gameScale.Point.X, Y: 0}
	uiTextPos := geometry.Position{X: uiPos.X, Y: boxSize.Height / 2}

	// add the top box
	wld.Add(entity.New(
		shapes.Box{Size: boxSize, Scale: 1},
		uiPos,
		color.DarkBlue.Alpha(200),
		uiLayer,
	))

	// add the text
	uiText = wld.Add(entity.New(
		text.Text{
			String:     "xx is Front",
			VAlignment: text.MiddleVAlignment,
			HAlignment: text.LeftHAlignment,
			Size:       uiFontSize * gameScale.Min,
			Spacing:    (uiFontSize / 10) * gameScale.Min,
		},
		uiTextPos,
		color.SkyBlue,
		uiLayer,
	))

	// add the bottom text
	wld.Add(entity.New(
		text.Text{
			String:     "press <ESC> to close",
			HAlignment: text.CenterHAlignment,
			VAlignment: text.BottomVAlignment,
			Size:       uiFontSize * gameScale.Min,
			Spacing:    (uiFontSize / 10) * gameScale.Min,
		},
		geometry.Position{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height * gameScale.Point.Y,
		},
		effects.AlternateColor{
			Time: .25,
			From: color.DarkBlue,
			To:   color.DarkBlue.Alpha(150),
		},
		uiLayer,
	))

	// add the items
	addItems(itemsToAdd, wld, gameScale)

	// at the layout system
	wld.AddSystem(&swapLayersOnTimeSystem{})
	return nil
}

// item type
type itemType int

// define item types
const (
	itemTypeText = itemType(iota)
	itemTypeSprite
	totalItemsTypes
)

// addItem add a set of random items to the world
func addItems(toAdd int, wld *world.World, scl geometry.Scale) {
	// for as many items we like to add
	for i := 0; i < toAdd; i++ {
		// pick a random group
		gn := int32(rand.Float32() * float32(len(groups)))
		gr := groups[gn]

		// position is random in within the screen
		x := rand.Float32() * designResolution.Width * scl.Point.X
		y := rand.Float32() * designResolution.Height * scl.Point.Y

		it := itemType(rand.Float32() * float32(totalItemsTypes))

		// item position
		pos := geometry.Position{X: x, Y: y}

		switch it {
		case itemTypeText:
			// text shadow position
			stp := geometry.Position{X: x + 5, Y: y + 5}

			// add the texts shadow
			wld.Add(entity.New(
				text.Text{
					String:     gr.name,
					VAlignment: text.MiddleVAlignment,
					HAlignment: text.CenterHAlignment,
					Size:       itemFontSize * scl.Min,
					Spacing:    itemFontSize / 10 * scl.Min,
				},
				stp,
				color.DarkGray.Alpha(127),
				gr.layer,
			))

			// add the text
			wld.Add(entity.New(
				text.Text{
					String:     gr.name,
					VAlignment: text.MiddleVAlignment,
					HAlignment: text.CenterHAlignment,
					Size:       itemFontSize * scl.Min,
					Spacing:    itemFontSize / 10 * scl.Min,
				},
				pos,
				gr.clr,
				gr.layer,
			))
		case itemTypeSprite:
			wld.Add(entity.New(
				sprite.Sprite{
					Sheet: "resources/gamer.json",
					Name:  "gamer.png",
					Scale: 0.25 * scl.Min,
				},
				pos,
				gr.clr,
				gr.layer,
			))
		}
	}
}

// swapLayersOnTimeSystem will change the groups layers base on time
type swapLayersOnTimeSystem struct {
	time float32
}

// Update the system checking when need to swap layers
func (sls *swapLayersOnTimeSystem) Update(w *world.World, delta float32) error {
	// increase time
	sls.time += delta

	// if we need to change
	if sls.time > timeToChange {
		// reset time
		sls.time = 0
		// swap layers
		if err := sls.swapLayers(w, 1, 2); err != nil {
			return err
		}
	}

	// update the ui
	if err := sls.updateUI(w); err != nil {
		return err
	}

	return nil
}

// Notify is trigger on events
func (sls *swapLayersOnTimeSystem) Notify(_ *world.World, _ interface{}, _ float32) error {
	return nil
}

// swapLayers swap the layers of two groups
func (sls *swapLayersOnTimeSystem) swapLayers(w *world.World, old, new int) error {
	for _, ent := range w.Entities() {
		if ent.Contains(effects.TYPE.Layer) {
			ly := effects.Get.Layer(ent)
			if ly.Depth == groups[old].layer.Depth {
				ent.Set(groups[new].layer)
			} else if ly.Depth == groups[new].layer.Depth {
				ent.Set(groups[old].layer)
			}
		}
	}

	// update the groups
	aux := groups[old].layer
	groups[old].layer = groups[new].layer
	groups[new].layer = aux

	return nil
}

// updateUI update the ui telling what group is top and when is going to change
func (sls *swapLayersOnTimeSystem) updateUI(_ *world.World) error {
	var top = 0
	for i := 0; i < 3; i++ {
		if groups[i].layer.Depth == topLayer.Depth {
			top = i
			break
		}
	}

	txt := text.Get(uiText)
	txt.String = fmt.Sprintf("top layer %s, change in %02.02f", groups[top].name, timeToChange-sls.time)
	uiText.Set(txt)

	return nil
}
