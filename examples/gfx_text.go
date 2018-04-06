package main

import (
	a "gojira2d/pkg/app"
	"gojira2d/pkg/text"
)

func main() {
	app := a.InitApp(640, 480, true, "Quad")
	defer a.TerminateApp()

	font := text.NewFontFromFiles(
		"examples/assets/fonts/mono.fnt",
		"examples/assets/fonts/mono.png",
	)
	t := font.RenderText("abcd")

	app.MainLoop(func(speed float64) {}, func() {
		t.EnqueueForDrawing(app.Context)
	})
}
