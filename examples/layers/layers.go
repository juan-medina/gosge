package main

import (
	"fmt"
	"github.com/juan-medina/goecs"
	"github.com/juan-medina/gosge"
	"github.com/juan-medina/gosge/components/color"
	"github.com/juan-medina/gosge/components/effects"
	"github.com/juan-medina/gosge/components/geometry"
	"github.com/juan-medina/gosge/components/shapes"
	"github.com/juan-medina/gosge/components/sprite"
	"github.com/juan-medina/gosge/components/ui"
	"github.com/juan-medina/gosge/options"
	"github.com/rs/zerolog/log"
	"math/rand"
)

// game options
var opt = options.Options{
	Title:      "GOSGE Layers Game",
	BackGround: color.Black,
	Icon:       "resources/icon.png",
	// Uncomment this for using windowed mode
	// Windowed: true,
	// Width:    2048,
	// Height:   1536,
}

func main() {
	// create and load the game
	if err := gosge.Run(opt, load); err != nil {
		log.Fatal().Err(err).Msg("error running the game")
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
	uiText *goecs.Entity

	// designResolution is how our game is designed
	designResolution = geometry.Size{Width: 1920, Height: 1080}
)

// game constants
const (
	itemsToAdd   = 200                     // itemsToAdd is how many items we are going to add
	timeToChange = float32(10)             // timeToChange is how many seconds to change layers
	uiFontSize   = 75                      // uiFontSize is the size for UI font
	itemFontSize = 50                      // itemFontSize is the size for Items font
	fontName     = "resources/go_mono.fnt" // fontName is the font that we use
)

// load the game
func load(eng *gosge.Engine) error {
	// Preload font
	if err := eng.LoadFont(fontName); err != nil {
		return err
	}

	// get the world
	world := eng.World()

	gameScale := eng.GetScreenSize().CalculateScale(designResolution)

	// preload sprite sheet
	if err := eng.LoadSpriteSheet("resources/gamer.json"); err != nil {
		return err
	}

	ss := eng.GetScreenSize()

	// calculate the UI Points
	boxWidth := ss.Width / gameScale.Max * 0.85
	boxStartX := ((designResolution.Width * gameScale.Point.X) - (boxWidth * gameScale.Max)) / 2

	// set the Point with scale
	boxSize := geometry.Size{Width: boxWidth, Height: uiFontSize}
	uiPos := geometry.Point{X: boxStartX, Y: 0}
	uiTextPos := geometry.Point{X: (designResolution.Width / 2) * gameScale.Point.X, Y: boxSize.Height / 2 * gameScale.Point.Y}

	// add the top box
	world.AddEntity(
		shapes.SolidBox{Size: boxSize, Scale: gameScale.Max},
		uiPos,
		color.DarkBlue.Alpha(200),
		uiLayer,
	)

	// add the text
	uiText = world.AddEntity(
		ui.Text{
			String:     "xx is Front",
			VAlignment: ui.MiddleVAlignment,
			HAlignment: ui.CenterHAlignment,
			Font:       fontName,
			Size:       uiFontSize * gameScale.Max,
		},
		uiTextPos,
		color.SkyBlue,
		uiLayer,
	)

	// add the bottom text
	world.AddEntity(
		ui.Text{
			String:     "press <ESC> to close",
			HAlignment: ui.CenterHAlignment,
			VAlignment: ui.BottomVAlignment,
			Font:       fontName,
			Size:       uiFontSize * gameScale.Max,
		},
		geometry.Point{
			X: designResolution.Width / 2 * gameScale.Point.X,
			Y: designResolution.Height * gameScale.Point.Y,
		},
		effects.AlternateColor{
			Time: .25,
			From: color.DarkBlue,
			To:   color.DarkBlue.Alpha(150),
		},
		uiLayer,
	)

	// add the items
	addItems(itemsToAdd, world, gameScale)

	// at the layout system
	world.AddSystem(swapLayersOnTimeSystem)
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
func addItems(toAdd int, world *goecs.World, scl geometry.Scale) {
	// for as many items we like to add
	for i := 0; i < toAdd; i++ {
		// pick a random group
		gn := int32(rand.Float32() * float32(len(groups)))
		gr := groups[gn]

		// Point is random in within the screen
		x := rand.Float32() * designResolution.Width * scl.Point.X
		y := rand.Float32() * designResolution.Height * scl.Point.Y

		it := itemType(rand.Float32() * float32(totalItemsTypes))

		// item Point
		pos := geometry.Point{X: x, Y: y}

		switch it {
		case itemTypeText:
			// text shadow Point
			stp := geometry.Point{X: x + 5, Y: y + 5}

			// add the texts shadow
			world.AddEntity(
				ui.Text{
					String:     gr.name,
					VAlignment: ui.MiddleVAlignment,
					HAlignment: ui.CenterHAlignment,
					Font:       fontName,
					Size:       itemFontSize * scl.Min,
				},
				stp,
				color.DarkGray.Alpha(127),
				gr.layer,
			)

			// add the text
			world.AddEntity(
				ui.Text{
					String:     gr.name,
					VAlignment: ui.MiddleVAlignment,
					HAlignment: ui.CenterHAlignment,
					Font:       fontName,
					Size:       itemFontSize * scl.Min,
				},
				pos,
				gr.clr,
				gr.layer,
			)
		case itemTypeSprite:
			world.AddEntity(
				sprite.Sprite{
					Sheet: "resources/gamer.json",
					Name:  "gamer.png",
					Scale: 0.25 * scl.Min,
				},
				pos,
				gr.clr,
				gr.layer,
			)
		}
	}
}

var (
	time float32 // time to change layers
)

// Update the system checking when need to swap layers
func swapLayersOnTimeSystem(w *goecs.World, delta float32) error {
	// increase time
	time += delta

	// if we need to change
	if time > timeToChange {
		// reset time
		time = 0
		// swap layers
		if err := swapLayers(w, 1, 2); err != nil {
			return err
		}
	}

	// update the ui
	if err := updateUI(w); err != nil {
		return err
	}

	return nil
}

// swapLayers swap the layers of two groups
func swapLayers(w *goecs.World, old, new int) error {
	for it := w.Iterator(effects.TYPE.Layer); it != nil; it = it.Next() {
		ent := it.Value()
		ly := effects.Get.Layer(ent)
		if ly.Depth == groups[old].layer.Depth {
			ent.Set(groups[new].layer)
		} else if ly.Depth == groups[new].layer.Depth {
			ent.Set(groups[old].layer)
		}
	}

	// update the groups
	aux := groups[old].layer
	groups[old].layer = groups[new].layer
	groups[new].layer = aux

	return nil
}

// updateUI update the ui telling what group is top and when is going to change
func updateUI(_ *goecs.World) error {
	var top = 0
	for i := 0; i < 3; i++ {
		if groups[i].layer.Depth == topLayer.Depth {
			top = i
			break
		}
	}

	txt := ui.Get.Text(uiText)
	txt.String = fmt.Sprintf("top layer %s, change in %02.02f", groups[top].name, timeToChange-time)
	uiText.Set(txt)

	return nil
}
