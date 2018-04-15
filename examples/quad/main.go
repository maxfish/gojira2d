package main

import (
	"github.com/go-gl/mathgl/mgl32"
	a "github.com/maxfish/gojira2d/app"
	g "github.com/maxfish/gojira2d/graphics"
)

func main() {
	app := a.InitApp(640, 480, true, "Quad")
	defer a.TerminateApp()

	Quad := g.NewQuadPrimitive(mgl32.Vec3{0, 0, 0}, mgl32.Vec2{200, 200})
	Quad.SetAnchorToCenter()
	Quad.SetTexture(g.NewTextureFromFile("examples/assets/texture.png"))

	app.MainLoop(func(speed float64) {
		// NOP
	}, func() {
		Quad.EnqueueForDrawing(app.Context)
	})
}
