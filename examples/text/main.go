package main

import (
	"fmt"
	a "gojira2d/pkg/app"
	"gojira2d/pkg/graphics"
	"gojira2d/pkg/ui"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	app := a.InitApp(800, 600, true, "Text")
	app.Context.SetClearColor(graphics.Color{0.3, 0.3, 0.3, 1})
	defer a.TerminateApp()

	font := ui.NewFontFromFiles(
		"examples/assets/fonts/roboto-mono-regular.fnt",
		"examples/assets/fonts/roboto-mono-regular.png",
	)

	font2 := ui.NewFontFromFiles(
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
		fp := ui.FontPropSmall
		fp.Color = graphics.Color{
			rand.Float32(),
			rand.Float32(),
			rand.Float32(),
			0.6 + 0.4*rand.Float32(),
		}
		tc[i] = f.NewText(
			"The quick brown fox jumps over the lazy dog",
			mgl32.Vec3{-400, -300 + vo, 0},
			mgl32.Vec2{j * 8, j * 8},
			fp,
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
