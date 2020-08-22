package main

import (
	"fmt"
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"github.com/juan-medina/gosge/pkg/components/color"
	"github.com/juan-medina/gosge/pkg/components/effects"
	"github.com/juan-medina/gosge/pkg/components/geometry"
	"github.com/juan-medina/gosge/pkg/components/shapes"
	"github.com/juan-medina/gosge/pkg/components/text"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/events"
	"github.com/juan-medina/gosge/pkg/game"
	"github.com/juan-medina/gosge/pkg/options"
	"log"
	"math/rand"
	"reflect"
)

var opt = options.Options{
	Title: "Layers Game",
	Size: geometry.Size{
		Width:  1920,
		Height: 1080,
	},
	ClearColor: color.Black,
}

var (
	backLayer   = effects.Layer{Depth: 3}
	middleLayer = effects.Layer{Depth: 2}
	topLayer    = effects.Layer{Depth: 1}
	uiLayer     = effects.Layer{Depth: 0}

	clrs   []color.Solid
	cns    []string
	lyr    []effects.Layer
	uiText *entity.Entity
)

const (
	itemsToAdd   = 100
	timeToChange = float32(10)
)

func addItems(toAdd int, wld *world.World) {
	for i := 0; i < toAdd; i++ {
		cn := int32(rand.Float32() * 3)
		x := rand.Float32() * opt.Size.Width
		y := rand.Float32() * opt.Size.Height

		txp := geometry.Position{X: x, Y: y}
		stp := geometry.Position{X: x + 5, Y: y + 5}

		wld.Add(entity.New(
			text.Text{String: cns[cn], VAlignment: text.MiddleVAlignment, HAlignment: text.CenterHAlignment},
			relativePosition{original: stp, size: 100},
			color.DarkGray.Alpha(127),
			lyr[cn],
		))

		wld.Add(entity.New(
			text.Text{String: cns[cn], VAlignment: text.MiddleVAlignment, HAlignment: text.CenterHAlignment},
			relativePosition{original: txp, size: 100},
			clrs[cn],
			lyr[cn],
		))
	}
}

func load(eng engine.Engine) error {
	wld := eng.World()

	clrs = make([]color.Solid, 3)
	cns = make([]string, 3)
	lyr = make([]effects.Layer, 3)

	clrs[0] = color.Green
	cns[0] = "Green"
	lyr[0] = backLayer

	clrs[1] = color.Red
	cns[1] = "Red"
	lyr[1] = topLayer

	clrs[2] = color.Yellow
	cns[2] = "Yellow"
	lyr[2] = middleLayer

	boxSize := geometry.Size{Width: opt.Size.Width / 1.5, Height: opt.Size.Height / 12}
	uiPos := geometry.Position{X: (opt.Size.Width / 2) - (boxSize.Width / 2), Y: 0}
	uiTextPos := geometry.Position{X: opt.Size.Width / 2, Y: boxSize.Height / 2}

	wld.Add(entity.New(
		shapes.Box{
			Size:  boxSize,
			Scale: 1,
		},
		relativePosition{original: uiPos},
		color.White.Alpha(127),
		uiLayer,
	))

	uiText = wld.Add(entity.New(
		text.Text{String: "xx is Front", VAlignment: text.MiddleVAlignment, HAlignment: text.CenterHAlignment},
		relativePosition{original: uiTextPos, size: 80},
		color.Black,
		uiLayer,
	))

	addItems(itemsToAdd, wld)

	wld.AddSystem(&layoutSystem{})
	return nil
}

func main() {
	if err := game.Run(opt, load); err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}

type relativePosition struct {
	original geometry.Position
	size     float32
}

var relativeType = reflect.TypeOf(relativePosition{})

type layoutSystem struct {
	time float32
	ssc  events.ScreenSizeChangeEvent
}

func (ls *layoutSystem) swapLayers(w *world.World, old, new int) error {
	for _, ent := range w.Entities() {
		if ent.Contains(relativeType, effects.TYPE.Layer) {
			ly := effects.Get.Layer(ent)
			if ly.Depth == lyr[old].Depth {
				ent.Set(lyr[new])
			} else if ly.Depth == lyr[new].Depth {
				ent.Set(lyr[old])
			}
		}
	}

	aux := lyr[old]
	lyr[old] = lyr[new]
	lyr[new] = aux

	return nil
}

func (ls *layoutSystem) changeUIText(_ *world.World) error {
	var top = 0
	for i := 0; i < 3; i++ {
		if lyr[i].Depth == topLayer.Depth {
			top = i
			break
		}
	}

	txt := text.Get(uiText)
	txt.String = fmt.Sprintf("top %s, change %.f", cns[top], timeToChange-ls.time)
	uiText.Set(txt)

	return nil
}

func (ls *layoutSystem) Update(w *world.World, delta float32) error {
	ls.time += delta

	if ls.time > timeToChange {
		ls.time = 0
		if err := ls.swapLayers(w, 1, 2); err != nil {
			return err
		}
	}

	if err := ls.changeUIText(w); err != nil {
		return err
	}
	return nil
}

func (ls layoutSystem) doLayout(w *world.World) error {
	for _, ent := range w.Entities(relativeType) {
		rl := ent.Get(relativeType).(relativePosition)
		pos := geometry.Position{}
		pos.X = rl.original.X * ls.ssc.Scale.Point.X
		pos.Y = rl.original.Y * ls.ssc.Scale.Point.Y
		ent.Set(pos)

		if ent.Contains(text.TYPE) {
			tx := text.Get(ent)
			tx.Size = rl.size * ls.ssc.Scale.Min
			tx.Spacing = tx.Size / 4
			ent.Set(tx)
		}

		if ent.Contains(shapes.TYPE.Box) {
			box := shapes.Get.Box(ent)
			box.Scale = ls.ssc.Scale.Min
			ent.Set(box)
		}
	}
	return nil
}

func (ls *layoutSystem) Notify(w *world.World, event interface{}, _ float32) error {
	switch e := event.(type) {
	case events.ScreenSizeChangeEvent:
		ls.ssc = e
		return ls.doLayout(w)
	}

	return nil
}
