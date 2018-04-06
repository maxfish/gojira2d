package main

import (
	a "gojira2d/pkg/app"
	"gojira2d/pkg/text"
	"github.com/go-gl/mathgl/mgl32"
	"gojira2d/pkg/graphics"
)

func main() {
	app := a.InitApp(800, 600, true, "Text")
	defer a.TerminateApp()

	font := text.NewFontFromFiles(
		"examples/assets/fonts/mono.fnt",
		"examples/assets/fonts/mono.png",
	)

	font2 := text.NewFontFromFiles(
		"examples/assets/fonts/regular.fnt",
		"examples/assets/fonts/regular.png",
	)

	var f *text.Font

	tc := [20]*graphics.Primitive2D{}
	vo := float32(0)
	for i := range tc {
		j := float32(i)
		if i % 2 == 0 {
			f = font
		} else {
			f = font2
		}
		tc[i] = f.RenderText(
			"The quick brown fox jumps over the lazy dog",
			mgl32.Vec3{-400, -300+vo, 0},
			mgl32.Vec2{j*8, j*8},
			text.FontPropSmall,
		)
		vo += j * 6
	}

	app.MainLoop(func(speed float64) {}, func() {
		for _, t := range tc {
			t.EnqueueForDrawing(app.Context)
		}
	})
}
