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
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"math/rand"
	"reflect"
)

// game options
var opt = options.Options{
	Title: "Layers Game",
	Size: geometry.Size{
		Width:  1920,
		Height: 1080,
	},
	ClearColor: color.Black,
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
			clr:   color.White,
			name:  "White",
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
)

// game constants
const (
	itemsToAdd   = 100         // itemsToAdd is how many items we are going to add
	timeToChange = float32(10) // timeToChange is how many seconds to change layers
)

// item type
type itemType int

// define item types
const (
	itemTypeText = itemType(iota)
	itemTypeSprite
	totalItemsTypes
)

// addItem add a set of random items to the world
func addItems(toAdd int, wld *world.World) {
	// for as many items we like to add
	for i := 0; i < toAdd; i++ {
		// pick a random group
		gn := int32(rand.Float32() * float32(len(groups)))
		gr := groups[gn]

		// position is random in within the screen
		x := rand.Float32() * opt.Size.Width
		y := rand.Float32() * opt.Size.Height

		it := itemType(rand.Float32() * float32(totalItemsTypes))

		// item position
		pos := geometry.Position{X: x, Y: y}

		switch it {
		case itemTypeText:
			// text shadow position
			stp := geometry.Position{X: x + 5, Y: y + 5}

			// add the texts shadow
			wld.Add(entity.New(
				text.Text{String: gr.name, VAlignment: text.MiddleVAlignment, HAlignment: text.CenterHAlignment},
				relativePosition{original: stp, size: 100},
				color.DarkGray.Alpha(127),
				gr.layer,
			))

			// add the text
			wld.Add(entity.New(
				text.Text{String: gr.name, VAlignment: text.MiddleVAlignment, HAlignment: text.CenterHAlignment},
				relativePosition{original: pos, size: 100},
				gr.clr,
				gr.layer,
			))
		case itemTypeSprite:
			wld.Add(entity.New(
				sprite.Sprite{Sheet: "resources/gamer.json", Name: "gamer.png"},
				relativePosition{original: pos, size: 0.25},
				gr.clr,
				gr.layer,
			))
		}
	}
}

// load the game
func load(eng engine.Engine) error {
	// get the world
	wld := eng.World()

	// preload sprite sheet
	if err := eng.LoadSpriteSheet("resources/gamer.json"); err != nil {
		return err
	}

	// calculate the UI positions
	boxSize := geometry.Size{Width: opt.Size.Width / 1.15, Height: opt.Size.Height / 12}
	uiPos := geometry.Position{X: (opt.Size.Width / 2) - (boxSize.Width / 2), Y: 0}
	uiTextPos := geometry.Position{X: uiPos.X, Y: boxSize.Height / 2}

	// add the top box
	wld.Add(entity.New(
		shapes.Box{
			Size:  boxSize,
			Scale: 1,
		},
		relativePosition{original: uiPos},
		color.DarkBlue.Alpha(200),
		uiLayer,
	))

	// add the text
	uiText = wld.Add(entity.New(
		text.Text{String: "xx is Front", VAlignment: text.MiddleVAlignment, HAlignment: text.LeftHAlignment},
		relativePosition{original: uiTextPos, size: 80},
		color.SkyBlue,
		uiLayer,
	))

	// add the items
	addItems(itemsToAdd, wld)

	// at the layout system
	wld.AddSystem(&layoutSystem{})
	return nil
}

func main() {
	// create and load the game
	if err := game.Run(opt, load); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

// relativePosition contains the original position and size for changing layout on screen size change
type relativePosition struct {
	original geometry.Position
	size     float32
}

// relativeType is the reflect.Type of your relativePosition
var relativeType = reflect.TypeOf(relativePosition{})

// layoutSystem will resize and re-position the game on screen size change
type layoutSystem struct {
	time float32
}

// swapLayers swap the layers of two groups
func (ls *layoutSystem) swapLayers(w *world.World, old, new int) error {
	for _, ent := range w.Entities() {
		if ent.Contains(relativeType, effects.TYPE.Layer) {
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
func (ls *layoutSystem) updateUI(_ *world.World) error {
	var top = 0
	for i := 0; i < 3; i++ {
		if groups[i].layer.Depth == topLayer.Depth {
			top = i
			break
		}
	}

	txt := text.Get(uiText)
	txt.String = fmt.Sprintf("top layer %s, change in %02.02f", groups[top].name, timeToChange-ls.time)
	uiText.Set(txt)

	return nil
}

// Update the system checking when need to swap layers
func (ls *layoutSystem) Update(w *world.World, delta float32) error {
	// increase time
	ls.time += delta

	// if we need to change
	if ls.time > timeToChange {
		// reset time
		ls.time = 0
		// swap layers
		if err := ls.swapLayers(w, 1, 2); err != nil {
			return err
		}
	}

	// update the ui
	if err := ls.updateUI(w); err != nil {
		return err
	}

	return nil
}

// layout game elements
func (ls layoutSystem) layout(w *world.World, ssc events.ScreenSizeChangeEvent) error {
	// get the entities that has a relativePosition
	for _, ent := range w.Entities(relativeType) {
		// get the original position
		rl := ent.Get(relativeType).(relativePosition)

		// calculate the new position
		pos := geometry.Position{}
		pos.X = rl.original.X * ssc.Scale.Point.X
		pos.Y = rl.original.Y * ssc.Scale.Point.Y
		ent.Set(pos)

		// if is a text scale text size and spacing
		if ent.Contains(text.TYPE) {
			tx := text.Get(ent)
			tx.Size = rl.size * ssc.Scale.Min
			tx.Spacing = tx.Size / 4
			ent.Set(tx)
		}

		// if is a Box shape update its scale
		if ent.Contains(shapes.TYPE.Box) {
			box := shapes.Get.Box(ent)
			box.Scale = ssc.Scale.Min
			ent.Set(box)
		}

		// if is a Sprite update its scale
		if ent.Contains(sprite.TYPE) {
			spr := sprite.Get(ent)
			spr.Scale = rl.size * ssc.Scale.Min
			ent.Set(spr)
		}
	}
	return nil
}

// Notify is trigger on events
func (ls *layoutSystem) Notify(w *world.World, event interface{}, _ float32) error {
	// switch event type
	switch e := event.(type) {
	//if the screen size change
	case events.ScreenSizeChangeEvent:
		// re-layout the game
		return ls.layout(w, e)
	}
	return nil
}
