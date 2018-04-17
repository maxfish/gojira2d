package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/maxfish/gojira2d/pkg/graphics"
	"github.com/maxfish/gojira2d/pkg/ui"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/maxfish/gojira2d/pkg/app"
)

func main() {
	app.Init(800, 600, true, "Text")
	app.SetClearColor(graphics.Color{0, 0, 0, 1})
	defer app.Terminate()

	font := ui.NewFontFromFiles(
		"mono",
		"examples/assets/fonts/roboto-mono-regular.fnt",
		"examples/assets/fonts/roboto-mono-regular.png",
	)

	font2 := ui.NewFontFromFiles(
		"regular",
		"examples/assets/fonts/roboto-regular.fnt",
		"examples/assets/fonts/roboto-regular.png",
	)

	var f *ui.Font

	tc := [20]*ui.Text{}
	vo := float32(0)
	for i := range tc {
		j := float32(i)
		if i%2 == 0 {
			f = font
		} else {
			f = font2
		}
		color := graphics.Color{
			rand.Float32(),
			rand.Float32(),
			rand.Float32(),
			0.6 + 0.4*rand.Float32(),
		}
		tc[i] = ui.NewText(
			"The quick brown fox jumps over the lazy dog",
			f,
			mgl32.Vec3{-400, -300 + vo, 0},
			mgl32.Vec2{j * 8, j * 8},
			color,
			mgl32.Vec4{0, 0, 0, -.17},
		)
		vo += j * 6
	}

	app.MainLoop(
		func(speed float64) {
			tc[4].SetText(fmt.Sprintf("%v", time.Now()))
		},
		func() {
			for _, t := range tc {
				t.EnqueueForDrawing(app.Context)
			}
		})
}
