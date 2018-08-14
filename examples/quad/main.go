package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/maxfish/gojira2d/pkg/app"
	g "github.com/maxfish/gojira2d/pkg/graphics"
)

func main() {
	app.Init(640, 480, "Quad")
	defer app.Terminate()

	app.Context.Camera2D.SetCentered(true)

	Quad := g.NewQuadPrimitive(mgl32.Vec3{0, 0, 0}, mgl32.Vec2{200, 200})
	Quad.SetAnchorToCenter()
	Quad.SetTexture(g.NewTextureFromFile("examples/assets/texture.png"))

	app.MainLoop(func(speed float64) {
		// NOP
	}, func() {
		Quad.Draw(app.Context)
	})
}
